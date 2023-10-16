package backup

import "github.com/basicrum/front_basicrum_go/types"

// IBackup interface for backup service
type IBackup interface {
	SaveAsync(event *types.Event, batcherInstance string)
	SaveUnknown(event *types.Event)
	SaveExpired(event *types.Event)
	Flush()
}
