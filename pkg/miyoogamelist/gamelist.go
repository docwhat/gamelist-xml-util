package miyoogamelist

import "encoding/xml"

type GameList struct {
	XMLName xml.Name `xml:"gameList"`

	Games []Game `xml:"game"`
}

type Game struct {
	XMLName xml.Name `xml:"game"`

	Path  string `xml:"path"`
	Name  string `xml:"name"`
	Image string `xml:"image"`
}

func NewGameList(games ...Game) GameList {
	return GameList{
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
