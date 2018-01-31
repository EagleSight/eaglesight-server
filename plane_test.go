package main

import (
	"math"
	"testing"
)

// Run this with : go test -bench=Benchmark* -cpu 1 -benchmem
func BenchmarkUpdatePlane(b *testing.B) {
	b.StopTimer()

	params := DefaultGameParameters()
	terrain, _ := LoadTerrain()

	arena := NewArena(params, terrain)

	const planesCount = 1
	planes := [planesCount]*Plane{}

	for x := 0; x < planesCount; x++ {
		planes[x] = NewPlane(uint8(x), terrain, PlaneModel{})
	}

	b.StartTimer()

	deltaT := 1.0 / 60.0

	for i := 0; i < b.N; i++ {

		for x := 0; x < planesCount; x++ {
			arena.snapshotInputs[uint8(x)] = &PlayerInput{plane: planes[x], data: []byte{3, 4, 4, 4, 4}}
		}

		generateSnapshot(arena, deltaT)

	}

}

func TestPlaneMovement(t *testing.T) {

	terrain, _ := LoadTerrain()

	p := NewPlane(1, terrain, PlaneModel{})

	p.Location.Y = 0
	p.Speed.Z = 10

	p.calculateSpeed(1 / 60)

}

func TestLocalSpeed(t *testing.T) {

	terrain, _ := LoadTerrain()

	plane := NewPlane(1, terrain, PlaneModel{})

	plane.Speed = Vector3D{
		X: 0,
		Y: 0,
		Z: 10,
	}

	plane.Model.MaxRotations.Y = 1
	plane.InputsAxes.Yaw = math.Pi / 2

	plane.Orientation = plane.calculateRotation(1)

	plane.Model.Mass = 1

	localSpeed := plane.getLocalSpeed()

	predictedLocalSpeed := Vector3D{
		X: -10,
		Y: 0,
		Z: 0,
	}

	if math.Abs(localSpeed.X-predictedLocalSpeed.X) > 0.0001 {
		t.Errorf("LocalSpeed.x is different from prediction. %f != %f", float32(localSpeed.X), float32(predictedLocalSpeed.X))
	}

	if math.Abs(localSpeed.Y-predictedLocalSpeed.Y) > 0.0001 {
		t.Errorf("LocalSpeed.y is different from prediction. %f != %f", float32(localSpeed.Y), float32(predictedLocalSpeed.Y))
	}

	if math.Abs(localSpeed.Z-predictedLocalSpeed.Z) > 0.0001 {
		t.Errorf("LocalSpeed.z is different from prediction. %f != %f", float32(localSpeed.Z), float32(predictedLocalSpeed.Z))
	}

}
