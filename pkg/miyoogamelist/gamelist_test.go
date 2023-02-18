package miyoogamelist_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
	"github.com/stretchr/testify/suite"
)

type GameListSuite struct {
	suite.Suite
	TestData string
}

func (suite *GameListSuite) SetupTest() {
	topdir, err := os.Getwd()
	suite.Require().NoError(err)

	suite.TestData = filepath.Join(filepath.Dir(filepath.Dir(topdir)), "testdata", "miyoogamelist")
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

	gameList, err := miyoogamelist.Load(strings.NewReader(xmlText))
	suite.Require().NoError(err)

	game := gameList.Games[0]
	provider := gameList.Provider

	suite.Equal("Game Boy Color", provider.System)
	suite.Equal("Skraper", provider.Software)
	suite.Equal("ScreenScraper.fr", provider.Database)
	suite.Equal("http://www.screenscraper.fr", provider.Web)

	suite.Equal("007: The World is Not Enough", game.Name)
	suite.Equal("./007 - The World Is Not Enough (USA, Europe).zip", game.Path)
	suite.Equal("E038E666", game.Hash)
	suite.Equal("./Imgs/007 - The World Is Not Enough (USA, Europe).png", game.Image)
}

func (suite *GameListSuite) TestReadingWithTestData() {
	err := filepath.Walk(suite.TestData, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "miyoogamelist.xml" {
			gameList, err := miyoogamelist.LoadFile(path)
			suite.Require().NoError(err)

			suite.NotEmpty(gameList)
			suite.NotEmpty(gameList.Provider.System)
			suite.NotEmpty(gameList.Games)
		}

		return nil
	})

	suite.NoError(err)
}

func TestFileTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(GameListSuite))
}
