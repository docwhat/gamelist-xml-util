package miyoogamelist

type Game struct {
	Path  string `xml:"path"`
	Name  string `xml:"name"`
	Image string `xml:"image"`
}
