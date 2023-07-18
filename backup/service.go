package backup

import (
	"encoding/json"
	"log"
	"time"

	"github.com/basicrum/front_basicrum_go/types"
	"github.com/eapache/go-resiliency/batcher"
)

// FileBackup saves the events on the file system
type FileBackup struct {
	batcher *batcher.Batcher
}

// NewFileBackup creates file system backup service
func NewFileBackup(
	backupInterval time.Duration,
	directory string,
) *FileBackup {
	b := batcher.New(backupInterval, func(params []any) error {
		do(params, directory)
		return nil
	})
	return &FileBackup{
		batcher: b,
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
}
