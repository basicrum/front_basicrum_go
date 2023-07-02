package dao

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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

	"github.com/basicrum/front_basicrum_go/beacon"
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
func (p *DAO) Save(theEvent beacon.RumEvent) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (
			created_at,
			hostname,
			event_type,
			browser_name,
			browser_version,
			ua_vnd,
			ua_plt,
			device_type,
			device_manufacturer,
			operating_system,
			operating_system_version,
			user_agent,
			next_hop_protocol,
			visibility_state,
			session_id,
			session_length,
			url,
			connect_duration,
			dns_duration,
			first_byte_duration,
			redirect_duration,
			redirects_count,
			first_contentful_paint,
			first_paint,
			cumulative_layout_shift,
			first_input_delay,
			largest_contentful_paint,
			geo_country_code,
			geo_city_name,
			page_id,
			data_saver_on,
			boomerang_version,
			screen_width,
			screen_height,
			dom_res,
			dom_doms,
			mem_total,
			mem_limit,
			mem_used,
			mem_lsln,
			mem_ssln,
			mem_lssz,
			scr_bpp,
			scr_orn,
			cpu_cnc,
			dom_ln,
			dom_sz,
			dom_ck,
			dom_img,
			dom_img_uniq,
			dom_script,
			dom_iframe,
			dom_link,
			dom_link_css,
			mob_etype,
			mob_dl,
			mob_rtt
		) VALUES(
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)`,
		p.table,
	)
	err := p.conn.Exec(context.Background(), query,
		theEvent.Created_At,
		theEvent.Hostname,
		theEvent.Event_Type,
		theEvent.Browser_Name,
		nullString(theEvent.Browser_Version),
		nullString(theEvent.Ua_Vnd),
		nullString(theEvent.Ua_Plt),
		theEvent.Device_Type,
		nullString(theEvent.Device_Manufacturer),
		theEvent.Operating_System,
		nullString(theEvent.Operating_System_Version),
		nullString(theEvent.User_Agent),
		theEvent.Next_Hop_Protocol,
		theEvent.Visibility_State,
		theEvent.Session_Id,
		theEvent.Session_Length,
		theEvent.Url,
		nullString(theEvent.Connect_Duration),
		nullString(theEvent.Dns_Duration),
		nullString(theEvent.First_Byte_Duration),
		nullString(theEvent.Redirect_Duration),
		theEvent.Redirects_Count,
		nullString(theEvent.First_Contentful_Paint),
		nullString(theEvent.First_Paint),
		nullNumber(theEvent.Cumulative_Layout_Shift),
		nullNumber(theEvent.First_Input_Delay),
		nullString(theEvent.Largest_Contentful_Paint),
		theEvent.Geo_Country_Code,
		nullString(theEvent.Geo_City_Name),
		theEvent.Page_Id,
		nullNumber(theEvent.Data_Saver_On),
		theEvent.Boomerang_Version,
		nullString(theEvent.Screen_Width),
		nullString(theEvent.Screen_Height),
		nullString(theEvent.Dom_Res),
		nullString(theEvent.Dom_Doms),
		nullString(theEvent.Mem_Total),
		nullString(theEvent.Mem_Limit),
		nullString(theEvent.Mem_Used),
		nullString(theEvent.Mem_Lsln),
		nullString(theEvent.Mem_Ssln),
		nullString(theEvent.Mem_Lssz),
		nullString(theEvent.Scr_Bpp),
		nullString(theEvent.Scr_Orn),
		nullString(theEvent.Cpu_Cnc),
		nullString(theEvent.Dom_Ln),
		nullString(theEvent.Dom_Sz),
		nullString(theEvent.Dom_Ck),
		nullString(theEvent.Dom_Img),
		nullString(theEvent.Dom_Img_Uniq),
		nullString(theEvent.Dom_Script),
		nullString(theEvent.Dom_Iframe),
		nullString(theEvent.Dom_Link),
		nullString(theEvent.Dom_Link_Css),
		nullString(theEvent.Mob_Etype),
		nullNumber(theEvent.Mob_Dl),
		nullNumber(theEvent.Mob_Rtt),
	)
	if err != nil {
		return fmt.Errorf("clickhouse insert failed: %w", err)
	}
	return nil
}

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

func nullNumber(value json.Number) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: value.String(),
		Valid:  true,
	}
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
