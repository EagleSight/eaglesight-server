package main

import (
	"math"
	"testing"
)

func TestMatrix3Inv(t *testing.T) {

	m := newMatrix3()

	m._11 = 1
	m._12 = 0
	m._13 = 2

	m._21 = 5
	m._22 = 10
	m._23 = 3

	m._31 = 2
	m._32 = 3
	m._33 = 1

	inverse := m.getInverse()

	var truth matrix3

	truth._11 = -1.0 / 9.0
	truth._12 = -2.0 / 3.0
	truth._13 = 20.0 / 9.0

	truth._21 = -1.0 / 9.0
	truth._22 = 1.0 / 3.0
	truth._23 = -7.0 / 9.0

	truth._31 = 5.0 / 9.0
	truth._32 = 1.0 / 3.0
	truth._33 = -10.0 / 9.0

	var diff float64

	diff += inverse._11 - truth._11
	diff += inverse._12 - truth._12
	diff += inverse._13 - truth._13
	diff += inverse._21 - truth._21
	diff += inverse._22 - truth._22
	diff += inverse._23 - truth._23
	diff += inverse._31 - truth._31
	diff += inverse._32 - truth._32
	diff += inverse._33 - truth._33

	if math.Abs(diff) > 0.0001 {
		t.Errorf("Matrix3 Inverse : Difference between result and truth bigger than 0. Result = %f", diff)
	}

}

func TestMatrix3Mul(t *testing.T) {

	var m1 matrix3
	var m2 matrix3

	// Fill m1
	m1._11 = 1
	m1._12 = 0
	m1._13 = 2

	m1._21 = 5
	m1._22 = 10
	m1._23 = 3

	m1._31 = 2
	m1._32 = 3
	m1._33 = 1

	// Fill m2
	m2._11 = 4
	m2._12 = 2
	m2._13 = 1

	m2._21 = 6
	m2._22 = 4
	m2._23 = 2

	m2._31 = 5
	m2._32 = 3
	m2._33 = 2

	result := m1.Mul(m2)

	var truth matrix3

	truth._11 = 14
	truth._12 = 8
	truth._13 = 5

	truth._21 = 95
	truth._22 = 59
	truth._23 = 31

	truth._31 = 31
	truth._32 = 19
	truth._33 = 10

	var diff float64

	diff += result._11 - truth._11
	diff += result._12 - truth._12
	diff += result._13 - truth._13

	diff += result._21 - truth._21
	diff += result._22 - truth._22
	diff += result._23 - truth._23

	diff += result._31 - truth._31
	diff += result._32 - truth._32
	diff += result._33 - truth._33

	if diff != 0 {
		t.Fail()
	}

}

func TestMatrixToEulerAngle(t *testing.T) {

	x1 := 0.00212
	y1 := 0.00212
	z1 := 0.00212

	pitchMat := makeMatrix3X(x1)
	yawMat := makeMatrix3Y(y1)
	rollMat := makeMatrix3Z(z1)

	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	v2 := localRotMat.ToEulerAngle()

	if math.Abs(x1-v2.x) > 0.00001 {
		t.Errorf("x1 is different from x2: %f != %f", x1, v2.x)
	}

	if math.Abs(y1-v2.y) > 0.00001 {
		t.Errorf("y1 is different from y2: %f != %f", y1, v2.y)
	}

	if math.Abs(z1-v2.z) > 0.001 {
		t.Errorf("z1 is different from z2: %f != %f", z1, v2.z)
	}

	// For singular

	x1 = math.Pi / 2
	y1 = 0
	z1 = 0

	pitchMat = makeMatrix3X(x1)
	yawMat = makeMatrix3Y(y1)
	rollMat = makeMatrix3Z(z1)

	localRotMat = yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	v2 = localRotMat.ToEulerAngle()

	if x1 != v2.x {
		t.Errorf("x1 is different for singular from x2: %f != %f", x1, v2.x)
	}

	if y1 != v2.y {
		t.Errorf("y1 is different for singular  from y2: %f != %f", y1, v2.y)
	}

	if z1 != 0 {
		t.Errorf("z1 is different for singular from 0: %f != 0.0", v2.z)
	}

}
