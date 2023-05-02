package templatemigrations

import (
	"embed"
)

// SQLMigrations contains the embed template sql migrations
// nolint: gochecknoglobals
//
//go:embed *.sql
var SQLMigrations embed.FS
