package world

import (
	"log"
	"testing"
	"time"
)

func getTestWorld() *World {
	terrain, err := LoadTerrain("../map.esmap")

	if err != nil {
		log.Fatalln(err)
	}

	return NewWorld(terrain)
}

func TestNewWorld(t *testing.T) {

	w := getTestWorld()
	log.Print("world created")
	go func(world *World) {
		time.Sleep(time.Second)
		log.Print("Should end now")
		w.End <- false
	}(w)

	go func(world *World) {
		for {
			<-w.Snapshots
		}
	}(w)

	w.Run(time.Second/100, time.Second/20)
}

func TestAddBullet(t *testing.T) {

	w := getTestWorld()

	b := Bullet{}

	w.addBullet(&b)

	if len(w.bullets) == 0 {
		t.Fail()
	}
}

func TestAddPlane(t *testing.T) {
	w := getTestWorld()

	w.addPlane(1, PlaneModel{}, w.gun)

	if len(w.planes) == 0 {
		t.Fail()
	}

	w.applyInput(&PlayerInput{UID: 1, Data: []byte{0x3, 12, 12, 12, 26, 0x80}})

	if !w.planes[1].input.IsFiring {
		t.Fail()
	}
}
