package backup

import "github.com/basicrum/front_basicrum_go/types"

// NullBackup is disabled backup implementation
type NullBackup struct {
}

// NewNullBackup creates disabled backup implementation
func NewNullBackup() *NullBackup {
	return &NullBackup{}
}

// SaveAsync disabled implementation
func (*NullBackup) SaveAsync(_ *types.Event, batcherInstance string) {}

// SaveExpired disabled implementation
func (*NullBackup) SaveExpired(_ *types.Event) {}

// SaveUnknown disabled implementation
func (*NullBackup) SaveUnknown(_ *types.Event) {}

// Flush disabled implementation
func (*NullBackup) Flush() {}
