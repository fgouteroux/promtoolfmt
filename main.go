/*
Copyright © 2023 François Gouteroux <francois.gouteroux@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

const (
	successExitCode = 0
	failureExitCode = 1
	// Exit code 3 is used for "one or more lint issues detected".
	lintErrExitCode = 3
)

var (
	cliVersion = "0.0.1"
)

// uniqueStringSlice returns unique items in a slice
func uniqueStringSlice(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

// checkMetrics performs a linting pass on input metrics.
// https://github.com/prometheus/prometheus/blob/6ddadd98b44cca7d55b27c20477123afac2201d7/cmd/promtool/main.go#L755
func checkMetrics(input io.Reader) (int, []string) {
	var errors []string
	var buf bytes.Buffer
	tee := io.TeeReader(input, &buf)

	l := promlint.New(tee)
	problems, err := l.Lint()
	if err != nil {
		errors = append(errors, fmt.Sprintf("error while linting: %v", err))
		return failureExitCode, errors
	}

	for _, p := range problems {
		errors = append(errors, fmt.Sprintln(p.Metric, p.Text))
	}

	if len(problems) > 0 {
		return lintErrExitCode, errors
	}

	return successExitCode, errors
}

// metricsFormat
func metricsFormat(input string) string {
	// remove duplicates lines
	strSlice := uniqueStringSlice(strings.Split(input, "\n"))

	// end with a newline
	metrics := fmt.Sprintf("%s\n", strings.Join(strSlice, "\n"))

	// remove return carriage
	metrics = strings.ReplaceAll(metrics, "\r", "")

	return metrics
}

// parseText read text and returns MetricFamily
func parseText(input io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(input)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

// parseMetricFamily read MetricFamily and returns text
func parseMetricFamily(mfs map[string]*dto.MetricFamily) (string, error) {
	var buf bytes.Buffer
	for _, mf := range mfs {
		_, err := expfmt.MetricFamilyToText(&buf, mf)
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func output(silent bool, metrics string, exitCode int, errors []string) {
	if !silent {
		if successExitCode != exitCode {
			for _, error := range errors {
				fmt.Fprintln(os.Stderr, error)
			}
		} else {
			// this part will output metrics text as following convention
			// https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format

			// text to metric family
			mfs, err := parseText(strings.NewReader(metrics))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			// metric family to text
			metricsFmt, err := parseMetricFamily(mfs)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println(metricsFmt)
			}
		}
	}
}

func main() {
	silent := flag.Bool("silent", false, "Silent or quiet mode.")
	version := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *version {
		fmt.Println(cliVersion)
		os.Exit(0)
	}

	var buf bytes.Buffer
	raw := io.TeeReader(os.Stdin, &buf)
	data, _ := io.ReadAll(raw)
	metrics := metricsFormat(string(data))
	exitCode, errors := checkMetrics(strings.NewReader(metrics))

	output(*silent, metrics, exitCode, errors)
	os.Exit(exitCode)
}
