package gamelist

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

func Equal(t *testing.T, attr string, expected, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		if attr == "" {
			t.Fatalf("Expected %v, got %v", expected, got)
		} else {
			t.Fatalf("Expected %s to be %v, got %v", attr, expected, got)
		}
	}
}

func StartsWith(t *testing.T, attr string, expected, got string) {
	if !strings.HasPrefix(expected, got) {
		if attr == "" {
			t.Fatalf("Expected %v, got %v", expected, got)
		} else {
			t.Fatalf("Expected %s to be %v, got %v", attr, expected, got)
		}
	}
}

func subWrapper[T1 any, T2 any](f func(*testing.T, string, T1, T2)) func(*testing.T, string, T1, T2) {
	return func(t *testing.T, attr string, expected T1, got T2) {
		t.Run(attr, func(t *testing.T) {
			f(t, attr, expected, got)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	// An XML test string

	game := Game{}
	if err := xml.Unmarshal([]byte(tetris), &game); err != nil {
		t.Error(err)
	}

	sEqual := subWrapper(Equal)
	sStartsWith := subWrapper(StartsWith)

	sEqual(t, "ID", game.ID, int(2976))
	sEqual(t, "Source", game.Source, "ScreenScraper.fr")
	sEqual(t, "Path", game.Path, "./Tetris (World) (Rev A).zip")
	sEqual(t, "Name", game.Name, "Tetris")
	sStartsWith(t, "Desc", game.Desc, "This version of Tetris is")
	sEqual(t, "Rating", game.Rating, float32(0.8))
	sEqual(t, "ReleaseDate", game.ReleaseDate, "19890602T000000")
	sEqual(t, "Lastplayed", game.Lastplayed, "20220401T092851")
	sEqual(t, "Developer", game.Developer, "Nintendo")
	sEqual(t, "Publisher", game.Publisher, "Nintendo")
	sEqual(t, "Genre", game.Genre, "Puzzle-Game / Fall-Puzzle-Game")
	sEqual(t, "Players", game.Players, "1-2")
	sEqual(t, "Hash", game.Hash, "46DF91AD")
	sEqual(t, "Image", game.Image, "./Imgs/Tetris (World) (Rev A).png")
	sEqual(t, "GenreID", game.GenreID, int(2816))
}

func TestRoundTrip(t *testing.T) {
	game := Game{
		ID:     2976,
		Source: "ScreenScraper.fr",
		Path:   "./Tetris (World) (Rev A).zip",
		Name:   "Tetris",
		Desc: `This version of Tetris is one of many conversions of the famous block-stacking game, and was included with the Game Boy upon its release in several territories. The goal is to place pieces made up of four tiles in a ten-by-twenty well, organizing them into complete rows, which then disappear. As rows are cleared, the pace of the game increases and the background changes, and the game ends if the stack reaches the top of the well.

The game is very similar to Nintendo's own NES version of the game, featuring the same "Type A" endless and "Type B" set-clear modes. The game also features a 2-player versus mode that can be played with two Game Boys, two copies of Tetris, and a Game Boy link cable. Clearing lines in this mode will cause the other player's stack to rise, with the goal being to make the other player lose.`,
		Rating:      0.8,
		ReleaseDate: "19890602T000000",
		Lastplayed:  "20220401T092851",
		Developer:   "Nintendo",
		Publisher:   "Nintendo",
		Genre:       "Puzzle-Game / Fall-Puzzle-Game",
		Players:     "1-2",
		Hash:        "46DF91AD",
		Image:       "./Imgs/Tetris (World) (Rev A).png",
		GenreID:     2816,
	}

	xmlText, err := xml.MarshalIndent(game, "", "  ")
	if err != nil {
		t.Error(err)
	}

	roundTripGame := Game{}
	if err := xml.Unmarshal([]byte(xmlText), &roundTripGame); err != nil {
		t.Error(err)
	}

	Equal(t, "XML", game, roundTripGame)
}
