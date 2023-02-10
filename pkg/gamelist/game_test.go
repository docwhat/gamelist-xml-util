package gamelist_test

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/gamelist"
	"github.com/stretchr/testify/suite"
)

type GameListSuite struct {
	suite.Suite
	testdata string
}

func (suite *GameListSuite) SetupTest() {
	topdir, err := os.Getwd()
	suite.Require().NoError(err)

	suite.testdata = filepath.Join(filepath.Dir(filepath.Dir(topdir)), "testdata", "miyoogamelist")
}

func (suite *GameListSuite) TestUnmarshal() {
	//nolint:exhaustruct
	game := gamelist.Game{}

	if err := xml.Unmarshal([]byte(tetris), &game); err != nil {
		suite.FailNow("error unmarshalling: %w", err)
	}

	suite.Equal(game.ID, int(2976))
	suite.Equal(game.Source, "ScreenScraper.fr")
	suite.Equal(game.Path, "./Tetris (World) (Rev A).zip")
	suite.Equal(game.Name, "Tetris")
	suite.Contains(game.Desc, "This version of Tetris is")
	suite.Equal(game.Rating, float32(0.8))
	suite.Equal(game.ReleaseDate, "19890602T000000")
	suite.Equal(game.Lastplayed, "20220401T092851")
	suite.Equal(game.Developer, "Nintendo")
	suite.Equal(game.Publisher, "Nintendo")
	suite.Equal(game.Genre, "Puzzle-Game / Fall-Puzzle-Game")
	suite.Equal(game.Players, "1-2")
	suite.Equal(game.Hash, "46DF91AD")
	suite.Equal(game.Image, "./Imgs/Tetris (World) (Rev A).png")
	suite.Equal(game.GenreID, int(2816))
}

func (suite *GameListSuite) TestRoundTrip() {
	game := gamelist.Game{
		XMLName: xml.Name{Space: "", Local: "game"},
		ID:      2976,
		Source:  "ScreenScraper.fr",
		Path:    "./Tetris (World) (Rev A).zip",
		Name:    "Tetris",
		Desc: `This version of Tetris is one of many conversions of the ` +
			`famous block-stacking game, and was included with the Game Boy ` +
			`upon its release in several territories. The goal is to place ` +
			`pieces made up of four tiles in a ten-by-twenty well, organizing ` +
			`them into complete rows, which then disappear. As rows are ` +
			`cleared, the pace of the game increases and the background changes, ` +
			`and the game ends if the stack reaches the top of the well.` +
			"\n\n" +
			`The game is very similar to Nintendo's own NES version of the ` +
			`game, featuring the same "Type A" endless and "Type B" set-clear ` +
			`modes. The game also features a 2-player versus mode that can be ` +
			`played with two Game Boys, two copies of Tetris, and a Game Boy link ` +
			`cable. Clearing lines in this mode will cause the other player's ` +
			`stack to rise, with the goal being to make the other player lose.`,
		Rating:      0.8,
		ReleaseDate: "19890602T000000",
		Lastplayed:  "20220401T092851",
		Developer:   "Nintendo",
		Publisher:   "Nintendo",
		Genre:       "Puzzle-Game / Fall-Puzzle-Game",
		Players:     "1-2",
		Hash:        "46DF91AD",
		Image:       "./Imgs/Tetris (World) (Rev A).png",
		Thumbnail:   "./Imgs/Tetris (World) (Rev A) Thumbnail.png",
		GenreID:     2816,
		PlayCount:   321,
	}

	var roundTripGame gamelist.Game

	if xmlText, err := xml.MarshalIndent(game, "", "  "); err != nil {
		suite.FailNow("error marshalling: %w", err)
	} else if err := xml.Unmarshal(xmlText, &roundTripGame); err != nil {
		suite.FailNow("error unmarshalling %v: %w", game, err)
	}

	suite.Equal(game, roundTripGame)
}

func (suite *GameListSuite) TestReadingWithTestData() {
	err := filepath.Walk(suite.testdata, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "miyoogamelist.xml" {
			suite.Run(filepath.Base(filepath.Dir(path)), func() {
				var data gamelist.GameList

				if xmlgamelist, err := os.ReadFile(path); err != nil {
					suite.FailNowf("error reading %s: %w", path, err)
				} else if err := xml.Unmarshal(xmlgamelist, &data); err != nil {
					suite.FailNowf("error unmarshalling %s: %w", path, err)
				}

				suite.NotEmpty(data.Games)
			})
		}

		return nil
	})

	suite.NoError(err)
}

func TestGameList(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(GameListSuite))
}
