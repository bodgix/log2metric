package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type statefulLogFile struct {
	logFile   file
	stateFile file
}

// Close save the current position and close the file
func (lf *statefulLogFile) Close() error {
	defer lf.logFile.Close()
	defer lf.stateFile.Close()

	// read the current position of the log file and save it to the state file
	pos, err := lf.logFile.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(lf.stateFile, "%d", pos)
	return err
}

func readLogFile(name, stateFile string, outCh chan<- string, errCh chan<- error) {
	logFile, err := openLogFile(name, stateFile)
	defer close(outCh)

	if err != nil {
		errCh <- err
	} else {
		defer logFile.Close()

		var line string
		reader := bufio.NewReader(logFile.logFile)
		for {
			line, err = reader.ReadString('\n')
			if err != nil {
				if err != io.EOF { // report all errors except io.EOF
					errCh <- err
					break
				} else { // EOF reached. Send the last line and stop reading
					outCh <- line
					break
				}
			}
			outCh <- line
		}
	}
}

func openLogFile(name, stateFile string) (*statefulLogFile, error) {
	sfLog := &statefulLogFile{}
	f, err := fs.OpenFile(name, os.O_RDONLY, 0660)
	if err != nil {
		return sfLog, err
	}
	sfLog.logFile = f

	f, err = openStateFile(stateFile, fs)
	if err != nil {
		return sfLog, err
	}
	sfLog.stateFile = f

	var lastPos int64
	lastPos, err = getLastPos(sfLog.stateFile)
	if err != nil {
		return sfLog, err
	}
	sfLog.stateFile.Seek(0, os.SEEK_SET)
	sfLog.logFile.Seek(lastPos, os.SEEK_SET)

	return sfLog, err
}

func openStateFile(name string, fs fileSystem) (file, error) {
	if _, err := fs.Stat(name); os.IsNotExist(err) {
		return fs.Create(name)
	}
	return fs.OpenFile(name, os.O_RDWR, 0660)
}

func getLastPos(stateFile io.Reader) (int64, error) {
	var lastPos int64
	n, err := fmt.Fscanf(stateFile, "%d", &lastPos)
	if n == 0 {
		lastPos = 0
	}
	if err == io.EOF {
		err = nil
	}
	return lastPos, err
}
