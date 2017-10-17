package main

import "testing"
import "math"
import "log"

func TestMultiplyByMatrix3(t *testing.T) {

	pitchMat := makeMatrix3X(0)
	yawMat := makeMatrix3Y(math.Pi / 2)
	rollMat := makeMatrix3Z(0)

	orientation := newMatrix3()

	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	localRotMat = localRotMat.getInverse()

	rotMat := localRotMat.Mul(localRotMat)

	rotMat = rotMat.Mul(orientation)

	v := vector3D{0, 0, 1}

	r := v.multiplyByMatrix3(&localRotMat)

	log.Println(r)
}
