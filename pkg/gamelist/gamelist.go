package gamelist

import "encoding/xml"

type GameList struct {
	XMLName xml.Name `xml:"gameList"`

	Games []Game `xml:"game"`
}

type Game struct {
	XMLName xml.Name `xml:"game"`

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
