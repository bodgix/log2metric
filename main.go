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
	doneCh := make(chan bool)

	defer close(errCh)

	go readLogFile(opts.logFile, opts.stateFile, logLinesCh, errCh)
	go parseLogFile(logLinesCh, metricsCh, opts.regexp)
	go printer.printer(metricsCh, doneCh)

	done := false

	for !done {
		select {
		case <-doneCh:
			done = true
		case err := <-errCh:
			log.Println("Received an error: ", err)
			done = true
		}
	}
}
