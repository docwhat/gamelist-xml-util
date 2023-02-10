package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

const (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "gamelist2miyoogamelist"
	app.Usage = "Cleans up miyoogamelist.xml files"
	app.Author = "Christian Höltje"
	app.UsageText = "miyoogamelist-cleaner [options] <miyoo-gamelist.xml>"
	app.Version = version

	app.Flags = []cli.Flag{
		//nolint:exhaustruct
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "remove all bugs",
		},
	}

	app.Action = func(cmd *cli.Context) error {
		filename := cmd.Args().First()

		if filename == "" {
			if err := cli.ShowAppHelp(cmd); err != nil {
				log.Fatal(err)
			}

			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "Hello %v\n", cmd.Bool("debug"))

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
