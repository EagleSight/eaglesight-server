package main

import (
	"encoding/binary"
	"log"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestSnapshotGeneration(t *testing.T) {

	arena := NewArena()
	const planesCount = 32

	arena.snapshotInputs[uint32(1)] = &PlayerInput{plane: NewPlane(uint32(1)), data: nil}
	arena.snapshotInputs[uint32(2)] = &PlayerInput{plane: NewPlane(uint32(2)), data: nil}

	deltaT := time.Now().Add(time.Second / 60)

	b := generateSnapshot(arena, deltaT)

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

	log.Println(b)

}

func TestConnectPlayer(t *testing.T) {

	arena := NewArena()

	player := NewPlayer(1, arena, new(websocket.Conn))

	arena.connectPlayer(player)

	if len(arena.players) != 1 {
		t.Fail()
	}

	arena.deconnectPlayer(player)

	if len(arena.players) != 0 {
		t.Fail()
	}

}
