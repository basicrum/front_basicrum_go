package backup

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

type zstdFactory struct {
	level CompressionLevel
}

func newZStdFactory(level CompressionLevel) zstdFactory {
	return zstdFactory{
		level: level,
	}
}

// Create returns a compression writer
func (f zstdFactory) Create(w io.Writer) (io.WriteCloser, error) {
	return zstd.NewWriter(w, zstd.WithEncoderLevel(f.makeLevel()))
}

func (f zstdFactory) makeLevel() zstd.EncoderLevel {
	switch f.level {
	case NoCompressionLevel:
		return zstd.SpeedFastest
	case BestSpeedCompressionLevel:
		return zstd.SpeedBetterCompression
	case BestCompressionCompressionLevel:
		return zstd.SpeedBestCompression
	case HuffmanOnlyCompressionLevel:
		return zstd.SpeedDefault
	case DefaultCompressionLevel:
		return zstd.SpeedDefault
	default:
		return zstd.SpeedDefault
	}
}

// Filename returns the new filename based on the compression
func (zstdFactory) Filename(originalFilename string) string {
	return originalFilename + ".zst"
}
