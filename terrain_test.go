package main

import "testing"

func TestLoad(t *testing.T) {

	terrain, err := LoadTerrain("")

	if err != nil {
		t.Error("Problem while loading terrain")
	}

	if len(terrain.points) != int(terrain.width*terrain.depth) {
		t.Error("Not enough points loaded")
	}

}
