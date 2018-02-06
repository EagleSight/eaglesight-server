package mathutils

import "testing"
import "math"

func TestMultiplyByMatrix3(t *testing.T) {

	pitchMat := MakeMatrix3X(0)
	yawMat := MakeMatrix3Y(math.Pi / 2)
	rollMat := MakeMatrix3Z(0)

	orientation := NewMatrix3()

	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	localRotMat = localRotMat.GetInverse()

	rotMat := localRotMat.Mul(localRotMat)

	rotMat = rotMat.Mul(orientation)

	//v := Vector3D{0, 0, 1}

	//	r := v.multiplyByMatrix3(&localRotMat)

}

func TestHeightOnTriangle(t *testing.T) {

	point := Vector3D{
		X: 1.0,
		Y: 4.0,
		Z: 1.0,
	}

	triangle := [3]Vector3D{
		Vector3D{
			X: 0.0,
			Y: 2.0,
			Z: 0.0,
		},
		Vector3D{
			X: 2.0,
			Y: 0.0,
			Z: -1.0,
		},
		Vector3D{
			X: 2.0,
			Y: 0.0,
			Z: 1.0,
		},
	}

	y := HeightOnTriangle(point, &triangle)

	if y != 1.0 {
		t.Errorf("y = %f, not 1.0", y)
	}

	// PART 2

	triangle = [3]Vector3D{
		Vector3D{
			X: 0.0,
			Y: 3.0,
			Z: 2.0,
		},
		Vector3D{
			X: 2.0,
			Y: 1.0,
			Z: -3.0,
		},
		Vector3D{
			X: 2.0,
			Y: 1.0,
			Z: 1.0,
		},
	}

	y = HeightOnTriangle(point, &triangle)

	if y != 2.0 {
		t.Errorf("y = %f, not 2.0", y)
	}

}
