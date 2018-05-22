package world

import (
	"encoding/binary"
	"math"

	"github.com/eaglesight/eaglesight-server/mathutils"
)

const (
	// PlaneSnapshotSize : uint8 (planeId) + float32 * 3 (location) + float32 * 4 (rotation) + 1 bit for firing + 7 bits damage
	PlaneSnapshotSize = 1 + 1 + (3 * 4) + (4 * 4)
)

// PlaneInput ...
type PlaneInput struct {
	Roll, Pitch, Yaw, Thrust float64
	IsFiring                 bool
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

// Plane describe a plane with all its properties
type Plane struct {
	UID                uint8
	input              PlaneInput
	model              PlaneModel
	location           mathutils.Vector3D // Absolute Location in the world
	speed              mathutils.Vector3D // unit / seconds
	orientation        mathutils.Matrix3
	orientationInverse mathutils.Matrix3
	isNoMore           bool
	gun                chan<- Bullet
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint8, model PlaneModel, gun chan<- Bullet) (plane *Plane) {

	plane = &Plane{
		UID: uid,
		input: PlaneInput{
			Roll:     0,
			Pitch:    0,
			Yaw:      0,
			Thrust:   0,
			IsFiring: false,
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
		gun:      gun,
	}
	return plane
}

func (p *Plane) Write(data []byte) (n int, err error) {

	if len(data) == 6 { // 0x3|Roll|Pitch|Yaw|Thrust|(IsFiring|...)
		// convert binary message to PlaneInput
		p.input = PlaneInput{
			Roll:     -float64(int8(data[1])) / 127,
			Pitch:    float64(int8(data[2])) / 127,
			Yaw:      float64(int8(data[3])) / 127,
			Thrust:   float64(uint8(data[4])) / 255,
			IsFiring: data[5] >= 0x80,
		}
		return len(data), nil
	}
	// Nothing matched
	return 0, nil
}

func (p *Plane) Read(snapshot []byte) (n int, err error) {
	// UID
	snapshot[0] = p.UID
	// TODO: Dammage
	snapshot[1] = 0
	// Location
	binary.BigEndian.PutUint32(snapshot[2:], math.Float32bits(float32(p.location.X)))
	binary.BigEndian.PutUint32(snapshot[6:], math.Float32bits(float32(p.location.Y)))
	binary.BigEndian.PutUint32(snapshot[10:], math.Float32bits(float32(p.location.Z)))
	// Rotation
	rotation := p.orientation.ToQuaternion()
	binary.BigEndian.PutUint32(snapshot[14:], math.Float32bits(float32(rotation.X)))
	binary.BigEndian.PutUint32(snapshot[18:], math.Float32bits(float32(rotation.Y)))
	binary.BigEndian.PutUint32(snapshot[22:], math.Float32bits(float32(rotation.Z)))
	binary.BigEndian.PutUint32(snapshot[26:], math.Float32bits(float32(rotation.W)))
	return PlaneSnapshotSize, nil
}

// Update updates the plane's properties from new parameters
func (p *Plane) Update(deltaT float64, terrain *Terrain) {
	// Update the rotation
	p.orientation = p.calculateRotation(deltaT)
	// Update the speed
	p.speed = p.calculateSpeed(deltaT)
	p.location = p.location.Add(p.speed.MulScalar(deltaT))
	p.CorrectFromCollision(terrain)
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
	globalAcceleration := localAcceleration.MultiplyByMatrix3(&p.orientation)
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
	if p.location.Y >= mathutils.HighestInTriangle(&triangle)+margin {
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
	p.orientation.Inverse(&p.orientationInverse)
	return p.speed.MultiplyByMatrix3(&p.orientationInverse)
}

// getAirDensity returns the air density at the current altitude
func getAirDensity(altitude float64) float64 {
	// HACK: always returns 1.2
	return 1.2
}

func (p *Plane) isDead() bool {

	return p.isNoMore
}
