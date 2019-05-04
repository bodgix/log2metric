package main

import (
	"errors"
	"flag"
	"regexp"
	"strconv"
	"time"
)

type logParser interface {
	parseLogFile(input <-chan string, output chan<- metric, opts options)
	validateOptions(opts options) error
}

func init() {
	// Add parser-related cmd-line options
	flag.StringVar(&opts.regexp, "regexp", "", "regexp with named captures")
	flag.StringVar(&opts.endRegexp, "endregexp", "", "regexp matching the end of an event - only in duration mode")
	flag.StringVar(&opts.tsFormat, "tsformat", "", "timestamp layout in golang's time.Parse format - only needed in duration mode")
	flag.StringVar(&opts.durationCacheFile, "durationcachefile", "", "cache file to save unmatched events in duration mode")

	flag.BoolVar(&opts.histogram, "histogram", false, "run in the histogram mode")
	flag.BoolVar(&opts.duration, "duration", false, "run in the duration mode")
}

func getLogParser(opts options) logParser {
	if opts.histogram {
		return new(histogramLogParser)
	} else if opts.duration {
		return new(eventLogParser)
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

func (simpleLogParser) parseLogFile(input <-chan string, output chan<- metric, opts options) {
	defer close(output)

	matches := make(chan []string)

	go parseLines(input, matches, opts.regexp)
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

func (histogramLogParser) parseLogFile(input <-chan string, output chan<- metric, opts options) {
	defer close(output)

	histogram := make(map[string]int)
	matches := make(chan []string)

	go parseLines(input, matches, opts.regexp)

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

func (eventLogParser) parseLogFile(input <-chan string, output chan<- metric, opts options) {
	defer close(output)

	events := make(map[string]time.Time)
	matches := make(chan []string)
	go parseLines(input, matches, opts.regexp)

	for {
		match, next := <-matches
    switch(match[0]) {
      case "ts"
    }
	}
}

func any(col []string, what string) bool {
	for _, elem := range col {
		if elem == what {
			return true
		}
	}
	return false
}

func validateDurationRegexp(regExp string) error {
	exp, err := regexp.Compile(regExp)
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

func (eventLogParser) validateOptions(opts options) error {
	if err := validateDurationRegexp(opts.regexp); err != nil {
		return err
	}

	if opts.endRegexp == "" {
		return errors.New("endregexp is required in duration mode")
	}
	if opts.tsFormat == "" {
		return errors.New("tsformat is required in duration mode")
	}
	if err := validateDurationRegexp(opts.endRegexp); err != nil {
		return err
	}
	return nil
}
