package main

import (
	"bytes"
	"testing"
	"time"
)

func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()

	deltaT := time.Now()
	arena := newArena()
	const planesCount = 1
	planes := [planesCount]*Plane{}

	for x := 0; x < planesCount; x++ {
		planes[x] = NewPlane(uint32(x))
	}

	snapshotBuffer := new(bytes.Buffer)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		deltaT.Add(time.Second / 60)

		for x := 0; x < planesCount; x++ {
			arena.snapshotInputs[uint32(x)] = &PlayerInput{plane: planes[x], data: nil}
		}

		generateSnapshotBytes(snapshotBuffer, arena, deltaT)

		snapshotBuffer.Reset()
	}

}
