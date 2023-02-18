package gamelist

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type GameList struct {
	XMLName xml.Name `xml:"gameList" exhaustruct:"optional"`

	Provider Provider `xml:"provider,omitempty" exhaustruct:"optional"`
	Games    []Game   `xml:"game"`
}

type Provider struct {
	XMLName xml.Name `xml:"provider" exhaustruct:"optional"`

	System   string `xml:"System,omitempty" exhaustruct:"optional"`
	Software string `xml:"software,omitempty" exhaustruct:"optional"`
	Database string `xml:"database,omitempty" exhaustruct:"optional"`
	Web      string `xml:"web,omitempty" exhaustruct:"optional"`
}

type Game struct {
	XMLName xml.Name `xml:"game" exhaustruct:"optional"`

	Path string `xml:"path"`
	Name string `xml:"name"`
	Desc string `xml:"desc"`

	Image       string  `xml:"image,omitempty" exhaustruct:"optional"`
	Thumbnail   string  `xml:"thumbnail,omitempty" exhaustruct:"optional"`
	Rating      float32 `xml:"rating,omitempty" exhaustruct:"optional"`
	ReleaseDate string  `xml:"releasedate,omitempty" exhaustruct:"optional"`
	Developer   string  `xml:"developer,omitempty" exhaustruct:"optional"`
	Publisher   string  `xml:"publisher,omitempty" exhaustruct:"optional"`
	Genre       string  `xml:"genre,omitempty" exhaustruct:"optional"`
	Players     string  `xml:"players,omitempty" exhaustruct:"optional"`
	PlayCount   int     `xml:"playcount,omitempty" exhaustruct:"optional"`
	Lastplayed  string  `xml:"lastplayed,omitempty" exhaustruct:"optional"`

	// Miyoo Mini/OnionOS Extensions
	Hash    string `xml:"hash,omitempty" exhaustruct:"optional"`
	GenreID int    `xml:"genreid,omitempty" exhaustruct:"optional"`
	ID      int    `xml:"id,attr,omitempty" exhaustruct:"optional"`
	Source  string `xml:"source,attr,omitempty" exhaustruct:"optional"`
}

// NewGameList returns a new empty &GameList.
func NewGameList() *GameList {
	return &GameList{
		Provider: Provider{},
		Games:    []Game{},
	}
}

func NewGame(name, path, desc string) Game {
	return Game{
		Path: path,
		Name: name,
		Desc: desc,
	}
}

func (g *GameList) AddGame(path string) error {
	// Strip directory
	name := filepath.Base(path)
	// Strip extension
	name = name[:len(name)-len(filepath.Ext(name))]

	game := NewGame(name, path, "")
	g.Games = append(g.Games, game)

	return nil
}

func (g *Game) AddImage(romsDir string) error {
	return nil
}

// Load takes an io.Reader and returns a GameList.
func Load(r io.Reader) (*GameList, error) {
	gamelist := NewGameList()

	if err := xml.NewDecoder(r).Decode(gamelist); err != nil {
		return nil, fmt.Errorf("unable to parse gamelist: %w", err)
	}

	return gamelist, nil
}

// LoadFile reads the XML data from the given file and returns a GameList.
func LoadFile(path string) (*GameList, error) {
	gameListFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open gamelist XML file: %w", err)
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
