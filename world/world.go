package world

import (
	"log"
	"time"
)

// PlayerInput contains input data and the uid to which it is attributed
type PlayerInput struct {
	UID  uint8
	Data []byte
}

// World World is which everything happens
type World struct {
	Input     chan PlayerInput
	Snapshots chan []byte
	join      chan struct {
		UID   uint8
		Model PlaneModel
	}
	leave   chan uint8
	End     chan bool // End the world
	gun     chan *Bullet
	terrain *Terrain
	planes  map[uint8]*Plane
	bullets []*Bullet
}

// NewWorld Creates a new world
func NewWorld(terrain *Terrain) *World {

	world := &World{
		terrain:   terrain,
		planes:    make(map[uint8]*Plane),
		Snapshots: make(chan []byte, 1),
		Input:     make(chan PlayerInput, 1),
		join: make(chan struct {
			UID   uint8
			Model PlaneModel
		}, 1),
		leave:   make(chan uint8, 1),
		End:     make(chan bool),
		gun:     make(chan *Bullet, 1),
		bullets: []*Bullet{},
	}
	return world
}

// Join ...
func (w *World) Join(uid uint8, model PlaneModel) {

	w.join <- struct {
		UID   uint8
		Model PlaneModel
	}{UID: uid, Model: model}
}

// Leave ...
func (w *World) Leave(uid uint8) {

	w.leave <- uid
}

// addPlane add a plane to the world
func (w *World) addPlane(uid uint8, model PlaneModel, gun chan<- *Bullet) {
	// Check if the plane already exists in the world
	w.planes[uid] = NewPlane(uid, model, gun)
}

// addBullet add the bullet in the world
func (w *World) addBullet(bullet *Bullet) {

	w.bullets = append(w.bullets, bullet)
}

// applyInput apply a UserInput to a plane in the world
func (w *World) applyInput(input *PlayerInput) {

	plane, exists := w.planes[input.UID]

	if exists {
		plane.Write(input.Data) // So the plane can process the data by itself
	}
}

// updateWorld updates the states of all the entities in the world
func (w *World) updateWorld(deltaT float64) {

	bulletsStillAlive := []*Bullet{}

	// Update all the bullets
	for _, bullet := range w.bullets {

		if bullet.Update(deltaT) {
			// The bullet is still alive. It goes to the next round.
			bulletsStillAlive = append(bulletsStillAlive, bullet)
		}
	}
	w.bullets = bulletsStillAlive

	// Update all the planes
	for _, plane := range w.planes {
		plane.Update(deltaT, w.terrain)
	}
}

func (w *World) removePlane(uid uint8) {

	if _, exists := w.planes[uid]; exists {
		delete(w.planes, uid)
	}
}

// generateSnapshots generate a snapshot of the whole world
func (w *World) generateSnapshots() []byte {

	const snapshotSizeOverhead = 1 // opcode's length
	snapshot := make([]byte, snapshotSizeOverhead+len(w.planes)*PlaneSnapshotSize)
	snapshot[0] = 0x3
	offset := snapshotSizeOverhead

	for _, plane := range w.planes {
		plane.Read(snapshot[offset:])
		offset += PlaneSnapshotSize
	}
	return snapshot
}

// Run starts and run the world
func (w *World) Run(simulationInterval time.Duration, snapshotInterval time.Duration) {

	simulationTimer := time.Tick(simulationInterval)
	snapshotTimer := time.Tick(snapshotInterval)
	lastTick := time.Now()

	for {

		select {
		case <-w.End:
			log.Println("STOP!")
			return
		case <-snapshotTimer:
			w.Snapshots <- w.generateSnapshots()
		case now := <-simulationTimer:
			w.updateWorld(now.Sub(lastTick).Seconds())
			lastTick = now
		case input := <-w.Input:
			w.applyInput(&input)
		case bullet := <-w.gun:
			w.addBullet(bullet)
		case plane := <-w.join:
			log.Println("Plane joining")
			w.addPlane(plane.UID, plane.Model, w.gun)
		case uid := <-w.leave:
			log.Println("Plane leaving")
			w.removePlane(uid)
		}

	}
}
