package backup

//go:generate mockgen -source=${GOFILE} -destination=mocks/${GOFILE} -package=backupmocks

import "github.com/basicrum/front_basicrum_go/types"

// IBackup interface for all backup sub directories
type IBackup interface {
	SaveAsync(event *types.Event)
	SaveUnknown(event *types.Event)
	SaveExpired(event *types.Event)
	Flush()
}

// IBackupSingle is single directory backup interface
type IBackupSingle interface {
	SaveAsync(event *types.Event)
	Flush()
	Compress()
}
