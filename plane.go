package main

import (
	"math"
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
	maxRot      Vector3D // All in radians / seconds
	absRot      Vector3D // Absolute rotation of the plane
	speed       Vector3D // unit / seconds
	maxThrust   float64
	mass        float64
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
		absRot: Vector3D{
			x: 0,
			y: 0,
			z: 0,
		},
		speed: Vector3D{
			x: 0,
			y: 0,
			z: 100,
		},
		maxThrust: 70000,
		maxRot: Vector3D{
			x: 0.2,
			y: 1.5,
			z: 1.2,
		},
		mass: 4000,
	}
}

// Update updates the plane's properties from new parameters
// and puts them into a buffer (first arg)
func (p *Plane) Update(inputs []byte, deltaT float64) {

	if len(inputs) > 0 { // We update those only if we have data
		p.inputsAxes.roll = -float64(int8(inputs[1])) / 127
		p.inputsAxes.pitch = float64(int8(inputs[2])) / 127
		p.inputsAxes.yaw = float64(int8(inputs[3])) / 127

		p.inputThrust = float64(uint8(inputs[4])) / 255
	}

	// Update the rotation
	p.orientation = p.calculateRotation(deltaT)

	p.absRot = p.orientation.ToEulerAngle()

	// Update the speed
	p.speed = p.calculateSpeed(deltaT)

	p.location.x += p.speed.x * deltaT
	p.location.y += p.speed.y * deltaT
	p.location.z += p.speed.z * deltaT

	// Update the position if there is colision
	p.correctFromCollision()
}

func (p *Plane) calculateRotation(deltaT float64) matrix3 {
	pitchMat := makeMatrix3X(p.maxRot.x * p.inputsAxes.pitch * deltaT)
	yawMat := makeMatrix3Y(p.maxRot.y * p.inputsAxes.yaw * deltaT)
	rollMat := makeMatrix3Z(p.maxRot.z * p.inputsAxes.roll * deltaT)

	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	return p.orientation.Mul(localRotMat)
}

func (p *Plane) calculateSpeed(deltaT float64) Vector3D {

	localAcceleration := Vector3D{
		x: 0,
		y: p.calculateLift(),
		z: p.calculateThrust() - p.calculateDrag(),
	}

	globalAcceleration := localAcceleration.multiplyByMatrix3(&p.orientation)

	// Apply gravity
	globalAcceleration.y += -9.8

	acceleration := globalAcceleration.MulScalar(deltaT)

	return p.speed.Add(&acceleration)
}

func (p *Plane) calculateLift() float64 {

	// The p.inputsAxes.pitch should affect the amount of lift

	return 15
}

// calculateDrag calculate the amount of drag
func (p *Plane) calculateDrag() float64 {

	// The p.inputsAxes should affect the amount of drag

	return 0
}

func (p *Plane) calculateThrust() float64 {

	return p.inputThrust * p.maxThrust / p.mass

}

func (p *Plane) correctFromCollision() {

	const margin = 5

	triangle := p.arena.arenaMap.OverredTriangle(p.location)

	// We are out of bound
	if math.IsNaN(triangle[0].x) {
		return
	}

	// Small optimization
	if p.location.y >= highestInTriangle(triangle)+margin {
		return
	}

	// The real thing
	h := heightOnTriangle(p.location, &triangle)

	// We are under the surface
	if p.location.y < h+margin {
		// We go back to the surface
		p.location.y = h + margin
		p.speed.y = 0
	}

}

func (p *Plane) localSpeed() Vector3D {

	inverseOrientation := p.orientation.getInverse()

	return p.speed.multiplyByMatrix3(&inverseOrientation)

}
