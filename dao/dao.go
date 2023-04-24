package dao

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"

	"github.com/uptrace/go-clickhouse/ch"
	"github.com/uptrace/go-clickhouse/chdebug"
	"github.com/uptrace/go-clickhouse/chmigrate"
)

const (
	baseTableName          = "webperf_rum_events"
	tablePrefixPlaceholder = "{prefix}"
	migrationsTemplateDir  = "template_migrations"
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

func (p *DAO) copyMigrations(tempDir string) error {
	srcDir := migrationsTemplateDir

	// read files in migrations template directory
	files, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("cannot read files in directory[%v] err[%w]", srcDir, err)
	}

	for _, f := range files {
		if f.IsDir() {
			// copy only files
			// skip directories
			continue
		}
		err = p.processMigrationFile(srcDir, tempDir, f.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DAO) processMigrationFile(srcDir, tempDir, filename string) error {
	// build source and destination file paths
	srcFile := filepath.Join(srcDir, filename)
	dstFile := filepath.Join(tempDir, filename)

	// copy migration file
	_, err := p.copyFile(srcFile, dstFile)
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

func (*DAO) copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
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
