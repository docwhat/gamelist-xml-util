package miyoogamelist

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// LoadXML reads the XML data from the given reader and returns a GameList.
func LoadXML(r io.Reader) (*GameList, error) {
	gameList := NewGameList()

	if err := xml.NewDecoder(r).Decode(&gameList); err != nil {
		return nil, fmt.Errorf("unable to parse miyoogamelist XML: %w", err)
	}

	return gameList, nil
}

// LoadXMLFile reads the XML data from the given file and returns a GameList.
func LoadXMLFile(path string) (*GameList, error) {
	gameListFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open miyoogamelist XML file: %w", err)
	}

	return LoadXML(gameListFile)
}

// WriteXML takes an io.Writer and writes the GameList to it.
func (g *GameList) WriteXML(writer io.Writer) error {
	_, err := writer.Write([]byte(xml.Header))
	if err != nil {
		return fmt.Errorf("unable to write XML header: %w", err)
	}

	enc := xml.NewEncoder(writer)
	enc.Indent("", "  ")

	if err := enc.Encode(g); err != nil {
		return fmt.Errorf("unable to write gamelist: %w", err)
	}

	return nil
}

// WriteXMLFile writes the GameList to the given file path.
func (g *GameList) WriteXMLFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	return g.WriteXML(file)
}
