package miyoogamelist_test

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
	testdata string
}

func (suite *FileTestSuite) SetupTest() {
	topdir, err := os.Getwd()
	suite.Require().NoError(err)

	suite.testdata = filepath.Join(filepath.Dir(filepath.Dir(topdir)), "testdata", "miyoogamelist")
}

func (suite *FileTestSuite) TestReadingWithTestData() {
	err := filepath.Walk(suite.testdata, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "miyoogamelist.xml" {
			var data miyoogamelist.GameList

			if xmlgamelist, err := os.ReadFile(path); err != nil {
				suite.T().Fatalf("error reading %s: %v", path, err)
			} else if err := xml.Unmarshal(xmlgamelist, &data); err != nil {
				suite.T().Fatalf("error unmarshalling %s: %v", path, err)
			}

			suite.NotEmpty(data.Games)
		}

		return nil
	})

	suite.NoError(err)
}

func TestFileTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(FileTestSuite))
}
