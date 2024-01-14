package backup

import (
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/robfig/cron/v3"
)

// FileBackup saves the events on the file system
type FileBackup struct {
	archive IBackupSingle
	cron    *cron.Cron
}

// NewFileBackup creates file system backup service
// nolint: revive
func NewFileBackup(
	archive IBackupSingle,
) (*FileBackup, error) {
	c := cron.New()
	result := &FileBackup{
		archive: archive,
		cron:    c,
	}
	// 01:00:00 each day
	_, err := c.AddFunc("CRON_TZ=UTC  0 1 * * *", result.compressDay)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *FileBackup) compressDay() {
	b.archive.Compress()
}

// SaveAsync saves an event with default batcher
func (b *FileBackup) SaveAsync(event *types.Event) {
	b.archive.SaveAsync(event)
}

// Flush is called before shutdown to force process of the last batch
func (b *FileBackup) Flush() {
	b.archive.Flush()
	b.cron.Stop()
}
