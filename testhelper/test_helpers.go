package testhelper

import (
	"path"
	"runtime"
)

// GetProjectRoot returns the root directory of the project
func GetProjectRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get current filename")
	}
	return path.Join(path.Dir(filename), "..")
}
