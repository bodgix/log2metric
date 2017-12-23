package main

import (
	"fmt"
	"log"
)

func main() {
	lf, err := openLogFile("apache.log", "/tmp/apache_log_state", fs)
	if err != nil {
		log.Fatal("Cannot open logfile ", err)
	}
	parser, err := buildParser("simple")
	if err != nil {
		log.Fatal(err)
	}
	metrics, err := parser.parseLogFile(lf, "test", `(?P<resp_time>[\d.]+)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", metrics)
	err = lf.Close()
	if err != nil {
		log.Fatal(err)
	}
}
