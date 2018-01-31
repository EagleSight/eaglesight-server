package main

import (
	"encoding/json"
	"log"
	"os"
)

// GameParameters contains all the parameters of the game
type GameParameters struct {
	GameID  string   `json:"gameId"`
	Players []Player `json:"players"`
}

// LoadGameParametersFromFile load the parameters of a game from a local JSON file (TEST THIS?)
func LoadGameParametersFromFile() GameParameters {

	reader, err := os.Open("./game.json")
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	data := GameParameters{}

	// Decode the json
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		log.Fatal(err)
	}

	return data
}
