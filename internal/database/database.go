package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yaitoo/sqle"
)

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
