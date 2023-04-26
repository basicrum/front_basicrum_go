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
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
	"github.com/uptrace/go-clickhouse/chmigrate"

	"github.com/basicrum/front_basicrum_go/templatemigrations"
)

const (
	baseTableName          = "webperf_rum_events"
	tablePrefixPlaceholder = "{prefix}"
	bufferSize             = 1024
)

// DAO is data access object for clickhouse database
type DAO struct {
	conn               clickhouse.Conn
	table              string
	migrateDatabaseURL string
	prefix             string
}

// New creates persistance service
// nolint: revive
func New(s server, a auth, opts *opts) (*DAO, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{s.addr()},
		Auth: clickhouse.Auth{
			Database: s.db,
			Username: a.user,
			Password: a.pwd,
		},
		Debug:           false,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse connection failed: %w", err)
	}
	migrateDatabaseURL := fmt.Sprintf("clickhouse://%v:%v@%v/%v?sslmode=disable",
		a.user,
		a.pwd,
		s.addr(),
		s.db,
	)
	table := opts.prefix + baseTableName
	return &DAO{
		conn:               conn,
		table:              table,
		migrateDatabaseURL: migrateDatabaseURL,
		prefix:             opts.prefix,
	}, nil
}

// Save stores data into table in clickhouse database
func (p *DAO) Save(data string) error {
	if data == "" {
		return fmt.Errorf("clickhouse invalid data for table %s: %s", p.table, data)
	}
	query := fmt.Sprintf(
		"INSERT INTO %s SETTINGS input_format_skip_unknown_fields = true FORMAT JSONEachRow %s",
		p.table,
		data,
	)
	err := p.conn.AsyncInsert(context.Background(), query, false)
	if err != nil {
		return fmt.Errorf("clickhouse insert failed: %w", err)
	}
	return nil
}

// Migrate applies all pending database migrations
func (p *DAO) Migrate() error {
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
func (p *DAO) copyMigrations(tempDir string) error {
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

func (p *DAO) processMigrationFile(srcFile fs.File, tempDir, filename string) error {
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
func (*DAO) copyFile(src fs.File, dst string) error {
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

func (*DAO) replaceTextInFile(file, find, replace string) error {
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

func (p *DAO) migrateUp(sourcePath string) error {
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

	migrator := chmigrate.NewMigrator(db, migrations)

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
