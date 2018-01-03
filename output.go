package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

var output = nagiosPrinter{}

type nagiosPrinter struct{}

func (nagiosPrinter) printer(metricsCh <-chan metric, doneCh chan<- func()) {
	done := false
	metrics := make(map[string][]float64)
	for !done {
		m, more := <-metricsCh
		if !more {
			break
		}
		metrics[m.name] = append(metrics[m.name], m.value)
	}
	fmt.Println("OK |", perfDataString(metrics))
	exitCode := 0
	doneCh <- func() { os.Exit(exitCode) }
}

func perfDataString(metrics map[string][]float64) string {
	var result bytes.Buffer
	for name, vals := range metrics {
		result.WriteString(name)
		result.WriteString("=")
		result.WriteString(strconv.FormatFloat(avg(vals), 'f', -1, 64))
		result.WriteString(";")
	}
	return result.String()
}

func avg(vals []float64) float64 {
	var sum float64
	var i float64

	for _, val := range vals {
		sum += val
		i++
	}

	if i > 0 {
		return sum / i
	}
	return 0
}
