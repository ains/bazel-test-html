package main

import (
	"fmt"
	"log"
	"os"

	"io/ioutil"
	"path/filepath"

	"github.com/GeertJohan/go.rice"
	"github.com/ains/bazel-test-html/lib"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Incorrect command line arguments")
		fmt.Println("Usage: bazel-test-html [bazel_testlogs_dir] [build_error_file] [output_file]")
		os.Exit(1)
	}

	testLogsDirectory := os.Args[1]
	buildErrorFile := os.Args[2]
	outputFile := os.Args[3]

	buildErr, err := os.Open(buildErrorFile)
	check(err)

	summary, err := lib.Parse(testLogsDirectory, buildErr)
	check(err)

	templateBox := rice.MustFindBox("template")
	html, err := lib.GenerateHTML(templateBox.MustString("template.html"), summary)
	check(err)

	err = ioutil.WriteFile(outputFile, []byte(html), 0644)
	check(err)

	outputFilePath, err := filepath.Abs(outputFile)
	check(err)

	fmt.Printf("Test results written to '%s'\n", outputFilePath)
}
