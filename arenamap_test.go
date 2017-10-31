package main

import "testing"

func TestOveredTriangle(t *testing.T) {

}

func TestLoad(t *testing.T) {
	am := LoadArenaMap()

	if len(am.points) != int(am.width*am.depth) {
		t.Error("Not enough points loaded")
	}

}
