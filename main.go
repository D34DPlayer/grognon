package main

import (
	"context"
	"log/slog"
	"os"

	"d34d.one/grognon/internal/backend"
	"d34d.one/grognon/internal/background"
	"d34d.one/grognon/internal/database"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Data    string
	SsrHost string
	DBUrl   string
}

func action(ctx context.Context, cfg Config) error {
	// Create data folder if it does not exist
	if _, err := os.Stat(cfg.Data); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.Data, 0755); err != nil {
			return cli.Exit(err, 1)
		}
	}

	db, err := database.Setup(ctx, cfg.DBUrl)
	if err != nil {
		return cli.Exit(err, 1)
	}
	cons, err := database.SetupConnections(db)
	if err != nil {
		return cli.Exit(err, 1)
	}

	background.SetupReflection(ctx, db, cons)
	background.SetupCronJobs(ctx, db, cons)

	if err := backend.Setup(db, cons, cfg.SsrHost); err != nil {
		return cli.Exit(err, 1)
	}

	<-ctx.Done()
	return nil
}

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "data",
			Value: "./data",
			Usage: "Folder where the data is stored",
		},
		&cli.StringFlag{
			Name:  "ssr",
			Value: "http://127.0.0.1:13714",
			Usage: "Hostname for SSR",
		},
		&cli.StringFlag{
			Name:     "db",
			Required: true,
			Usage:    "Database connection string",
		},
	}
	cmd := cli.Command{
		Name:  "grognon",
		Usage: "Scavage for statistics",
		Flags: flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config := Config{
				Data:    cmd.String("data"),
				SsrHost: cmd.String("ssr"),
				DBUrl:   cmd.String("db"),
			}
			return action(ctx, config)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("Error: ", slog.Any("error", err))
	}
}
