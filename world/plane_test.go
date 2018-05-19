package world

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/eaglesight/eaglesight-backend/mathutils"
)

func dummyPlane(uid uint8) *Plane {
	model := PlaneModel{
		MaxThrust: 50000,
		Mass:      4000,
		MaxRotations: mathutils.Vector3D{
			X: 0.314159265358979,
			Y: 0.314159265358979,
			Z: 1,
		},
		DragFactors: mathutils.Vector3D{
			X: 0.05,
			Y: 0.005,
			Z: 0.05,
		},
		LiftMin:      0.0005,
		LiftMax:      0.0007,
		DefaultSpeed: 150,
	}

	gun := make(chan Bullet)

	return NewPlane(uid, model, gun)
}

func TestSnapshot(t *testing.T) {
	var uid uint8 = 5
	plane := dummyPlane(uid)

	plane.location.X = 32
	plane.location.Y = 19
	plane.location.Z = 90

	snap := make([]byte, PlaneSnapshotSize)

	plane.Read(snap)

	if uint8(snap[0]) != uid {
		t.Fail()
	}

	// Location
	if math.Float32frombits(binary.BigEndian.Uint32(snap[2:6])) != float32(plane.location.X) {
		t.Fail()
	}

	if math.Float32frombits(binary.BigEndian.Uint32(snap[6:10])) != float32(plane.location.Y) {
		t.Fail()
	}

	if math.Float32frombits(binary.BigEndian.Uint32(snap[10:14])) != float32(plane.location.Z) {
		t.Fail()
	}

	rot := plane.orientation.ToEulerAngle()

	// Rotation
	if math.Float32frombits(binary.BigEndian.Uint32(snap[14:18])) != float32(rot.X) {
		t.Fail()
	}

	if math.Float32frombits(binary.BigEndian.Uint32(snap[18:22])) != float32(rot.Y) {
		t.Fail()
	}

	if math.Float32frombits(binary.BigEndian.Uint32(snap[22:26])) != float32(rot.Z) {
		t.Fail()
	}
}

func TestGetAirDensity(t *testing.T) {
	if getAirDensity(500) != 1.2 {
		t.Fail()
	}
}
