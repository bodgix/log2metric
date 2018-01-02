package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	opts := parseOptions()
	if opts.help {
		flag.Usage()
		os.Exit(0)
	}
	if err := validateOptions(opts); err != nil {
		log.Fatal(err)
	}

	logLinesCh := make(chan string)
	errCh := make(chan error)
	metricsCh := make(chan metric)

	defer close(errCh)

	go readLogFile(opts.logFile, opts.stateFile, logLinesCh, errCh)
	go parseLogFile(logLinesCh, metricsCh, opts.regexp)

	fin := false

	for !fin {
		select {
		case m, ok := <-metricsCh:
			if ok {
				log.Println("Received a new metric: ", m)
			} else {
				log.Println("Metrics channel was closed")
				fin = true
			}
		case err := <-errCh:
			log.Println("Received an error: ", err)
			fin = true
		}
	}
}
