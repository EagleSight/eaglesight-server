package world

import (
	"log"
	"math"
	"sync"

	"github.com/eaglesight/eaglesight-backend/mathutils"
)

// PlaneInput ...
type PlaneInput struct {
	Roll, Pitch, Yaw, Thrust float64
}

// PlaneModel are all the constant properties that can easily be loaded from a JSON object
type PlaneModel struct {
	MaxThrust    float64            `json:"maxThrust"`
	Mass         float64            `json:"mass"`
	MaxRotations mathutils.Vector3D `json:"maxRotations"` // All in radians / seconds
	DragFactors  mathutils.Vector3D `json:"dragFactors"`  // A drag factor are "Area * Drag Coeficient * 0.5" of a side
	LiftMin      float64            `json:"liftMin"`
	LiftMax      float64            `json:"liftMax"`
	DefaultSpeed float64            `json:"defaultSpeed"` // Default speed of the plane on the Z axis
}

// PlaneState reperent the LocRot of a plane
type PlaneState struct {
	UID      uint8
	Location mathutils.Vector3D // Absolute Location in the world
	Rotation mathutils.Vector3D // Absolute rotation of the plane (Radians, Euler's angles)
}

// Plane describe a plane with all its properties
type Plane struct {
	UID         uint8
	input       PlaneInput
	model       PlaneModel
	location    mathutils.Vector3D // Absolute Location in the world
	speed       mathutils.Vector3D // unit / seconds
	orientation mathutils.Matrix3
	isNoMore    bool
	inputs      chan PlaneInput
	mux         sync.Mutex
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint8, model PlaneModel) (plane *Plane) {
	plane = &Plane{
		UID: uid,
		input: PlaneInput{
			Roll:   0,
			Pitch:  0,
			Yaw:    0,
			Thrust: 0,
		},
		location: mathutils.Vector3D{
			X: 0,
			Y: 1500, // Find a way to set location
			Z: 0,
		},
		orientation: mathutils.NewMatrix3(),
		speed: mathutils.Vector3D{
			X: 0,
			Y: 0,
			Z: model.DefaultSpeed,
		},
		model:    model,
		isNoMore: false,
		inputs:   make(chan PlaneInput),
	}

	go func(p *Plane) {
		for input := range p.inputs {
			p.mux.Lock()
			p.input = input
			p.mux.Unlock()
		}

		p.mux.Lock()
		p.isNoMore = true
		p.mux.Unlock()
	}(plane)

	return plane
}

// Tick updates the plane's properties from new parameters
func (p *Plane) Tick(deltaT float64, terrain *Terrain) (state PlaneState) {

	p.mux.Lock()

	// Update the rotation
	p.orientation = p.calculateRotation(deltaT)

	// Update the speed
	p.speed = p.calculateSpeed(deltaT)

	p.location = p.location.Add(p.speed.MulScalar(deltaT))

	p.CorrectFromCollision(terrain)

	state = PlaneState{
		UID:      p.UID,
		Location: p.location,
		Rotation: p.orientation.ToEulerAngle(),
	}

	p.mux.Unlock()

	return state

}

func (p *Plane) calculateRotation(deltaT float64) mathutils.Matrix3 {

	// Generate the matrices that represent the rotation change
	pitchMat := mathutils.MakeMatrix3X(p.model.MaxRotations.X * p.input.Pitch * deltaT)
	yawMat := mathutils.MakeMatrix3Y(p.model.MaxRotations.Y * p.input.Yaw * deltaT)
	rollMat := mathutils.MakeMatrix3Z(p.model.MaxRotations.Z * p.input.Roll * deltaT)

	// Multiply them together in the right order
	localRotMat := yawMat.Mul(pitchMat)
	localRotMat = rollMat.Mul(localRotMat)

	return p.orientation.Mul(localRotMat)
}

func (p *Plane) calculateSpeed(deltaT float64) mathutils.Vector3D {

	localAcceleration := mathutils.Vector3D{
		X: 0,
		Y: p.calculateLift(),
		Z: p.calculateThrust(),
	}

	// Compute the drag
	localDrag := p.calculateDrag()

	// Divide the force by the mass
	localDrag = localDrag.DivScalar(p.model.Mass)

	// Add the drag to the local acceleration
	localAcceleration = localAcceleration.Add(localDrag)

	// Convert to global
	globalAcceleration := localAcceleration.MultiplyByMatrix3(p.orientation)

	log.Println("Global Acceleration Y:", globalAcceleration.Y)

	// Apply gravity
	globalAcceleration.Y += -9.8

	return p.speed.Add(globalAcceleration.MulScalar(deltaT))
}

func (p *Plane) calculateLift() float64 {

	return p.model.LiftMin + -p.input.Pitch*(p.model.LiftMax-p.model.LiftMin)
}

// calculateDrag calculate the amount of drag.
// Please note that the force here is expressed in newtons
func (p *Plane) calculateDrag() (drag mathutils.Vector3D) {

	localSpeed := p.getLocalSpeed()

	airDensity := getAirDensity(p.location.Y)

	// We inveres so the forces are applied in the right directions
	drag.X = -(p.model.DragFactors.X * (localSpeed.X * localSpeed.X) * airDensity)
	drag.Y = -(p.model.DragFactors.Y * (localSpeed.Y * localSpeed.Y) * airDensity)
	drag.Z = -(p.model.DragFactors.Z * (localSpeed.Z * localSpeed.Z) * airDensity)

	return
}

func (p *Plane) calculateThrust() float64 {
	return (p.input.Thrust * p.model.MaxThrust) / p.model.Mass
}

// CorrectFromCollision update the position of the plane if there is a collision with the terrain
func (p *Plane) CorrectFromCollision(terrain *Terrain) {

	const margin = 5

	triangle := terrain.OverredTriangle(p.location)

	// We are out of bound
	if math.IsNaN(triangle[0].X) {
		return
	}

	// Small optimization
	if p.location.Y >= mathutils.HighestInTriangle(triangle)+margin {
		return
	}

	// The real thing
	h := mathutils.HeightOnTriangle(p.location, &triangle)

	// We are under the surface
	if p.location.Y < h+margin {
		// We go back to the surface
		p.location.Y = h + margin
		p.speed.Y = 0
	}

}

func (p *Plane) getLocalSpeed() mathutils.Vector3D {

	return p.speed.MultiplyByMatrix3(p.orientation.GetInverse())

}

// getAirDensity returns the air density at the current altitude
func getAirDensity(altitude float64) float64 {

	// HACK: always returns 1.2
	return 1.2
}

func (p *Plane) isDead() bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.isNoMore
}
