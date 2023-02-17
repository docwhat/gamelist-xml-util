package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"docwhat.org/gamelist-xml-util/pkg/gamelist"
	cli "github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	app := cli.NewApp()
	app.Usage = "CLI for managing and manipulating game list data"
	app.Authors = []*cli.Author{
		{
			Name:  "Christian Höltje",
			Email: "docwhat@gerf.org",
		},
	}
	app.Flags = []cli.Flag{
		//nolint:exhaustruct
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "-",
			Usage:   "path to output file (use '-' for stdout)",
		},
		//nolint:exhaustruct
		&cli.StringFlag{
			Name:    "roms",
			Aliases: []string{"r"},
			Value:   "",
			Usage:   "path to ROMs directory (default: gamelist.xml dir or current dir if gamelist.xml is stdin)",
		},
		//nolint:exhaustruct
		&cli.BoolFlag{
			Name:    "add-roms",
			Aliases: []string{"a"},
			Usage:   "add games from ROMs directory to gamelist.xml that are not in the gamelist.xml",
		},
		//nolint:exhaustruct
		&cli.BoolFlag{
			Name:    "add-images",
			Aliases: []string{"i"},
			Usage:   "add images to games in gamelist.xml that are missing an image but the image exists in ROMs directory",
		},
	}
	app.Version = version

	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	// Determine gamelist.xml path
	gamelistPath := "-"
	if ctx.NArg() > 0 {
		gamelistPath = ctx.Args().First()
	}

	// Open gamelist.xml
	gamelistReader, err := getReader(gamelistPath)
	if err != nil {
		return err
	}
	defer gamelistReader.Close()

	// Determine output writer
	outputWriter, err := getWriter(ctx.String("output"))
	if err != nil {
		return err
	}
	defer outputWriter.Close()

	// Determine ROMs path
	romsPath := getRomsDir(ctx.String("roms"), gamelistPath)

	// Load gameList.xml
	gameList, err := gamelist.Load(gamelistReader)
	if err != nil {
		return fmt.Errorf("unable to load gamelist.xml: %w", err)
	}

	// Add ROMs to gamelist.xml
	if ctx.Bool("add-roms") {
		err = gameList.AddGame(romsPath)
		if err != nil {
			return fmt.Errorf("unable to add ROMs: %w", err)
		}
	}

	// Add missing images
	if ctx.Bool("add-images") {
		if err := addImages(gameList, romsPath); err != nil {
			return err
		}
	}

	// Write gamelist.xml to output
	err = gameList.Write(outputWriter)
	if err != nil {
		return fmt.Errorf("unable to write gamelist.xml: %w", err)
	}

	return nil
}

func getRomsDir(romsPath, gamelistPath string) string {
	if romsPath == "" {
		if gamelistPath == "-" {
			romsPath = "."
		} else {
			romsPath = filepath.Dir(gamelistPath)
		}
	}

	return romsPath
}

// addImages adds missing images to the given gameList.
func addImages(gameList *gamelist.GameList, romsPath string) error {
	for _, game := range gameList.Games {
		err := game.AddImage(romsPath)
		if err != nil {
			return fmt.Errorf("unable to add image for %s: %w", game.Path, err)
		}
	}

	return nil
}

// getWriter returns a writer for the given path or stdout if path is "-".
func getWriter(outputPath string) (*os.File, error) {
	outputWriter := os.Stdout

	if outputPath != "-" {
		var err error
		outputWriter, err = os.Create(outputPath)

		if err != nil {
			return nil, fmt.Errorf("unable to write to %s: %w", outputPath, err)
		}
	}

	return outputWriter, nil
}

// getReader returns a reader for the given path or stdin if path is "-".
func getReader(gamelistPath string) (*os.File, error) {
	gamelistReader := os.Stdin

	if gamelistPath != "-" {
		var err error
		gamelistReader, err = os.Open(gamelistPath)

		if err != nil {
			return nil, fmt.Errorf("unable to open %s: %w", gamelistPath, err)
		}
	}

	return gamelistReader, nil
}
