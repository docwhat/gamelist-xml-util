package miyoogamelist_test

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

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
			var data interface{}

			xmlgamelist, err := os.ReadFile(path)
			suite.Require().NoError(err)
			suite.Require().Empty(xmlgamelist)

			_ = xml.Unmarshal(xmlgamelist, &data)

			suite.Require().NotEmpty(data)
		}

		return nil
	})

	suite.NoError(err)
}

func TestFileTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(FileTestSuite))
}
