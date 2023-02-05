package miyoogamelist

// Imports any struct that is miyoogamelist shaped by using reflection.
func Import(unknownStruct map[string]interface{}) (GameList, error) {
	var gamelist GameList
	var games []Game
	for _, game := range unknownStruct["gameList"].([]interface{}) {
		games = append(games, Game{
			Path:  game.(map[string]interface{})["path"].(string),
			Name:  game.(map[string]interface{})["name"].(string),
			Image: game.(map[string]interface{})["image"].(string),
		})
	}
	gamelist.Games = games
	return gamelist, nil
}
