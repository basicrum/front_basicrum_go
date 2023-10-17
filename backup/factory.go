package backup

import (
	"fmt"
	"os"
	"path"
	"time"
)

type backupType string

const (
	// archive is the default batcher for general request archiving
	archive backupType = "archive"
	// expired is the batcher for expired subscription events
	expired backupType = "expired"
	// unknown is the batcher for unknown subscription events
	unknown backupType = "unknown"
)

// New is factory for backup service
// nolint: revive
func New(
	enabled bool,
	backupInterval time.Duration,
	baseDirectory string,
	compressionFactory CompressionWriterFactory,
) (IBackup, error) {
	if !enabled {
		return NewNullBackup(), nil
	}
	archiveBackup, err := makeSingle(archive, backupInterval, baseDirectory, compressionFactory)
	if err != nil {
		return nil, err
	}
	expiredBackup, err := makeSingle(expired, backupInterval, baseDirectory, compressionFactory)
	if err != nil {
		return nil, err
	}
	unknownBackup, err := makeSingle(unknown, backupInterval, baseDirectory, compressionFactory)
	if err != nil {
		return nil, err
	}
	return NewFileBackup(archiveBackup, expiredBackup, unknownBackup)
}

func makeSingle(singleBackupType backupType,
	backupInterval time.Duration,
	baseDirectory string,
	compressionFactory CompressionWriterFactory,
) (IBackupSingle, error) {
	directory := path.Join(baseDirectory, string(singleBackupType))
	if err := os.MkdirAll(directory, os.ModeDir.Perm()); err != nil {
		return nil, fmt.Errorf("cannot create directory[%v], %w", directory, err)
	}
	return NewSingleFileBackup(backupInterval, directory, compressionFactory), nil
}
