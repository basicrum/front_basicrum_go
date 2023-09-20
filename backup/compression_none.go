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
func (noneFactory) Create(w io.WriteCloser) (io.WriteCloser, error) {
	return w, nil
}
