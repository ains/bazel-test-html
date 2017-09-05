package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	directory string
	summary   *TestSummary
}

var testCases = []testCase{
	{
		directory: "001",
		summary: &TestSummary{
			TotalTests:  2,
			BuildErrors: "",
			Results: Results{
				PASS: []*Test{
					{
						TestName: "//go/src/pass:go_default_test",
						Time:     0,
						Output:   "pass",
						Status:   PASS,
					},
				},
				FAIL: []*Test{
					{
						TestName: "//go/src/failytest:go_default_test",
						Time:     0,
						Output:   "fail",
						Status:   FAIL,
					},
				},
				SKIP: []*Test{},
			},
		},
	},
}

func TestSummaryParser(t *testing.T) {
	for _, testCase := range testCases {
		t.Logf("Running: (%s)", testCase.directory)

		stdoutErr, err := os.Open(os.DevNull)
		require.NoError(t, err)

		actual, err := Parse("tests/"+testCase.directory, stdoutErr)
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.EqualValues(t, testCase.summary, actual)
	}
}
