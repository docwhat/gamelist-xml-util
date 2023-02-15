package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
	cli "github.com/urfave/cli/v2"
)

type fileWriter interface {
	io.Writer
	io.StringWriter
	io.Closer
}

func action(cmd *cli.Context) error {
	// Show help if no args given
	if !cmd.Args().Present() {
		log.Fatal(cli.ShowAppHelp(cmd))
	}

	gamelist, err := parseFile(cmd.Args().First())
	if err != nil {
		return err
	}

	return writeFile(os.Stdout, gamelist)
}

func parseFile(filename string) (miyoogamelist.GameList, error) {
	var data []byte

	var err error

	if data, err = os.ReadFile(filename); err != nil {
		return miyoogamelist.GameList{}, fmt.Errorf("unable to read %v: %w", filename, err)
	}

	var gamelist miyoogamelist.GameList

	if err := xml.Unmarshal(data, &gamelist); err != nil {
		return miyoogamelist.GameList{}, fmt.Errorf("unable to parse %v: %w", filename, err)
	}

	return gamelist, nil
}

func writeFile(file fileWriter, gamelist miyoogamelist.GameList) error {
	strippedData, err := xml.MarshalIndent(gamelist, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal %w", err)
	}

	defer file.Close()

	if _, err = file.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"); err != nil {
		//nolint:wrapcheck
		return err
	}

	if _, err = file.Write(strippedData); err != nil {
		//nolint:wrapcheck
		return err
	}

	if _, err = file.WriteString("\n"); err != nil {
		//nolint:wrapcheck
		return err
	}

	return nil
}
