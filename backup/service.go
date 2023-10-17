package backup

import (
	"encoding/json"
	"log"
	"time"

	"github.com/basicrum/front_basicrum_go/types"
	"github.com/eapache/go-resiliency/batcher"
	"github.com/robfig/cron/v3"
)

// BATCHER_DEFAULT is the default batcher for general archiving
const BATCHER_DEFAULT = "batcher_archive"

// BATCHER_EXPIRED is the batcher for expired subscription events
const BATCHER_EXPIRED = "batcher_expired"

// BATCHER_UNKNOWN is the batcher for unknown subscription events
const BATCHER_UNKNOWN = "batcher_unknown"

// FileBackup saves the events on the file system
type FileBackup struct {
	batcherBackup      *batcher.Batcher
	batcherExpired     *batcher.Batcher
	batcherUnknown     *batcher.Batcher
	cron               *cron.Cron
	directory          string
	expiredDirectory   string
	unknownDirectory   string
	compressionFactory CompressionWriterFactory
}

// NewFileBackup creates file system backup service
// nolint: revive
func NewFileBackup(
	backupInterval time.Duration,
	directory string,
	expiredDirectory string,
	unknownDirectory string,
	compressionFactory CompressionWriterFactory,
) (*FileBackup, error) {
	batcherBackup := batcher.New(backupInterval, func(params []any) error {
		do(params, directory)
		return nil
	})
	batcherExpired := batcher.New(backupInterval, func(params []any) error {
		do(params, expiredDirectory)
		return nil
	})
	batcherUnknown := batcher.New(backupInterval, func(params []any) error {
		do(params, unknownDirectory)
		return nil
	})
	c := cron.New()
	result := &FileBackup{
		batcherBackup:      batcherBackup,
		batcherExpired:     batcherExpired,
		batcherUnknown:     batcherUnknown,
		cron:               c,
		directory:          directory,
		expiredDirectory:   expiredDirectory,
		unknownDirectory:   unknownDirectory,
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

// SaveAsync saves an event with default batcher
func (b *FileBackup) SaveAsync(event *types.Event) {
	b.doSaveAsync(event, BATCHER_DEFAULT)
}

// nolint: revive
func (b *FileBackup) doSaveAsync(event *types.Event, batcherInstance string) {
	go func() {
		forArchiving := event.RequestParameters
		// Flatten headers later
		h, hErr := json.Marshal(forArchiving)
		if hErr != nil {
			log.Println(hErr)
		}
		forArchiving.Add("request_headers", string(h))
		switch batcherInstance {
		case BATCHER_EXPIRED:
			if err := b.batcherExpired.Run(forArchiving); err != nil {
				log.Printf("Error archiving expired url[%v] err[%v]", forArchiving, err)
			}
		case BATCHER_UNKNOWN:
			if err := b.batcherUnknown.Run(forArchiving); err != nil {
				log.Printf("Error archiving unknown url[%v] err[%v]", forArchiving, err)
			}
		case BATCHER_DEFAULT:
			if err := b.batcherUnknown.Run(forArchiving); err != nil {
				log.Printf("Error archiving unknown url[%v] err[%v]", forArchiving, err)
			}
		default:
			log.Printf("Error batcher instance unknown: %s", batcherInstance)
		}
	}()
}

// SaveExpired saves an expired event asynchronously
func (b *FileBackup) SaveExpired(event *types.Event) {
	b.doSaveAsync(event, BATCHER_EXPIRED)
}

// SaveUnknown saves an unknown event asynchronously
func (b *FileBackup) SaveUnknown(event *types.Event) {
	b.doSaveAsync(event, BATCHER_UNKNOWN)
}

// Flush is called before shutdown to force process of the last batch
func (b *FileBackup) Flush() {
	b.batcherBackup.Shutdown(true)
	b.batcherExpired.Shutdown(true)
	b.batcherUnknown.Shutdown(true)
	b.cron.Stop()
}
