package main

import (
	"errors"
	"regexp"
	"strconv"
)

type logParser interface {
	parseLogFile(input <-chan string, output chan<- metric, regExp string)
	validateOptions(opts options) error
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
	duration
)

type metric struct {
	t     metricType
	name  string
	value float64
}

func parseLines(lines <-chan string, matches chan<- []string, regExp string) {
	defer close(matches)

	exp := regexp.MustCompile(regExp)
	for line := range lines {
		ms := exp.FindStringSubmatch(line)
		if ms == nil {
			continue
		}
		for i, name := range exp.SubexpNames() {
			if name == "" {
				continue
			}
			matches <- []string{name, ms[i]}
		}
	}
}

type simpleLogParser struct{}

func (simpleLogParser) parseLogFile(input <-chan string, output chan<- metric, regExp string) {
	defer close(output)

	matches := make(chan []string)

	go parseLines(input, matches, regExp)
	for match := range matches {
		val, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			continue
		}
		output <- metric{t: simple, name: match[0], value: val}
	}
}

func (simpleLogParser) validateOptions(opts options) error {
	return nil
}

type histogramLogParser struct{}

func (histogramLogParser) parseLogFile(input <-chan string, output chan<- metric, regExp string) {
	defer close(output)

	histogram := make(map[string]int)
	matches := make(chan []string)

	go parseLines(input, matches, regExp)

	for match := range matches {
		histogram[match[0]+"_"+match[1]]++
	}
	for name, val := range histogram {
		output <- metric{t: histo, name: name, value: float64(val)}
	}
}

func (histogramLogParser) validateOptions(opts options) error {
	return nil
}

type eventLogParser struct{}

func (eventLogParser) parseLogFile(input <-chan string, output chan<- metric, regExp string) {
}

func any(col []string, what string) bool {
	for _, elem := range col {
		if elem == what {
			return true
		}
	}
	return false
}

func (eventLogParser) validateOptions(opts options) error {
	exp, err := regexp.Compile(opts.regexp)
	if err != nil {
		return err
	}
	if !any(exp.SubexpNames(), "event_id") {
		return errors.New("event_id named group must exist in the regexp")
	}
	if !any(exp.SubexpNames(), "ts") {
		return errors.New("ts named group must exist in the regexp")
	}
	return nil
}
