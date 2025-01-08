package main

import (
	"context"
	"log"
	"os"

	"d34d.one/grognon/internal/database"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Data string
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
	_, err = database.SetupConnections(db)
	if err != nil {
		return cli.Exit(err, 1)
	}

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
		log.Fatal(err)
	}
}
