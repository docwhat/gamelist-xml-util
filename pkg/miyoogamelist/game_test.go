package miyoogamelist

import (
	"reflect"
	"testing"
)

type RandomStruct struct {
	Path  string
	Name  string
	Image string
	Other string
}

func (r RandomStruct) GetMiyooData() Game {
	return Game{
		Path:  r.Path,
		Name:  r.Name,
		Image: r.Image,
	}
}

func TestImport(t *testing.T) {
	given := RandomStruct{
		Path:  "./Path.zip",
		Name:  "name",
		Image: "./Image.png",
		Other: "other",
	}

	expected := Game{
		Path:  "./Path.zip",
		Name:  "name",
		Image: "./Image.png",
	}

	gotten := GameImport(given)

	if !reflect.DeepEqual(expected, gotten) {
		t.Fatalf("Expected %v, got %v", expected, gotten)
	}
}
