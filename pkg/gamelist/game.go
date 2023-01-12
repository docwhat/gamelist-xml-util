package gamelist

type Game struct {
	Path string `xml:"path"`
	Name string `xml:"name"`
	Desc string `xml:"desc"`

	Image       string  `xml:"image,omitempty"`     // The path to the game's image.
	Thumbnail   string  `xml:"thumbnail,omitempty"` // The path to the game's thumbnail image.
	Rating      float32 `xml:"rating,omitempty"`
	ReleaseDate string  `xml:"releasedate,omitempty"`
	Developer   string  `xml:"developer,omitempty"`
	Publisher   string  `xml:"publisher,omitempty"`
	Genre       string  `xml:"genre,omitempty"`     // The (primary) genre for the game.
	Players     string  `xml:"players,omitempty"`   // The number of players the game supports.
	PlayCount   int     `xml:"playcount,omitempty"` // The number of times this game has been played
	Lastplayed  string  `xml:"lastplayed,omitempty"`

	// Miyoo Mini/OnionOS Extensions
	Hash    string `xml:"hash,omitempty"`
	GenreID int    `xml:"genreid,omitempty"`
	ID      int    `xml:"id,attr,omitempty"`
	Source  string `xml:"source,attr,omitempty"`
}
