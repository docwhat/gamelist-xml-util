package miyoogamelist_test

import (
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
)

func TestImport(t *testing.T) {
	dummy := map[string]interface{}{
		"gameList": []interface{}{
			map[string]string {
				"path": "path1",
				"name": "name1",
				"image": "image1",
			},
			map[string]string {
				"path": "path2",

				"name": "name2",
				"image": "image2",
			},
		},
	}

	gamelist, err := miyoogamelist.Import(dummy)
	if err != nil {
		t.Error(err)
	}
}
