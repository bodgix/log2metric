package main

import (
	"regexp"
	"strconv"
)

type logParser interface {
	parseLogFile(input <-chan string, output chan<- metric, regExp string)
}

func getLogParser(opts options) logParser {
	if opts.histogram {
		return new(histogramLogParser)
	}
	return new(simpleLogParser)
}

type metricType int

const (
	simple metricType = iota
	histo
)

type metric struct {
	t     metricType
	name  string
	value float64
}

type simpleLogParser struct{}

func (simpleLogParser) parseLogFile(input <-chan string, output chan<- metric, regExp string) {
	defer close(output)
	exp := regexp.MustCompile(regExp)
	for line := range input {
		matches := exp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		for i, name := range exp.SubexpNames() {
			if name == "" {
				continue
			}
			val, err := strconv.ParseFloat(matches[i], 64)
			if err != nil {
				continue
			}
			m := metric{t: simple, name: name, value: val}
			output <- m
		}
	}
}

type histogramLogParser struct{}

func (histogramLogParser) parseLogFile(input <-chan string, output chan<- metric, regExp string) {
	histogram := make(map[string]int)
	defer close(output)

	exp := regexp.MustCompile(regExp)
	for line := range input {
		matches := exp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		for i, name := range exp.SubexpNames() {
			if name == "" {
				continue
			}
			histogram[name+"_"+matches[i]]++
		}
	}
	for name, val := range histogram {
		output <- metric{t: histo, name: name, value: float64(val)}
	}
}
