package lib

import (
	"path/filepath"
	"os"
	"fmt"
	"encoding/xml"
	"io/ioutil"
	"io"
)

type Results map[string][]*Test

const (
	PASS = "pass"
	FAIL = "fail"
	SKIP = "skip"
)

type Test struct {
	TestName string `json:"test_name"`
	Time     int    `json:"time"`
	Output   string `json:"output"`
	Status   string
}

type TestSummary struct {
	TotalTests  int     `json:"total_tests"`
	BuildErrors string  `json:"build_errors"`
	Results     Results `json:"results"`
}

func Parse(bazelTestLogsDirectory string, buildErrReader io.Reader) (*TestSummary, error) {
	results := Results{
		PASS: []*Test{},
		FAIL: []*Test{},
		SKIP: []*Test{},
	}

	walkRoot, err := filepath.EvalSymlinks(bazelTestLogsDirectory)
	if err != nil {
		return nil, err

	}

	totalTests := 0
	err = filepath.Walk(walkRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Base(path) == "test.xml" {
			totalTests += 1
			test, err := getTestFromXML(path)
			if err != nil {
				return err
			}

			results[test.Status] = append(results[test.Status], test)
		}
		return nil
	})

	buildErrorBytes, err := ioutil.ReadAll(buildErrReader)
	if err != nil {
		return nil, err
	}

	summary := &TestSummary{
		TotalTests:  totalTests,
		Results:     results,
		BuildErrors: string(buildErrorBytes),
	}

	return summary, err
}

type Error struct {
	Message string `xml:"message,attr"`
}

type TestCase struct {
	Duration int `xml:"duration,int64"`
	Error    *Error `xml:"error"`
}

type TestSuite struct {
	Name     string `xml:"name,attr"`
	TestCase *TestCase `xml:"testcase"`
	Output   string `xml:"system-out"`
}

type ResultXML struct {
	XMLName   xml.Name `xml:"testsuites"`
	TestSuite *TestSuite `xml:"testsuite"`
}

func getTestFromXML(testXml string) (*Test, error) {
	xmlContent, err := ioutil.ReadFile(testXml)
	if err != nil {
		return nil, err
	}

	var parsedTest ResultXML
	err = xml.Unmarshal(xmlContent, &parsedTest)
	if err != nil {
		return nil, err
	}

	suite := parsedTest.TestSuite
	if suite == nil {
		return nil, fmt.Errorf("no test suite found in %v", testXml)
	}

	status := PASS
	if suite.TestCase.Error != nil {
		status = FAIL
	}

	// For some reason targets are of the format "<package>/<target>" convert to "//<package>:target"
	pkg := filepath.Dir(suite.Name)
	target := filepath.Base(suite.Name)

	return &Test{
		TestName: fmt.Sprintf("//%s:%s", pkg, target),
		Time:     suite.TestCase.Duration,
		Output:   suite.Output,
		Status:   status,
	}, nil
}
