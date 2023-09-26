package testhelper

import (
	"compress/gzip"
	"io"
	"strings"

	"github.com/klauspost/compress/zstd"
)

func newReadCloser(file string, reader io.Reader) (io.ReadCloser, error) {
	if strings.HasSuffix(file, ".gz") {
		return gzip.NewReader(reader)
	}
	if strings.HasSuffix(file, ".zst") {
		result, err := zstd.NewReader(reader)
		if err != nil {
			return nil, err
		}
		return readerWrapper{result}, nil
	}
	return readerWrapper{reader}, nil
}

type readerWrapper struct {
	io.Reader
}

// Close implements the closer interface
func (readerWrapper) Close() error {
	return nil
}
