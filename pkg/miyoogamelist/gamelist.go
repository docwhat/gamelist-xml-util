package miyoogamelist

import (
	"encoding/xml"

	"docwhat.org/gamelist-xml-util/pkg/gamelist"
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

	ID     int    `xml:"id,attr,omitempty" exhaustruct:"optional"`
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

// Downgrade converts a gamelist.GameList into a miyoogamelist.GameList.
func Downgrade(gamelist *gamelist.GameList) *GameList {
	provider := Provider{
		XMLName:  xml.Name{Space: "", Local: "provider"},
		System:   gamelist.Provider.System,
		Software: gamelist.Provider.Software,
		Database: gamelist.Provider.Database,
		Web:      gamelist.Provider.Web,
	}

	games := make([]Game, len(gamelist.Games))
	for i, game := range gamelist.Games {
		games[i] = Game{
			XMLName: xml.Name{Space: "", Local: "game"},
			ID:      game.ID,
			Source:  game.Source,
			Path:    game.Path,
			Name:    game.Name,
			Image:   game.Image,
			Hash:    game.Hash,
		}
	}

	return &GameList{
		XMLName:  xml.Name{Space: "", Local: "gameList"},
		Provider: provider,
		Games:    games,
	}
}
