package miyoogamelist_test

import (
	"reflect"
	"testing"

	"docwhat.org/gamelist-xml-util/pkg/miyoogamelist"
)

type RandomStruct struct {
	Path  string
	Name  string
	Image string
	Other string
}

func (r RandomStruct) GetMiyooData() miyoogamelist.Game {
	return miyoogamelist.Game{
		Path:  r.Path,
		Name:  r.Name,
		Image: r.Image,
	}
}

func TestImport(t *testing.T) {
	t.Parallel()

	given := RandomStruct{
		Path:  "./Path.zip",
		Name:  "name",
		Image: "./Image.png",
		Other: "other",
	}

	expected := miyoogamelist.Game{
		Path:  "./Path.zip",
		Name:  "name",
		Image: "./Image.png",
	}

	gotten := miyoogamelist.GameImport(given)

	if !reflect.DeepEqual(expected, gotten) {
		t.Fatalf("Expected %v, got %v", expected, gotten)
	}
}
