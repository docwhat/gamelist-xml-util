package gamelist_test

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
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

// This tests the Load() function by passing in a Reader containing a gamelist
// XML.
func (suite *GameListSuite) TestLoad() {
	xmlText := `<?xml version="1.0" encoding="UTF-8"?>
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<gameList>
	<provider>
		<System>Game Boy Color</System>
		<software>Skraper</software>
		<database>ScreenScraper.fr</database>
		<web>http://www.screenscraper.fr</web>
	</provider>
	<game id="36846" source="ScreenScraper.fr">
		<path>./007 - The World Is Not Enough (USA, Europe).zip</path>
		<name>007: The World is Not Enough</name>
		<desc>It seems that an MI-6 agent has been killed just before</desc>
		<rating>0.75</rating>
		<releasedate>20010911T000000</releasedate>
		<developer>2n Productions</developer>
		<publisher>Electronic Arts</publisher>
		<genre>Shooter-Action</genre>
		<players>1-2</players>
		<hash>E038E666</hash>
		<image>./Imgs/007 - The World Is Not Enough (USA, Europe).png</image>
		<genreid>256</genreid>
	</game>
</gameList>`

	gameList, err := gamelist.Load(strings.NewReader(xmlText))
	suite.Require().NoError(err)

	game := gameList.Games[0]
	provider := gameList.Provider

	suite.Equal("Game Boy Color", provider.System)
	suite.Equal("Skraper", provider.Software)
	suite.Equal("ScreenScraper.fr", provider.Database)
	suite.Equal("http://www.screenscraper.fr", provider.Web)

	suite.Equal("007: The World is Not Enough", game.Name)
	suite.Equal("./007 - The World Is Not Enough (USA, Europe).zip", game.Path)
	suite.Equal("It seems that an MI-6 agent has been killed just before", game.Desc)
	suite.Equal("20010911T000000", game.ReleaseDate)
	suite.Equal(float32(0.75), game.Rating)
	suite.Equal("2n Productions", game.Developer)
	suite.Equal("Electronic Arts", game.Publisher)
	suite.Equal("Shooter-Action", game.Genre)
	suite.Equal("1-2", game.Players)
	suite.Equal("E038E666", game.Hash)
	suite.Equal("./Imgs/007 - The World Is Not Enough (USA, Europe).png", game.Image)
	suite.Equal(256, game.GenreID)
}

func (suite *GameListSuite) TestLoadWithInvalidXML() {
	_, err := gamelist.Load(strings.NewReader("invalid xml"))
	suite.Error(err)
}

func (suite *GameListSuite) TestWrite() {
	gameList := gamelist.NewGameList()

	path := "./path/to/rom.zip"
	suite.Require().NoError(gameList.AddGame(path))

	var buf bytes.Buffer

	suite.Require().NoError(gameList.Write(&buf))

	suite.Require().NotEmpty(buf.String())
	suite.Contains(buf.String(), path)
}

func (suite *GameListSuite) TestAddGame() {
	gameList := gamelist.NewGameList()
	path := "./path/to/rom.zip"

	suite.Require().NoError(gameList.AddGame(path))

	suite.Require().NotEmpty(gameList.Games)
	suite.Equal(path, gameList.Games[0].Path)
	suite.Equal("rom", gameList.Games[0].Name)
	suite.Equal("", gameList.Games[0].Desc)
}

func TestGameList(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(GameListSuite))
}
