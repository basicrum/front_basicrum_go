package dao

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
	"github.com/uptrace/go-clickhouse/chmigrate"

	"github.com/basicrum/front_basicrum_go/templatemigrations"
)

// MigrationDAO is data access object for clickhouse database
type MigrationDAO struct {
	migrateDatabaseURL string
	prefix             string
}

// New creates persistance service
// nolint: revive
func NewMigrationDAO(s server, a auth, opts *opts) *MigrationDAO {
	return &MigrationDAO{
		migrateDatabaseURL: migrateDBURL(s, a),
		prefix:             opts.prefix,
	}
}

func migrateDBURL(s server, a auth) string {
	return fmt.Sprintf("clickhouse://%v:%v@%v/%v?sslmode=disable",
		a.user,
		a.pwd,
		s.addr(),
		s.db,
	)
}

// Migrate applies all pending database migrations
func (p *MigrationDAO) Migrate() error {
	tempDir, err := os.MkdirTemp("", "migrations")
	if err != nil {
		return fmt.Errorf("cannot create temp directory migrations err[%w]", err)
	}
	defer os.RemoveAll(tempDir)
	err = p.copyMigrations(tempDir)
	if err != nil {
		return fmt.Errorf("cannot copy migrations err[%w]", err)
	}
	return p.migrateUp(tempDir)
}

// nolint: revive
func (p *MigrationDAO) copyMigrations(tempDir string) error {
	return fs.WalkDir(templatemigrations.SQLMigrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		f, err := templatemigrations.SQLMigrations.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("received close file[%v] err[%v]\n", path, err)
			}
		}()

		return p.processMigrationFile(f, tempDir, path)
	})
}

func (p *MigrationDAO) processMigrationFile(srcFile fs.File, tempDir, filename string) error {
	// build source and destination file paths
	dstFile := filepath.Join(tempDir, filename)

	// copy migration file
	err := p.copyFile(srcFile, dstFile)
	if err != nil {
		return fmt.Errorf("cannot copy file[%v] into temp directory[%v] err[%w]", srcFile, dstFile, err)
	}

	// replace table prefix in file
	err = p.replaceTextInFile(dstFile, tablePrefixPlaceholder, p.prefix)
	if err != nil {
		return fmt.Errorf("cannot replace table prefix in migration file[%v] err[%w]", dstFile, err)
	}

	return nil
}

// nolint: revive
func (*MigrationDAO) copyFile(src fs.File, dst string) error {
	sourceFileStat, err := src.Stat()
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, bufferSize)
	for {
		n, err := src.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

func (*MigrationDAO) replaceTextInFile(file, find, replace string) error {
	// read file permissions
	stat, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("cannot read file stat file[%v] err[%w]", file, err)
	}

	// read file
	input, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read source migration file[%v] err[%w]", file, err)
	}

	// replace text in file
	output := bytes.ReplaceAll(input, []byte(find), []byte(replace))

	// save file
	if err = os.WriteFile(file, output, stat.Mode()); err != nil {
		return fmt.Errorf("cannot write replaced migration file[%v] err[%w]", file, err)
	}
	return nil
}

func (p *MigrationDAO) migrateUp(sourcePath string) error {
	db := ch.Connect(ch.WithDSN(p.migrateDatabaseURL))

	db.AddQueryHook(chdebug.NewQueryHook(
		chdebug.WithEnabled(false),
		chdebug.FromEnv("CHDEBUG"),
	))

	sourceFS := os.DirFS(sourcePath)
	var migrations = chmigrate.NewMigrations()
	if err := migrations.Discover(sourceFS); err != nil {
		return fmt.Errorf("cannot discover migrations path[%v] err[%w]", sourcePath, err)
	}

	migrator := chmigrate.NewMigrator(db, migrations,
		chmigrate.WithTableName(p.prefix+"ch_migrations"),
		chmigrate.WithLocksTableName(p.prefix+"ch_migration_locks"),
	)

	ctx := context.Background()

	// create ch_migrations (changelog) and ch_migration_locks tables
	err := migrator.Init(ctx)
	if err != nil {
		return err
	}

	// lock the migrations
	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	// unlock the migrations
	defer func() {
		if err := migrator.Unlock(ctx); err != nil {
			log.Printf("received unlock err[%v]\n", err)
		}
	}()

	// apply the migrations
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}

	if group.IsZero() {
		log.Printf("there are no new migrations to run (database is up to date)\n")
		return nil
	}
	log.Printf("migrated to %s\n", group)

	return nil
}
