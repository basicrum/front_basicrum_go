package backup

import (
	"io"
)

// Compression the type of compression
type Compression string

const (
	// NoneCompression no compression
	NoneCompression Compression = "NONE"
	// GZIPCompression gzip
	GZIPCompression Compression = "GZIP"
	// ZStandardCompression Zstandard
	ZStandardCompression Compression = "Zstandard"
)

// CompressionLevel the level of compression
type CompressionLevel string

const (
	// NoCompressionLevel no compression
	NoCompressionLevel CompressionLevel = "No"
	// BestSpeedCompressionLevel best speed
	BestSpeedCompressionLevel CompressionLevel = "BestSpeed"
	// DefaultCompressionLevel default compression
	DefaultCompressionLevel CompressionLevel = "Default"
	// BestCompressionCompressionLevel best compression
	BestCompressionCompressionLevel CompressionLevel = "BestCompression"
	// HuffmanOnlyCompressionLevel HuffmanOnly gzip specific compression
	HuffmanOnlyCompressionLevel CompressionLevel = "HuffmanOnly"
)

// CompressionWriterFactory is factory for compression writer
type CompressionWriterFactory interface {
	// Create returns a compression writer
	Create(io.Writer) (io.WriteCloser, error)
	// Filename returns the new filename based on the compression
	Filename(originalFilename string) string
}

// NewCompressionWriterFactory constructor of CompressionWriterFactory
// nolint: revive
func NewCompressionWriterFactory(
	enabled bool,
	compression Compression,
	level CompressionLevel,
) CompressionWriterFactory {
	if !enabled {
		return newNoneFactory()
	}
	switch compression {
	case GZIPCompression:
		return newGZIPFactory(level)
	case ZStandardCompression:
		return newZStdFactory(level)
	case NoneCompression:
		return newNoneFactory()
	default:
		return newNoneFactory()
	}
}
