package miyoogamelist

type Gamey interface {
	GetMiyooData() Game
}

type Game struct {
	Path  string `xml:"path"`
	Name  string `xml:"name"`
	Image string `xml:"image"`
}

func GameImport[T Gamey](thing T) Game {
	return thing.GetMiyooData()
}
