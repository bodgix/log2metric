package main

import (
	"log"
)

func main() {
	logLinesCh := make(chan string)
	errCh := make(chan error)
	metricsCh := make(chan metric)

	defer close(errCh)

	go readLogFile("apache.log", "/tmp/apache_log_state", logLinesCh, errCh)
	go parseLogFile(logLinesCh, metricsCh, `(?P<resp_time>[\d.]+)`)

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
