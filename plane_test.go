package main

import (
	"math"
	"testing"
)

// Run this with : go test -bench=Benchmark* -cpu 1 -benchmem
func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()

	arena := NewArena()
	const planesCount = 1
	planes := [planesCount]*Plane{}

	for x := 0; x < planesCount; x++ {
		planes[x] = NewPlane(uint32(x), arena)
	}

	b.StartTimer()

	deltaT := 1.0 / 60.0

	for i := 0; i < b.N; i++ {

		for x := 0; x < planesCount; x++ {
			arena.snapshotInputs[uint32(x)] = &PlayerInput{plane: planes[x], data: []byte{3, 4, 4, 4, 4}}
		}

		generateSnapshot(arena, deltaT)

	}

}

func TestPlaneMovement(t *testing.T) {

	arena := NewArena()

	p := NewPlane(1, arena)

	p.location.y = 0
	p.speed.z = 10

	p.calculateSpeed(1 / 60)

}

func TestLocalSpeed(t *testing.T) {
	arena := NewArena()

	plane := NewPlane(1, arena)

	plane.speed = Vector3D{
		x: 0,
		y: 0,
		z: 10,
	}

	plane.maxRot.y = 1
	plane.inputsAxes.yaw = math.Pi / 2

	plane.orientation = plane.calculateRotation(1)

	plane.mass = 1

	localSpeed := plane.localSpeed()

	predictedLocalSpeed := Vector3D{
		x: -10,
		y: 0,
		z: 0,
	}

	if math.Abs(localSpeed.x-predictedLocalSpeed.x) > 0.0001 {
		t.Errorf("LocalSpeed.x is different from prediction. %f != %f", float32(localSpeed.x), float32(predictedLocalSpeed.x))
	}

	if math.Abs(localSpeed.y-predictedLocalSpeed.y) > 0.0001 {
		t.Errorf("LocalSpeed.y is different from prediction. %f != %f", float32(localSpeed.y), float32(predictedLocalSpeed.y))
	}

	if math.Abs(localSpeed.z-predictedLocalSpeed.z) > 0.0001 {
		t.Errorf("LocalSpeed.z is different from prediction. %f != %f", float32(localSpeed.z), float32(predictedLocalSpeed.z))
	}

}
