package miyoogamelist

import "fmt"

var ErrFailedToImport = fmt.Errorf("failed to import")

// Imports any struct that is miyoogamelist shaped by using reflection.
func Import(unknownStruct map[string]interface{}) (*GameList, error) {
	gameList := GameList{
		Games: []Game{},
	}

	unknownGameList, ok := unknownStruct["gameList"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: No such attribute 'gameList' in %v (%T)", ErrFailedToImport, unknownStruct, unknownStruct)
	}

	for _, unknownGame := range unknownGameList {
		game := Game{
			Path:  "",
			Name:  "",
			Image: "",
		}

		if path, ok := unknownGame.(map[string]string)["path"]; ok {
			game.Path = path
		} else {
			return nil, fmt.Errorf("%w: No such attribute 'path' in %v (%T)", ErrFailedToImport, unknownGame, unknownGame)
		}

		if name, ok := unknownGame.(map[string]string)["name"]; ok {
			game.Name = name
		} else {
			return nil, fmt.Errorf("%w: No such attribute 'name' in %v (%T)", ErrFailedToImport, unknownGame, unknownGame)
		}

		if image, ok := unknownGame.(map[string]string)["image"]; ok {
			game.Image = image
		} else {
			return nil, fmt.Errorf("%w: No such attribute 'image' in %v (%T)", ErrFailedToImport, unknownGame, unknownGame)
		}

		gameList.Games = append(gameList.Games, game)
	}

	return &gameList, nil
}
