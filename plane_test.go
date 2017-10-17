package main

import (
	"math"
	"testing"
	"time"
)

// Run this with : go test -bench=Benchmark* -cpu 1 -benchmem
func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()

	deltaT := time.Now()
	arena := NewArena()
	const planesCount = 1
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

func TestPlaneMovement(t *testing.T) {

	p := NewPlane(1)

	p.location.y = 0
	p.deltaRot.y = -math.Pi / 4
	p.deltaRot.x = -math.Pi / 4
	p.speed.z = 10

	p.calculateMovement()

}
