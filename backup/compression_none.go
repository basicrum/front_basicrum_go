package backup

import (
	"io"
)

type noneFactory struct {
}

func newNoneFactory() noneFactory {
	return noneFactory{}
}

// Create returns a compression writer
func (noneFactory) Create(w io.Writer) (io.WriteCloser, error) {
	return writerWrapper{w}, nil
}

// Filename returns the new filename based on the compression
func (noneFactory) Filename(originalFilename string) string {
	return originalFilename
}

type writerWrapper struct {
	io.Writer
}

// Close implements the closer interface
func (writerWrapper) Close() error {
	return nil
}
