package backup

import (
	"encoding/json"
	"log"
	"time"

	"github.com/basicrum/front_basicrum_go/types"
	"github.com/eapache/go-resiliency/batcher"
	"github.com/robfig/cron/v3"
)

// FileBackup saves the events on the file system
type FileBackup struct {
	batcher            *batcher.Batcher
	cron               *cron.Cron
	directory          string
	compressionFactory CompressionWriterFactory
}

// NewFileBackup creates file system backup service
func NewFileBackup(
	backupInterval time.Duration,
	directory string,
	compressionFactory CompressionWriterFactory,
) (*FileBackup, error) {
	b := batcher.New(backupInterval, func(params []any) error {
		do(params, directory)
		return nil
	})
	c := cron.New()
	result := &FileBackup{
		batcher:            b,
		cron:               c,
		directory:          directory,
		compressionFactory: compressionFactory,
	}
	// 01:00:00 each day
	_, err := c.AddFunc("CRON_TZ=UTC  0 1 * * *", result.compressDay)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *FileBackup) compressDay() {
	day := time.Now().UTC().AddDate(0, 0, -1)
	if err := archiveDay(b.directory, day, b.compressionFactory); err != nil {
		log.Printf("error archive day[%v] err[%v]", day, err)
	}
}

// SaveAsync saves an event asynchronously
func (b *FileBackup) SaveAsync(event *types.Event) {
	go func() {
		forArchiving := event.RequestParameters
		// Flatten headers later
		h, hErr := json.Marshal(forArchiving)
		if hErr != nil {
			log.Println(hErr)
		}
		forArchiving.Add("request_headers", string(h))
		if err := b.batcher.Run(forArchiving); err != nil {
			log.Printf("Error archiving url[%v] err[%v]", forArchiving, err)
		}
	}()
}

// Flush is called before shutdown to force process of the last batch
func (b *FileBackup) Flush() {
	b.batcher.Shutdown(true)
	b.cron.Stop()
}
