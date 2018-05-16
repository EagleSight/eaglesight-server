package game

import (
	"encoding/json"
	"log"
	"os"
)

// Parameters contains all the parameters of the game
type Parameters struct {
	GameID  string          `json:"gameId"`
	Players []PlayerProfile `json:"profiles"`
}

// LoadGameParametersFromFile load the parameters of a game from a local JSON file (TEST THIS?)
func LoadGameParametersFromFile() Parameters {
	reader, err := os.Open("./game.json")

	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	data := Parameters{}

	// Decode the json
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		log.Fatal(err)
	}
	return data
}
