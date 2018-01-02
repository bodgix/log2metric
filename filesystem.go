package main

import (
	"io"
	"os"
)

// This file contains a mockable filesystem abstraction that all other files
// use for all their disk IO needs.
// By default, the File system is backed by real golang's io package
// and all IO goes to disk for real.
// However, the design makes it possible to mock the filesystem operations
// which can be useful for testing.
// This code is inspired by:
// https://stackoverflow.com/questions/16742331/how-to-mock-abstract-filesystem-in-go

// The default filesystem is the Real OS fs so real IO is made by default.
// For testing, it's possible to use a mock file system though.
// This code is inspired by https://stackoverflow.com/questions/16742331/how-to-mock-abstract-filesystem-in-go
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
