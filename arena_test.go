package main

import (
	"encoding/binary"
	"testing"

	"github.com/gorilla/websocket"
)

var terrain, _ = LoadTerrain("")

func TestSnapshotGeneration(t *testing.T) {

	params := DefaultGameParameters()

	arena := NewArena(params, terrain)
	const planesCount = 32

	arena.snapshotInputs[uint32(1)] = &PlayerInput{plane: NewPlane(uint32(1), terrain), data: nil}
	arena.snapshotInputs[uint32(2)] = &PlayerInput{plane: NewPlane(uint32(2), terrain), data: nil}

	b := generateSnapshot(arena, 1/60)

	// # Instruction Type == 3
	if b[0] != 0x3 {
		t.Error("The instruction type is wrong")
	}

	// Tick always equals 1 for this test
	if binary.BigEndian.Uint32(b[1:5]) != 1 {
		t.Error("The tick is wrong")
	}

	// Players count
	if binary.BigEndian.Uint16(b[5:7]) != 2 {
		t.Error("Players count is wrong")
	}

}

func TestConnectPlayer(t *testing.T) {

	params := DefaultGameParameters()

	arena := NewArena(params, terrain)

	profil := DefaultPlayerProfile(1)
	plane := NewPlane(profil.Token, terrain)

	player := NewPlayer(profil, plane, new(websocket.Conn))

	arena.connectPlayer(player)

	if len(arena.players) != 1 {
		t.Fail()
	}

	arena.deconnectPlayer(player)

	if len(arena.players) != 0 {
		t.Fail()
	}

}
