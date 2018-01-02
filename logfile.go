package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type statefulLogFile struct {
	logFile   file
	logReader *bufio.Reader
	stateFile file
}

// Close save the current position and close the file
func (lf *statefulLogFile) Close() error {
	log.Print("Closing the log file")
	defer lf.logFile.Close()
	defer lf.stateFile.Close()

	pos, err := lf.logFile.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}
	log.Printf("The current position in the log file is: %d", pos)

	_, err = fmt.Fprintf(lf.stateFile, "%d", pos)
	return err
}

func readLogFile(name, stateFile string, outCh chan<- string, errCh chan<- error) {
	logFile, err := openLogFile(name, stateFile)
	defer close(outCh)

	if err != nil {
		log.Println("Error opening the log file ", err)
		errCh <- err
	} else {
		defer logFile.Close()

		var line string
		for {
			line, err = logFile.logReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- err
					break
				} else {
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
	log.Print("Opened the log file ", f)
	sfLog.logReader = bufio.NewReader(sfLog.logFile)

	f, err = openStateFile(stateFile, fs)
	if err != nil {
		log.Print("Error openning the state file ", err)
		return sfLog, err
	}
	log.Print("Opened the state file ", f)
	sfLog.stateFile = f

	var lastPos int64
	lastPos, err = getLastPos(sfLog.stateFile)
	if err != nil {
		log.Print("Error getting the last position")
		return sfLog, err
	}
	sfLog.stateFile.Seek(0, os.SEEK_SET)

	log.Print("Seeking the log file to position ", lastPos)
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
