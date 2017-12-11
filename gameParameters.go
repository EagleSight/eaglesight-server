package main

import (
	"encoding/json"
	"io"
)

// GameParameters contains all the parameters of the game
type GameParameters struct {
	GameID     string          `json:"gameId"`
	Players    []PlayerProfile `json:"players"`
	TerrainURL string          `json:"terrainURL"`
}

// TEST THIS!
// DefaultGameParameters return the default parameters of a game
func DefaultGameParameters() GameParameters {

	return GameParameters{
		GameID:     "",
		Players:    []PlayerProfile{},
		TerrainURL: "",
	}

}

// TEST THIS!
// DecodeAndUpdate decode the parameters of a game
func (gp *GameParameters) DecodeAndUpdate(src io.ReadCloser) error {

	defer src.Close()

	data := GameParameters{}

	// Decode the json
	if err := json.NewDecoder(src).Decode(data); err != nil {
		return nil
	}

	gp = &data

	return nil
}
