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
func (f zstdFactory) Create(w io.WriteCloser) (io.WriteCloser, error) {
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
