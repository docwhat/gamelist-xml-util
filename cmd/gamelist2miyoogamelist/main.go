package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
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
	app.UsageText = "miyoogamelist-cleaner [options]"
	app.Version = version

	app.Flags = []cli.Flag{
		//nolint:exhaustruct
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "remove all bugs.",
			EnvVar: "",
			Hidden: true,
		},
		cli.StringFlag{
			Name:        "gamelist, g",
			Usage:       "The gamelist.xml file to clean up. Use '-' for stdin.",
			EnvVar:      "",
			FilePath:    "",
			Required:    true,
			Hidden:      false,
			TakesFile:   true,
			Value:       "",
			Destination: new(string),
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "the file to write the results to. Use '-' for stdout.",
			EnvVar:      "",
			FilePath:    "",
			Required:    false,
			Hidden:      false,
			TakesFile:   true,
			Value:       "-",
			Destination: new(string),
		},
	}

	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(cmd *cli.Context) error {
	// Show help if no args given
	if !cmd.Args().Present() {
		log.Fatal(cli.ShowAppHelp(cmd))
	}

	inFile := os.Stdin
	inPath := cmd.String("gamelist")

	if inPath != "-" {
		var err error

		if inFile, err = os.OpenFile(inPath, os.O_RDONLY, 0); err != nil {
			return fmt.Errorf("unable to open %v for reading: %w", inPath, err)
		}

		defer inFile.Close()
	}

	var data []byte

	var err error

	if data, err = io.ReadAll(inFile); err != nil {
		return fmt.Errorf("unable to write to %v: %w", inPath, err)
	}

	var gamelist miyoogamelist.GameList

	if err := xml.Unmarshal(data, &gamelist); err != nil {
		return fmt.Errorf("unable to parse %v: %w", inFile, err)
	}

	outFile := os.Stdout
	outPath := cmd.String("out")

	if outPath != "-" {
		const mode = 0o755

		var err error

		if outFile, err = os.OpenFile(outPath, os.O_RDWR|os.O_CREATE, mode); err != nil {
			return fmt.Errorf("unable to open %v for writing: %w", outPath, err)
		}

		defer outFile.Close()
	}

	var strippedData []byte

	if strippedData, err = xml.MarshalIndent(gamelist, "", "\t"); err != nil {
		return fmt.Errorf("unable to marshal %w", err)
	}

	if _, err := outFile.Write(strippedData); err != nil {
		return fmt.Errorf("unable to write to %v: %w", outPath, err)
	}

	return nil
}
