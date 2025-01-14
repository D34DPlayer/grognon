package database

import (
	"context"
	"embed"

	"github.com/pkg/errors"
	"github.com/yaitoo/sqle/migrate"
)

//go:embed migrations
var migrations embed.FS

func (db *Database) Migrate(ctx context.Context) error {
	m := migrate.New(db.DB)

	if err := m.Discover(migrations); err != nil {
		return errors.Wrap(err, "failed to discover migrations")
	}

	err := m.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to initialize migrations")
	}

	err = m.Migrate(ctx)
	return errors.Wrap(err, "failed to migrate")
}
