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
) IBackup {
	if !enabled {
		return NewNullBackup()
	}
	return NewFileBackup(backupInterval, directory)
}
