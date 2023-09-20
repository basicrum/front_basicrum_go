package backup

import (
	"compress/gzip"
	"io"
)

type gzipFactory struct {
	level CompressionLevel
}

func newGZIPFactory(level CompressionLevel) gzipFactory {
	return gzipFactory{
		level: level,
	}
}

// Create returns a compression writer
func (f gzipFactory) Create(fw io.WriteCloser) (io.WriteCloser, error) {
	return gzip.NewWriterLevel(fw, f.makeLevel())
}

func (f gzipFactory) makeLevel() int {
	switch f.level {
	case NoCompressionLevel:
		return gzip.NoCompression
	case HuffmanOnlyCompressionLevel:
		return gzip.HuffmanOnly
	case BestSpeedCompressionLevel:
		return gzip.BestSpeed
	case BestCompressionCompressionLevel:
		return gzip.BestCompression
	case DefaultCompressionLevel:
		return gzip.DefaultCompression
	default:
		return gzip.DefaultCompression
	}
}
