package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	pgxStdlib "github.com/jackc/pgx/v5/stdlib"
	"github.com/yaitoo/sqle"
)

type Database struct {
	PgxPool *pgxpool.Pool
	*sqle.DB
}

func Setup(ctx context.Context, dbUrl string) (*Database, error) {
	slog.Debug("Opening database", slog.String("url", dbUrl))
	pgxConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("db: failed to parse postgres database URL %s: %w", dbUrl, err)
	}
	pgxPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("db: failed to create pgx pool: %w", err)
	}

	sqlDb := pgxStdlib.OpenDBFromPool(pgxPool)
	sqleDb := sqle.Open(sqlDb)
	db := &Database{PgxPool: pgxPool, DB: sqleDb}

	slog.Info("Checking for migrations")
	if err := db.Migrate(ctx); err != nil {
		return nil, fmt.Errorf("db: failed to migrate database: %w", err)
	}
	slog.Info("Database migrated")

	return db, nil
}
