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

	// Test
	go parseLogFile(logCh, metricsCh, regexp)
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
