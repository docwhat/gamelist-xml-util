package miyoogamelist_test

import (
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	t.Parallel()

	dummy := map[string]interface{}{
		"gameList": []interface{}{
			map[string]string{
				"path":  "path1",
				"name":  "name1",
				"image": "image1",
			},
			map[string]string{
				"path":  "path2",
				"name":  "name2",
				"image": "image2",
			},
		},
	}

	gamelist, err := miyoogamelist.Import(dummy)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(gamelist.Games), 2)
	assert.Equal(t, gamelist.Games[0].Path, "path1")
	assert.Equal(t, gamelist.Games[0].Name, "name1")
	assert.Equal(t, gamelist.Games[0].Image, "image1")
	assert.Equal(t, gamelist.Games[1].Path, "path2")
	assert.Equal(t, gamelist.Games[1].Name, "name2")
	assert.Equal(t, gamelist.Games[1].Image, "image2")
}
