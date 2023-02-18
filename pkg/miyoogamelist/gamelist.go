package miyoogamelist

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type GameList struct {
	XMLName xml.Name `xml:"gameList"`

	Provider Provider `xml:"provider,omitempty" exhaustruct:"optional"`
	Games    []Game   `xml:"game"`
}

type Provider struct {
	XMLName xml.Name `xml:"provider,omitempty"`

	System   string `xml:"System,omitempty" exhaustruct:"optional"`
	Software string `xml:"software,omitempty" exhaustruct:"optional"`
	Database string `xml:"database,omitempty" exhaustruct:"optional"`
	Web      string `xml:"web,omitempty" exhaustruct:"optional"`
}

type Game struct {
	XMLName xml.Name `xml:"game"`

	ID     string `xml:"id,attr,omitempty" exhaustruct:"optional"`
	Source string `xml:"source,attr,omitempty" exhaustruct:"optional"`

	Path  string `xml:"path"`
	Name  string `xml:"name"`
	Image string `xml:"image"`
	Hash  string `xml:"hash,omitempty" exhaustruct:"optional"`
}

func NewGameList(games ...Game) *GameList {
	return &GameList{
		XMLName: xml.Name{Space: "", Local: "gameList"},
		Games:   games,
	}
}

func NewGame(path, name, image string) Game {
	return Game{
		XMLName: xml.Name{Space: "", Local: "game"},
		Path:    path,
		Name:    name,
		Image:   image,
	}
}

// Load reads the XML data from the given reader and returns a GameList.
func Load(r io.Reader) (*GameList, error) {
	gameList := NewGameList()

	if err := xml.NewDecoder(r).Decode(&gameList); err != nil {
		return nil, fmt.Errorf("unable to parse miyoogamelist XML: %w", err)
	}

	return gameList, nil
}

// LoadFile reads the XML data from the given file and returns a GameList.
func LoadFile(path string) (*GameList, error) {
	gameListFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open miyoogamelist XML file: %w", err)
	}

	return Load(gameListFile)
}

// Write takes an io.Writer and writes the GameList to it.
func (g *GameList) Write(writer io.Writer) error {
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

// WriteFile writes the GameList to the given file path.
func (g *GameList) WriteFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	return g.Write(file)
}
