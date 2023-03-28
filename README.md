# promtoolfmt

promtoolfmt is made for linting and formating [prometheus metrics in text-based format](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format).


## How it works

promtoolfmt remove any duplicate metric/comment, performs a linting pass on metrics and format them using [expfmt](https://pkg.go.dev/github.com/prometheus/common/expfmt#MetricFamilyToText).


## Usage

```shell
Usage of promtoolfmt:
  -silent
      Silent or quiet mode.
  -version
      Show version.
```

## Run it

```shell
$ cat testdata/metrics-test.prom | promtoolfmt
# HELP test_metric_info A metric with a constant '1' value labeled by testlabel.
# TYPE test_metric_info gauge
test_metric_info{testlabel="testvalue1"} 1
test_metric_info{testlabel="testvalue2"} 1
# HELP test_metric_duplicate_info A metric with a constant '1' value labeled by testlabel.
# TYPE test_metric_duplicate_info gauge
test_metric_duplicate_info{testlabel="testvalue1"} 1
# HELP test_metric_with_spaces_info A metric with a constant '1' value labeled by testlabel.
# TYPE test_metric_with_spaces_info gauge
test_metric_with_spaces_info{testlabel="testvalue1",testanotherlabel="testanothervalue1"} 1
test_metric_with_spaces_info{testlabel="testvalue2",testanotherlabel="testanothervalue2"} 1

```

### With silent mode

Use this mode to get exit code only (suppress any output/error)

```shell
$ cat testdata/metrics-test.prom | promtoolfmt -silent
```


## Errors examples

```shell
$ cat testdata/invalid-metrics-no-help-text.prom | promtoolfmt
test_metric_info no help text
```

```shell
$ cat testdata/invalid-metrics-float-value.prom | promtoolfmt
error while linting: text format parsing error in line 4: expected float as value, got "1m"
```
