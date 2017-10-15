package main

import (
	"encoding/binary"
	"math"
	"time"
)

type axes struct {
	roll, pitch, yaw float64
}

// Plane describe a plane with all its properties
type Plane struct {
	uid         uint32 // Plane owner's player uid
	inputsAxes  axes
	inputThrust float64
	location    vector3D
	orientation matrix3
	deltaRot    vector3D
	speed       vector3D // unit / seconds
	maxSpeed    float64
	maxRotation vector3D // All in radians / seconds
	updatedLast time.Time
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint32) *Plane {
	return &Plane{
		uid: uid,

		inputsAxes: axes{
			roll:  0,
			pitch: 0,
			yaw:   0,
		},

		inputThrust: 0,

		location: vector3D{
			x: 0,
			y: 500,
			z: 0,
		},
		orientation: newMatrix3(),
		deltaRot: vector3D{
			x: 0,
			y: math.Pi / 2,
			z: 0,
		},
		speed: vector3D{
			x: 0,
			y: 0,
			z: 0,
		},
		maxSpeed: 20000,
		maxRotation: vector3D{
			x: 1.5,
			y: 1.5,
			z: 1.5,
		},
		updatedLast: time.Now(),
	}
}

// UpdateIntoBuffer updates the plane's properties from new parameters
// and puts them into a buffer (first arg)
func (p *Plane) UpdateIntoBuffer(buf []byte, offset int, params []byte, tick time.Time) {

	// Calculate the time since the last time updated
	deltaT := tick.Sub(p.updatedLast).Seconds()

	// Set updatedLast to the current tick
	p.updatedLast = tick

	if len(params) > 0 { // We update those only if we have data

		p.inputsAxes.roll = -float64(int8(params[1])) / 127
		p.inputsAxes.pitch = float64(int8(params[2])) / 127
		p.inputsAxes.yaw = float64(int8(params[3])) / 127

		p.inputThrust = float64(uint8(params[4])) / 255

		// HACK: speed multiplied from thrust
		p.speed.z = p.inputThrust * p.maxSpeed
	}

	// deltaRot is only used for display on the client side
	p.deltaRot.x = p.maxRotation.x * p.inputsAxes.pitch * deltaT
	p.deltaRot.y = p.maxRotation.y * p.inputsAxes.yaw * deltaT
	p.deltaRot.z = p.maxRotation.z * p.inputsAxes.roll * deltaT

	mov := p.calculateMovement()

	p.location.x += mov.x * deltaT
	p.location.y += (mov.y - 980) * deltaT
	p.location.z += mov.z * deltaT

	// Dump everything into the slice
	binary.BigEndian.PutUint32(buf[offset:], p.uid)

	binary.BigEndian.PutUint32(buf[offset+4:], math.Float32bits(float32(p.location.x)))
	binary.BigEndian.PutUint32(buf[offset+8:], math.Float32bits(float32(p.location.y)))
	binary.BigEndian.PutUint32(buf[offset+12:], math.Float32bits(float32(p.location.z)))

	binary.BigEndian.PutUint32(buf[offset+16:], math.Float32bits(float32(p.deltaRot.x))) // X
	binary.BigEndian.PutUint32(buf[offset+20:], math.Float32bits(float32(p.deltaRot.y))) // Y
	binary.BigEndian.PutUint32(buf[offset+24:], math.Float32bits(float32(p.deltaRot.z))) // Z

}

func (p *Plane) calculateMovement() vector3D {

	pitchMat := makeMatrix3X(p.deltaRot.x)
	yawMat := makeMatrix3Y(p.deltaRot.y)
	rollMat := makeMatrix3Z(p.deltaRot.z)

	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	// invertedOrientation := p.orientation.getInverse()

	// rotMat := invertedOrientation.Mul(localRotMat)

	rotMat := p.orientation.Mul(localRotMat)

	mov := p.speed.multiplyByMatrix3(&rotMat)

	p.orientation = rotMat

	return mov

}
