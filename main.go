package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "data",
			Value: "./data",
			Usage: "Folder where the data is stored",
		},
	}
	cmd := cli.Command{
		Name:  "grognon",
		Usage: "Scavage for statistics",
		Flags: flags,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
