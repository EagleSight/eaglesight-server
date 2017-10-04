package main

import (
	"testing"
	"time"
)

func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()

	deltaT := time.Now()
	arena := newArena()
	const planesCount = 100
	planes := [planesCount]*Plane{}

	for x := 0; x < planesCount; x++ {
		planes[x] = NewPlane(uint32(x))
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		deltaT.Add(time.Second / 60)

		for x := 0; x < planesCount; x++ {
			arena.snapshotInputs[uint32(x)] = &PlayerInput{plane: planes[x], data: nil}
		}

		generateSnapshotBytes(arena, deltaT)
	}

}
