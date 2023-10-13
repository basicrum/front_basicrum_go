package backup

import "github.com/basicrum/front_basicrum_go/types"

// IBackup interface for backup service
type IBackup interface {
	SaveAsync(event *types.Event)
	Flush()
	// TODO:
	// SaveUnknown(event *types.Event)
	// SaveExpired(event *types.Event)
}
