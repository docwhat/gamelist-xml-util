package cleaner

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"

	"docwhat.org/gamelist-xml-util/pkg/gamelist"
	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
)

func Clean(input io.ReadCloser, output io.WriteCloser) error {
	gameData := gamelist.GameList{}

	if text, err := ioutil.ReadAll(input); err != nil {
		return err
	} else if err := xml.Unmarshal(text, &gameData); err != nil {
		return err
	}

	miyooData := miyoogamelist.GameImport(gameData)
	if newText, err := xml.MarshalIndent(miyooData, "", "  "); err != nil {
		return err
	}
}

func foo() {
	f, err := os.Open("foo.xml")

	Clean(f, os.Stdout)
}
