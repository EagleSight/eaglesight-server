package main

// Player
type Player struct {
	Name  string     `json:"username"`
	UUID  string     `json:"token"`
	Model PlaneModel `json:"planeModel"`
}
