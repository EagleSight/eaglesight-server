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
	uid         uint32
	arena       *Arena
	inputsAxes  axes
	inputThrust float64
	location    Vector3D
	orientation matrix3
	deltaRot    Vector3D
	maxRot      Vector3D // All in radians / seconds
	absRot      Vector3D // Absolute rotation of the plane
	speed       Vector3D // unit / seconds
	maxSpeed    float64
	updatedLast time.Time
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint32, arena *Arena) *Plane {
	return &Plane{
		uid:   uid,
		arena: arena,
		inputsAxes: axes{
			roll:  0,
			pitch: 0,
			yaw:   0,
		},

		inputThrust: 0,

		location: Vector3D{
			x: 0,
			y: 1500,
			z: 0,
		},
		orientation: newMatrix3(),
		deltaRot: Vector3D{
			x: 0,
			y: math.Pi / 2,
			z: 0,
		},
		absRot: Vector3D{
			x: 0,
			y: 0,
			z: 0,
		},
		speed: Vector3D{
			x: 0,
			y: 0,
			z: 0,
		},
		maxSpeed: 3 * 685 * 0.27,
		maxRot: Vector3D{
			x: 1.5,
			y: 1.5,
			z: 1.2,
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
	p.deltaRot.x = p.maxRot.x * p.inputsAxes.pitch * deltaT
	p.deltaRot.y = p.maxRot.y * p.inputsAxes.yaw * deltaT
	p.deltaRot.z = p.maxRot.z * p.inputsAxes.roll * deltaT

	mov := p.calculateMovement()

	p.absRot = p.orientation.ToEulerAngle()

	p.location.x += mov.x * deltaT
	p.location.y += mov.y * deltaT
	p.location.z += mov.z * deltaT

	// Update the position if there is colision
	p.correctFromCollision()

	// Dump everything into the slice
	binary.BigEndian.PutUint32(buf[offset:], p.uid)

	binary.BigEndian.PutUint32(buf[offset+4:], math.Float32bits(float32(p.location.x)))
	binary.BigEndian.PutUint32(buf[offset+8:], math.Float32bits(float32(p.location.y)))
	binary.BigEndian.PutUint32(buf[offset+12:], math.Float32bits(float32(p.location.z)))

	binary.BigEndian.PutUint32(buf[offset+16:], math.Float32bits(float32(p.absRot.x)))
	binary.BigEndian.PutUint32(buf[offset+20:], math.Float32bits(float32(p.absRot.y)))
	binary.BigEndian.PutUint32(buf[offset+24:], math.Float32bits(float32(p.absRot.z)))

}

func (p *Plane) calculateMovement() Vector3D {

	pitchMat := makeMatrix3X(p.deltaRot.x)
	yawMat := makeMatrix3Y(p.deltaRot.y)
	rollMat := makeMatrix3Z(p.deltaRot.z)

	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	p.orientation = p.orientation.Mul(localRotMat)

	mov := p.speed.multiplyByMatrix3(&p.orientation)

	return mov

}

func (p *Plane) correctFromCollision() {

	triangle := p.arena.arenaMap.OverredTriangle(p.location)

	// We are out of bound
	if math.IsNaN(triangle[0].x) {
		return
	}

	h := heightOnTriangle(p.location, &triangle)

	// We are under the surface
	if p.location.y < h+5 {
		// We go back to the surface
		p.location.y = h + 5
	}

}
