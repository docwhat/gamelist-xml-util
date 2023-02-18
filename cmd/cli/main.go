package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"docwhat.org/gamelist-xml-util/pkg/gamelist"
	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

// ErrUnknownFormat is returned when the format specified via the command line is not known.
var ErrUnknownFormat = fmt.Errorf("unknown format")

//nolint:funlen // this is almost entirely configuring CLI flags.
func main() {
	//nolint:exhaustruct
	app := &cli.App{
		Name:      "gamelist-cli",
		Usage:     "CLI for managing and manipulating ROM gamelist data files",
		UsageText: "gamelist-cli [options] [gamelist.xml]",
		Authors: []*cli.Author{
			{
				Name:  "Christian Höltje",
				Email: "docwhat@gerf.org",
			},
		},
		Suggest:                true,
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			//nolint:exhaustruct
			&cli.PathFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "-",
				Usage:   "write output to `FILE`; use \"-\" for stdout",
			},
			//nolint:exhaustruct
			&cli.PathFlag{
				Name:    "roms",
				Aliases: []string{"r"},
				Value:   "",
				Usage:   "path to ROMs `DIR` (default: gamelist.xml dir or current dir if gamelist.xml is stdin)",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "gamelist-xml",
				Usage:   "output format (gamelist-xml, miyoogamelist-xml)",
				Action: func(ctx *cli.Context, format string) error {
					switch format {
					case "gamelist-xml", "miyoogamelist-xml":
						return nil
					default:
						return fmt.Errorf("unknown format %q: %w", format, ErrUnknownFormat)
					}
				},
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
		},
		Version: version,
		Action:  action,
	}

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
	gameList, err := gamelist.LoadXML(gamelistReader)
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

	return writeFormat(ctx.String("format"), gameList, outputWriter)
}

func writeFormat(format string, gameList *gamelist.GameList, outputWriter *os.File) error {
	var err error

	switch format {
	case "gamelist-xml":
		// Write gamelist.xml to output
		err = gameList.WriteXML(outputWriter)
	case "miyoogamelist-xml":
		// Write miyoo gamelist.xml to output
		err = miyoogamelist.Downgrade(gameList).WriteXML(outputWriter)
	default:
		panic("programming error; please recycle the programmer.")
	}

	if err != nil {
		return fmt.Errorf("unable to write %s: %w", format, err)
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
