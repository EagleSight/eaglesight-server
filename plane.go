package main

import (
	"math"
)

// Axes of a plane
type Axes struct {
	Roll, Pitch, Yaw float64
}

type minMax struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// PlaneFlightProps are all the constant properties that can easily be loaded from a JSON object
type PlaneFlightProps struct {
	MaxThrust    float64  `json:"maxThrust"`
	Mass         float64  `json:"mass"`
	MaxRotations Vector3D `json:"maxRotations"` // All in radians / seconds
	DragFactors  Vector3D `json:"dragFactors"`  // A drag factor are "Area * Drag Coeficient * 0.5" of a side
	LiftRange    minMax   `json:"liftRange"`
}

// Plane describe a plane with all its properties
type Plane struct {
	UID         uint32
	Arena       *Arena
	InputsAxes  Axes
	InputThrust float64
	Location    Vector3D
	Orientation matrix3
	Rotation    Vector3D // Absolute rotation of the plane
	Speed       Vector3D // unit / seconds
	Props       PlaneFlightProps
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint32, arena *Arena) *Plane {
	return &Plane{
		UID:   uid,
		Arena: arena,
		InputsAxes: Axes{
			Roll:  0,
			Pitch: 0,
			Yaw:   0,
		},

		InputThrust: 0,

		Location: Vector3D{
			X: 0,
			Y: 1500,
			Z: 0,
		},
		Orientation: newMatrix3(),
		Rotation: Vector3D{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Speed: Vector3D{
			X: 0,
			Y: 0,
			Z: 100,
		},
		Props: PlaneFlightProps{
			MaxThrust: 120000,
			MaxRotations: Vector3D{
				X: 0.2,
				Y: 1.5,
				Z: 1.2,
			},
			Mass: 4000,
			DragFactors: Vector3D{
				X: 0.02,
				Y: 0.02,
				Z: 0.02,
			},
		},
	}
}

// Update updates the plane's properties from new parameters
// and puts them into a buffer (first arg)
func (p *Plane) Update(inputs []byte, deltaT float64) {

	if len(inputs) > 0 { // We update those only if we have data
		p.InputsAxes.Roll = -float64(int8(inputs[1])) / 127
		p.InputsAxes.Pitch = float64(int8(inputs[2])) / 127
		p.InputsAxes.Yaw = float64(int8(inputs[3])) / 127

		p.InputThrust = float64(uint8(inputs[4])) / 255
	}

	// Update the rotation
	p.Orientation = p.calculateRotation(deltaT)

	p.Rotation = p.Orientation.ToEulerAngle()

	// Update the speed
	p.Speed = p.calculateSpeed(deltaT)

	p.Location.X += p.Speed.X * deltaT
	p.Location.Y += p.Speed.Y * deltaT
	p.Location.Z += p.Speed.Z * deltaT

	// Update the position if there is colision
	p.correctFromCollision()
}

func (p *Plane) calculateRotation(deltaT float64) matrix3 {

	// Generate the matrices that represent the rotation change
	pitchMat := makeMatrix3X(p.Props.MaxRotations.X * p.InputsAxes.Pitch * deltaT)
	yawMat := makeMatrix3Y(p.Props.MaxRotations.Y * p.InputsAxes.Yaw * deltaT)
	rollMat := makeMatrix3Z(p.Props.MaxRotations.Z * p.InputsAxes.Roll * deltaT)

	// Multiply them together in the right order
	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	return p.Orientation.Mul(localRotMat)
}

func (p *Plane) calculateSpeed(deltaT float64) Vector3D {

	localAcceleration := Vector3D{
		X: 0,
		Y: p.calculateLift(),
		Z: p.calculateThrust(),
	}

	// Compute the drag
	localDrag := p.calculateDrag()

	// Divide the force by the mapp
	localDrag.DivScalar(p.Props.Mass)

	// Add the drag to the local acceleration
	localAcceleration.Add(&localDrag)

	// Convert to global
	globalAcceleration := localAcceleration.multiplyByMatrix3(&p.Orientation)

	// Apply gravity
	globalAcceleration.Y += -9.8

	// Multiply by the time
	globalAcceleration.MulScalar(deltaT)

	return p.Speed.Add(&globalAcceleration)
}

func (p *Plane) calculateLift() float64 {

	// The p.inputsAxes.pitch should affect the amount of lift

	return p.Props.LiftRange.Min + -p.InputsAxes.Pitch*(p.Props.LiftRange.Max-p.Props.LiftRange.Min)
}

// calculateDrag calculate the amount of drag.
// Please note that the force here is expressed in newtons
func (p *Plane) calculateDrag() (drag Vector3D) {

	localSpeed := p.getLocalSpeed()

	airDensity := getAirDensity(p.Location.Y)

	// We inveres so the forces are applied in the right directions
	drag.X = -(p.Props.DragFactors.X * (localSpeed.X * localSpeed.X) * airDensity)
	drag.Y = -(p.Props.DragFactors.Y * (localSpeed.Y * localSpeed.Y) * airDensity)
	drag.Z = -(p.Props.DragFactors.Z * (localSpeed.Z * localSpeed.Z) * airDensity)

	return
}

func (p *Plane) calculateThrust() float64 {

	return p.InputThrust * p.Props.MaxThrust / p.Props.Mass

}

func (p *Plane) correctFromCollision() {

	const margin = 5

	triangle := p.Arena.arenaMap.OverredTriangle(p.Location)

	// We are out of bound
	if math.IsNaN(triangle[0].X) {
		return
	}

	// Small optimization
	if p.Location.Y >= highestInTriangle(triangle)+margin {
		return
	}

	// The real thing
	h := heightOnTriangle(p.Location, &triangle)

	// We are under the surface
	if p.Location.Y < h+margin {
		// We go back to the surface
		p.Location.Y = h + margin
		p.Speed.Y = 0
	}

}

func (p *Plane) getLocalSpeed() Vector3D {

	inverseOrientation := p.Orientation.getInverse()

	return p.Speed.multiplyByMatrix3(&inverseOrientation)

}

// getAirDensity returns the air density at the current altitude
func getAirDensity(altitude float64) float64 {

	// HACK: always returns 1.2 for now
	return 1.2
}
