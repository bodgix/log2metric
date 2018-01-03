package main

import (
	"testing"
)

func TestParseLogFile(t *testing.T) {
	// Setup
	regexp := `time:\s+(?P<resp_time>[^\s]+)\s+bytes:\s+(?P<resp_bytes>\d+)`
	lines := []string{"time: 0.567 bytes: 5000", "time: 0.432 bytes: 4567", "bad_metric: 0.999", "time: 0.432as bytes: 4567"}
	logCh := make(chan string)
	metricsCh := make(chan metric)
	expected := []metric{
		metric{simple, "resp_time", float64(0.567)},
		metric{simple, "resp_bytes", float64(5000)},
		metric{simple, "resp_time", float64(0.432)},
		metric{simple, "resp_bytes", float64(4567)},
		metric{simple, "resp_bytes", float64(4567)},
	}
	parser := new(simpleLogParser)

	// Test
	go parser.parseLogFile(logCh, metricsCh, regexp)
	go func() {
		for _, line := range lines {
			logCh <- line
		}
		close(logCh)
	}()
	i := 0
	for m := range metricsCh {
		if expected[i] != m {
			t.Error("Expected: ", expected[i], "Got: ", m)
		}
		i++
	}
}

func TestHistogramLogParser(t *testing.T) {
	regexp := `http:\s+(?P<http>[\d]+)`
	lines := []string{"http: 200", "http: 401", "http: 401", "http: 200", "http: 500"}
	logCh := make(chan string)
	metricsCh := make(chan metric)
	expected := map[string]float64{"http_200": float64(2), "http_401": float64(2), "http_500": float64(1)}
	parser := new(histogramLogParser)

	// Test
	go parser.parseLogFile(logCh, metricsCh, regexp)
	go func() {
		for _, line := range lines {
			logCh <- line
		}
		close(logCh)
	}()
	for m := range metricsCh {
		if expected[m.name] != m.value {
			t.Error("Expected: ", expected[m.name], "Got: ", m.value)
		}
	}
}
