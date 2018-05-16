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

	localRotMat = *localRotMat.GetInverse()

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

func TestAddVector3D(t *testing.T) {
	v1 := Vector3D{X: 3.0, Y: 4.0, Z: 5.0}
	v2 := Vector3D{X: 6.0, Y: 3.0, Z: 1.0}

	v3 := v1.Add(v2)

	if v3.X != v1.X+v2.X {
		t.Errorf("v3.X should be %f + %f. Not %f.", v1.X, v2.X, v3.X)
	}

	if v3.Y != v1.Y+v2.Y {
		t.Errorf("v3.Y should be %f + %f. Not %f.", v1.Y, v2.Y, v3.Y)
	}

	if v3.Z != v1.Z+v2.Z {
		t.Errorf("v3.Z should be %f + %f. Not %f.", v1.Z, v2.Z, v3.Z)
	}
}

func TestDivScalar(t *testing.T) {
	v := Vector3D{X: 3.0, Y: 4.0, Z: 5.0}

	v2 := v.DivScalar(2.0)

	if v2.X != v.X/2.0 {
		t.Errorf("v2.X should be %f not %f", v.X/2, v2.X)
	}

	if v2.Y != v.Y/2.0 {
		t.Errorf("v2.Y should be %f not %f", v.Y/2, v2.Y)
	}

	if v2.Z != v.Z/2.0 {
		t.Errorf("v2.Z should be %f not %f", v.Z/2, v2.Z)
	}
}

func TestMultiplyScalar(t *testing.T) {
	v := Vector3D{X: 3.0, Y: 4.0, Z: 5.0}

	v2 := v.MulScalar(2.0)

	if v2.X != v.X*2.0 {
		t.Fail()
	}

	if v2.Y != v.Y*2.0 {
		t.Fail()
	}

	if v2.Z != v.Z*2.0 {
		t.Fail()
	}
}

func TestHighestInTriangle(t *testing.T) {
	var triangle [3]Vector3D

	triangle[0] = Vector3D{X: 3.0, Y: 5.0, Z: 3.0}
	triangle[1] = Vector3D{X: 1.0, Y: 8.0, Z: -5.0} // The hieghest
	triangle[2] = Vector3D{X: -5.0, Y: 1.0, Z: 1.0}

	i := HighestInTriangle(&triangle)

	if i != 8.0 {
		t.Fail()
	}

	// Round 2
	triangle[0] = Vector3D{X: 3.0, Y: 2.0, Z: 3.0}
	triangle[1] = Vector3D{X: 1.0, Y: 3.0, Z: -5.0} // The hieghest
	triangle[2] = Vector3D{X: -5.0, Y: 4.0, Z: 1.0}

	i = HighestInTriangle(&triangle)

	if i != 4.0 {
		t.Fail()
	}
}
