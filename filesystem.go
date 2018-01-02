/*
A layer of abstraction around file system IO.

Functions that operate on files, use this abstraction instead of
calling functions from the os package directly.

By default, the functions in this abstraction call the io and os functions
to open files and do disk IO, but the implementation can be substituted by
a mock which is done in unit tests.

This code is inspired by:
https://stackoverflow.com/questions/16742331/how-to-mock-abstract-filesystem-in-go
*/

package main

import (
	"io"
	"os"
)

// All functions in this pakage interact with the file system via this package variable
var fs fileSystem = osFS{}

type fileSystem interface {
	Create(name string) (file, error)
	OpenFile(name string, flag int, perm os.FileMode) (file, error)
	Stat(name string) (os.FileInfo, error)
}

type file interface {
	io.Closer
	io.Reader
	io.Seeker
	io.Writer
}

type osFS struct{}

func (osFS) Open(name string) (file, error) {
	return os.Open(name)
}

func (osFS) Create(name string) (file, error) {
	return os.Create(name)
}

func (osFS) OpenFile(name string, flag int, perm os.FileMode) (file, error) {
	return os.OpenFile(name, flag, perm)
}

func (osFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
