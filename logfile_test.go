package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

var logFileLines = []string{"line1", "line2", "line3"}
var stateFileLines = []string{"0"}
var openLogError error
var openStateError error

type testFile struct {
	*bytes.Reader
	*bytes.Buffer
}

func (testFile) Close() error {
	return nil
}

func (f testFile) Write(p []byte) (int, error) {
	log.Print("Writing ", p, " to the testFile.")
	return f.Buffer.Write(p)
}

func (f testFile) Read(p []byte) (int, error) {
	return f.Reader.Read(p)
}

type mockFS struct{}

func (mockFS) Create(name string) (file, error) {
	return new(testFile), nil
}

func (fs mockFS) OpenFile(name string, flag int, perm os.FileMode) (file, error) {
	var lines []string
	var err error
	if name == "log" {
		lines = logFileLines
		err = openLogError
	} else {
		lines = stateFileLines
		err = openStateError
	}
	buf := testFile{bytes.NewReader([]byte(strings.Join(lines, "\n"))), &bytes.Buffer{}}
	return buf, err
}

func (mockFS) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

func TestReadLogFileNoErrors(t *testing.T) {
	fs = mockFS{}
	logLinesCh := make(chan string)
	errCh := make(chan error)
	go readLogFile("log", "state", logLinesCh, errCh)
	i := 0
	for line := range logLinesCh {
		line = strings.TrimRight(line, "\n")
		if line != logFileLines[i] {
			t.Errorf("Expected: %x got %x\n", logFileLines[i], line)
		}
		i++
	}
}

func TestReadLogFileOpenLogFileError(t *testing.T) {
	fs = mockFS{}
	logLinesCh := make(chan string)
	errCh := make(chan error)
	openLogError = errors.New("Error opening the log file")
	go readLogFile("log", "state", logLinesCh, errCh)
	err := <-errCh
	if err != openLogError {
		t.Errorf("Expected error %v, got %v\n", openLogError, err)
	}
}

func TestReadLogFileOpenStateFileError(t *testing.T) {
	fs = mockFS{}
	logLinesCh := make(chan string)
	errCh := make(chan error)
	openStateError = errors.New("Error opening the state file")
	go readLogFile("log", "state", logLinesCh, errCh)
	err := <-errCh
	if err != openLogError {
		t.Errorf("Expected error %v, got %v\n", openLogError, err)
	}
}

func TestCloseStatefulLogFile(t *testing.T) {
	logFile := new(statefulLogFile)
	testFile := &testFile{bytes.NewReader([]byte(strings.Join(logFileLines, "\n"))), new(bytes.Buffer)}
	logFile.logFile = testFile
	logFile.stateFile = testFile
	logFile.Close()
	got := testFile.String()
	if got != "0" {
		t.Error("Expected 0 written to the state file, got ", got)
	}
}
