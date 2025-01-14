package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yaitoo/sqle"
)

type EpochTime struct {
	time.Time
	Valid bool
}

func (t *EpochTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}
	t.Valid = true
	switch value := value.(type) {
	case int64:
		t.Time = time.Unix(value, 0)
	case string:
		t.Time, _ = time.Parse(time.RFC3339, value)
	default:
		return fmt.Errorf("unsupported type for EpochTime: %T", value)
	}
	return nil
}

type Database struct {
	*sqle.DB
}

func Setup(ctx context.Context, dataDir string) (*Database, error) {
	log.Println("Opening database...")
	sqldb, err := sql.Open("sqlite3", "file:"+dataDir+"/grognon.db?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("db: failed to open sqlite database %s: %w", dataDir, err)
	}

	sqleDb := sqle.Open(sqldb)

	db := &Database{sqleDb}

	log.Println(("Checking for migrations..."))
	if err := db.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("db: failed to migrate database: %w", err)
	}
	log.Println("Database migrated")

	return db, nil
}
