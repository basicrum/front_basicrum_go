package backup

import (
	"github.com/basicrum/front_basicrum_go/types"
	"github.com/robfig/cron/v3"
)

// FileBackup saves the events on the file system
type FileBackup struct {
	archive IBackupSingle
	expired IBackupSingle
	unknown IBackupSingle
	cron    *cron.Cron
}

// NewFileBackup creates file system backup service
// nolint: revive
func NewFileBackup(
	archive IBackupSingle,
	expired IBackupSingle,
	unknown IBackupSingle,
) (*FileBackup, error) {
	c := cron.New()
	result := &FileBackup{
		archive: archive,
		expired: expired,
		unknown: unknown,
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
	b.expired.Compress()
	b.unknown.Compress()
}

// SaveAsync saves an event with default batcher
func (b *FileBackup) SaveAsync(event *types.Event) {
	b.archive.SaveAsync(event)
}

// SaveExpired saves an expired event asynchronously
func (b *FileBackup) SaveExpired(event *types.Event) {
	b.expired.SaveAsync(event)
}

// SaveUnknown saves an unknown event asynchronously
func (b *FileBackup) SaveUnknown(event *types.Event) {
	b.unknown.SaveAsync(event)
}

// Flush is called before shutdown to force process of the last batch
func (b *FileBackup) Flush() {
	b.archive.Flush()
	b.expired.Flush()
	b.unknown.Flush()
	b.cron.Stop()
}
