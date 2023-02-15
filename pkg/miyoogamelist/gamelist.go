package miyoogamelist

import "encoding/xml"

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
