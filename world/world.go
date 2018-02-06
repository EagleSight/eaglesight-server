package world

import (
	"log"
)

// State ...
type State struct {
	Planes []PlaneState
}

// World World is which everything happens
type World struct {
	terrain Terrain
	planes  []*Plane
}

// NewWorld Creates a new world
func NewWorld() *World {

	// Try to load the terrain
	log.Print("Loading terrain... ")
	terrain, err := LoadTerrain()
	log.Println("DONE")

	if err != nil {
		log.Fatalln(err)
	}

	world := &World{
		terrain: terrain,
		planes:  []*Plane{},
	}

	return world
}

// AddPlane Add a plane to the world
func (w *World) AddPlane(uid uint8, model PlaneModel) (chan PlaneInput, error) {

	plane := NewPlane(uid, model)

	w.planes = append(w.planes, plane)

	return plane.inputs, nil
}

// Tick simulate a tick in the world
func (w *World) Tick(deltaT float64) (state State) {

	state.Planes = []PlaneState{}

	for i, plane := range w.planes {
		if plane.isDead() {
			w.planes = append(w.planes[:i], w.planes[i+1:]...)
			continue
		}

		state.Planes = append(state.Planes, plane.Tick(deltaT, &w.terrain))
	}

	return state
}
