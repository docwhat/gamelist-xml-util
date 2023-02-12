package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	app := cli.NewApp()
	app.Usage = "Cleans up miyoogamelist.xml files"
	app.Authors = []*cli.Author{
		{
			Name:  "Christian Höltje",
			Email: "docwhat@gerf.org",
		},
	}
	app.UsageText = "miyoogamelist-cleaner [options] <gamelist.xml>"
	app.Version = version

	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
