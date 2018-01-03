package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type options struct {
	regexp    string
	logFile   string
	stateFile string
	prefix    string
	histogram bool
	help      bool
}

func parseOptions() options {
	var opts options
	flag.StringVar(&opts.regexp, "regexp", "", "regexp with named captures")
	flag.StringVar(&opts.logFile, "logfile", "", "full path to the log file")
	flag.StringVar(&opts.stateFile, "statefile", "", "full path to the state file")
	flag.StringVar(&opts.prefix, "prefix", "", "prefix to add to metrics names")
	flag.BoolVar(&opts.histogram, "histogram", false, "run in the histogram mode")
	flag.BoolVar(&opts.help, "help", false, "print this help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	return opts
}

func validateOptions(opts options) error {
	switch "" {
	case opts.logFile:
		return errors.New("logfile is required")
	case opts.stateFile:
		return errors.New("statefile is required")
	case opts.regexp:
		return errors.New("regexp is required")
	}
	return validateRegexp(opts.regexp)
}

func validateRegexp(expr string) error {
	exp, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	namedCaptures := filter(exp.SubexpNames(), func(s string) bool {
		return s != ""
	})
	if len(namedCaptures) < 1 {
		return errors.New("regexp must have named cupture groups")
	}
	return nil
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
