package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"d34d.one/grognon/internal/backend"
	"d34d.one/grognon/internal/database"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Data string
}

func backgroundTask(ctx context.Context, duration time.Duration, task func()) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func action(ctx context.Context, cfg Config) error {
	// Create data folder if it does not exist
	if _, err := os.Stat(cfg.Data); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.Data, 0755); err != nil {
			return cli.Exit(err, 1)
		}
	}

	db, err := database.Setup(ctx, cfg.Data)
	if err != nil {
		return cli.Exit(err, 1)
	}
	cons, err := database.SetupConnections(db)
	if err != nil {
		return cli.Exit(err, 1)
	}
	database.ReflectAll(db, cons)

	backgroundTask(ctx, 15*time.Minute, func() {
		database.ReflectAll(db, cons)
	})

	// TODO: Use latest run time instead of fixed intervals
	backgroundTask(ctx, 1*time.Minute, func() {
		database.ExecuteCrons(db, cons, "minute")
	})
	backgroundTask(ctx, 1*time.Hour, func() {
		database.ExecuteCrons(db, cons, "hour")
	})
	backgroundTask(ctx, 1*24*time.Hour, func() {
		database.ExecuteCrons(db, cons, "day")
	})

	// DEBUG STUFF
	// err = database.AddCron(db, cons, database.Cron{
	// 	ConnectionId: 0,
	// 	Name:         "User count",
	// 	Command:      "SELECT count(*) as count FROM streamer",
	// 	Schedule:     "minute",
	// })
	// if err != nil {
	// 	return cli.Exit(err, 1)
	// }
	// end DEBUG STUFF

	if err := backend.Setup(); err != nil {
		return cli.Exit(err, 1)
	}

	<-ctx.Done()
	return nil
}

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "data",
			Value:   "./data",
			Aliases: []string{"d"},
			Usage:   "Folder where the data is stored",
		},
	}
	cmd := cli.Command{
		Name:  "grognon",
		Usage: "Scavage for statistics",
		Flags: flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			config := Config{
				Data: cmd.String("data"),
			}
			return action(ctx, config)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("Error: ", slog.Any("error", err))
	}
}
