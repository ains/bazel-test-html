# bazel-test-html
[![Build Status](https://travis-ci.org/ains/bazel-test-html.svg?branch=master)](https://travis-ci.org/ains/bazel-test-html)

Converts `bazel test` output from rules_go tests into a prettified HTML summary.

![Example HTML output](https://cloud.githubusercontent.com/assets/2838283/18885683/8f53b7fa-84e4-11e6-9ca8-f5e3101acf98.png)

## Installation

Go version 1.1 or higher is required. Install or update using the `go get`
command:

	go get -u github.com/ains/bazel-test-html

## Usage

The bazel-test-html command takes in files containing the stdout and stderr from the `go test` command and the location
of the HTML file to write.

    bazel-test-html [bazel_testlogs_dir] [build_error_file] [output_file]
