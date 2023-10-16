package backup

import (
	"time"
)

// New is factory for backup service
// nolint: revive
func New(
	enabled bool,
	backupInterval time.Duration,
	directory string,
	expiredDirectory string,
	unknownDirectory string,
	compressionFactory CompressionWriterFactory,
) (IBackup, error) {
	if !enabled {
		return NewNullBackup(), nil
	}
	return NewFileBackup(backupInterval, directory, expiredDirectory, unknownDirectory, compressionFactory)
}
