package database

import (
	"context"

	"github.com/jackc/tern/v2/migrate"
	"github.com/pkg/errors"
)

func (db *Database) Migrate(ctx context.Context) error {
	poolCon, err := db.PgxPool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire database connection")
	}
	defer poolCon.Release()

	m, err := migrate.NewMigrator(ctx, poolCon.Conn(), "migrations")
	if err != nil {
		return errors.Wrap(err, "failed to create migrator")
	}

	m.Migrations = []*migrate.Migration{
		{
			Sequence: 1,
			Name:     "v0.0.1",
			UpSQL: `
CREATE TABLE connections (
    connection_id     SERIAL PRIMARY KEY,
    db_type           TEXT    NOT NULL,
    connection_url    TEXT    NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,

    connected         BOOLEAN NOT NULL DEFAULT FALSE,
    last_connected_at TIMESTAMPTZ,
    last_error        TEXT
);

CREATE TABLE tables (
    connection_id INTEGER NOT NULL REFERENCES connections (connection_id),
    table_name    TEXT    NOT NULL,

    CONSTRAINT tables_pk PRIMARY KEY (connection_id, table_name)
);

CREATE TABLE columns (
    connection_id INTEGER NOT NULL REFERENCES connections (connection_id),
    table_name    TEXT    NOT NULL,
    name          TEXT    NOT NULL,
    type          TEXT    NOT NULL,
    "notnull"     BOOLEAN NOT NULL,
    dflt_value    TEXT,
    pk            INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT columns_pk PRIMARY KEY (connection_id, table_name, name)
);

CREATE TABLE crons (
    cron_id       SERIAL PRIMARY KEY,
    connection_id INTEGER   NOT NULL REFERENCES connections (connection_id),
    name          TEXT      NOT NULL,
    command       TEXT      NOT NULL,
    schedule      TEXT      NOT NULL,

    created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    last_run_at   TIMESTAMPTZ
);

CREATE TABLE cron_outputs (
    cron_id INTEGER NOT NULL REFERENCES crons (cron_id),
    name    TEXT    NOT NULL,
    type    TEXT    NOT NULL,

    CONSTRAINT cron_outputs_pk PRIMARY KEY (cron_id, name)
);
            `,
			DownSQL: `
DROP TABLE connections;
DROP TABLE tables;
DROP TABLE columns;
DROP TABLE crons;
DROP TABLE cron_outputs;
`,
		},
	}

	if err := m.Migrate(ctx); err != nil {
		return errors.Wrap(err, "failed to migrate")
	}

	return nil
}
