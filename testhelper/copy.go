package testhelper

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyFile create a destination file and copy the source file content
func CopyFile(source string, dest string) error {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err != nil {
		return err
	}

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.Chmod(dest, sourceinfo.Mode())
}

// CopyDir creates destination directory and copies all the files and sub directories from source directory
func CopyDir(source string, dest string) error {
	err := createDestinationDirectory(source, dest)
	if err != nil {
		return err
	}

	directory, err := os.Open(source)
	if err != nil {
		return err
	}
	defer directory.Close()

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		err := copyDirEntry(source, dest, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyDirEntry(source string, dest string, obj fs.FileInfo) error {
	sourcefilepointer := filepath.Join(source, obj.Name())
	destinationfilepointer := filepath.Join(dest, obj.Name())
	if obj.IsDir() {
		// create sub-directories - recursively
		err := CopyDir(sourcefilepointer, destinationfilepointer)
		if err != nil {
			return err
		}
	} else {
		// perform copy
		err := CopyFile(sourcefilepointer, destinationfilepointer)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDestinationDirectory(source string, dest string) error {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}
	return nil
}
