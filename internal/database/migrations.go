package database

import (
	"context"
	"embed"

	"github.com/yaitoo/sqle/migrate"
)

//go:embed migrations
var migrations embed.FS

func (db *Database) Migrate(ctx context.Context) error {
	m := migrate.New(db.DB)

	if err := m.Discover(migrations); err != nil {
		return err
	}

	err := m.Init(ctx)
	if err != nil {
		return err
	}

	err = m.Migrate(ctx)
	return err
}
