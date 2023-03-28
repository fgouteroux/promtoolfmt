package main

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckMetricsSilent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on windows")
	}

	f, err := os.Open("testdata/metrics-test.prom")
	require.NoError(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)

	metrics := metricsFormat(string(data))

	exitCode, errors := checkMetrics(strings.NewReader(metrics))
	require.Equal(t, 0, exitCode)
	output(false, metrics, exitCode, errors)
}

func TestCheckMetricsNoSilent(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on windows")
	}

	f, err := os.Open("testdata/metrics-test.prom")
	require.NoError(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)

	metrics := metricsFormat(string(data))

	exitCode, errors := checkMetrics(strings.NewReader(metrics))
	require.Equal(t, 0, exitCode)
	require.Equal(t, 0, len(errors))
	output(true, metrics, exitCode, errors)
}

func TestCheckMetricInvalidFloatValue(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on windows")
	}

	f, err := os.Open("testdata/invalid-metrics-float-value.prom")
	require.NoError(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)

	metrics := metricsFormat(string(data))

	exitCode, _ := checkMetrics(strings.NewReader(metrics))
	require.Equal(t, 1, exitCode)
}

func TestCheckMetricInvalidHelpText(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on windows")
	}

	f, err := os.Open("testdata/invalid-metrics-no-help-text.prom")
	require.NoError(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.NoError(t, err)

	metrics := metricsFormat(string(data))

	exitCode, _ := checkMetrics(strings.NewReader(metrics))
	require.Equal(t, 3, exitCode)
}
