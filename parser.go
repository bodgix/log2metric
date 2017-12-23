package main

import (
	"log"
	"regexp"
	"strconv"
)

type logParser interface {
	parseLogFile(lf logFile, instanceName, regExp string) (map[string][]float32, error)
}

type simpleLogParser struct{}

func (lp simpleLogParser) parseLogFile(lf logFile, instanceName, regExp string) (map[string][]float32, error) {
	metrics := make(map[string][]float32)
	exp := regexp.MustCompile(regExp)
	for line, err := lf.nextLine(); err != errNoMoreLines; line, err = lf.nextLine() {
		if err != nil {
			return metrics, err
		}
		log.Print("Parsing line: ", line)
		matches := exp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		for i, name := range exp.SubexpNames() {
			if name == "" {
				continue
			}
			log.Printf("Named capture number %d, %v", i, name)
			val, err := strconv.ParseFloat(matches[i], 32)
			if err != nil {
				continue
			}
			metrics[instanceName+name] = append(metrics[instanceName+name], float32(val))
		}
	}
	return metrics, nil
}

func buildParser(t string) (logParser, error) {
	var lp simpleLogParser
	return lp, nil
}
