package main

import (
	"log"
	"regexp"
	"strconv"
)

type metricType int

const (
	simple metricType = iota
)

type metric struct {
	t     metricType
	name  string
	value float64
}

func parseLogFile(input <-chan string, output chan<- metric, regExp string) {
	defer close(output)
	exp := regexp.MustCompile(regExp)
	for line := range input {
		log.Print("Parsing line: ", line)
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
