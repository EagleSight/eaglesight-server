package main

import (
	"encoding/binary"
	"log"
	"testing"
	"time"
)

func TestSnapshotGeneration(t *testing.T) {

	arena := newArena()
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

func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()

	deltaT := time.Now()
	arena := newArena()
	const planesCount = 32
	planes := [planesCount]*Plane{}

	for x := 0; x < planesCount; x++ {
		planes[x] = NewPlane(uint32(x))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		deltaT.Add(time.Second / 60)

		for x := 0; x < planesCount; x++ {
			arena.snapshotInputs[uint32(x)] = &PlayerInput{plane: planes[x], data: []byte{3, 4, 4, 4, 4}}
		}

		generateSnapshot(arena, deltaT)

	}

}
