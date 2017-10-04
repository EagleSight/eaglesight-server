package main

import (
	"encoding/binary"
	"math"
	"time"
)

type vector3D struct {
	x, y, z float64
}

// Plane describe a plane with all its properties
type Plane struct {
	uid           uint32 // Plane owner's player uid
	yaw           int8   // from -127 to 127
	pitch         int8   // from -127 to 127
	location      vector3D
	rotation      vector3D
	thrust        uint8   // from 0 to 255
	speed         float64 // unit / seconds
	maxSpeed      float64
	maxYawSpeed   float64 // radian / seconds
	maxPitchSpeed float64 // radian / seconds
	maxRollSpeed  float64 // radian / seconds
	updatedLast   time.Time
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint32) *Plane {
	return &Plane{
		uid:   uid,
		yaw:   0,
		pitch: 0,
		location: vector3D{
			x: 0,
			y: 500,
			z: 0,
		},
		rotation: vector3D{
			x: 0,
			y: math.Pi,
			z: 0,
		},
		thrust:        0,
		speed:         0,
		maxSpeed:      20000,
		maxYawSpeed:   1.5,
		maxPitchSpeed: 0.5,
		maxRollSpeed:  1.0,
		updatedLast:   time.Now(),
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

		p.yaw = int8(params[1])
		p.pitch = int8(params[2])

		p.thrust = uint8(params[4])

	}
	// HACK: speed multiplied from thrust
	p.speed = float64(p.thrust) / 255 * p.maxSpeed

	mov := p.speed * deltaT

	p.rotation.y += p.maxYawSpeed * float64(p.yaw) / 127 * deltaT
	p.rotation.x += p.maxPitchSpeed * float64(p.pitch) / 127 * deltaT

	p.location.x += math.Sin(p.rotation.y) * mov
	p.location.z += math.Cos(p.rotation.y) * mov

	binary.BigEndian.PutUint32(buf[offset:], p.uid)

	binary.BigEndian.PutUint32(buf[offset+4:], math.Float32bits(float32(p.location.x)))
	binary.BigEndian.PutUint32(buf[offset+8:], math.Float32bits(float32(p.location.y)))
	binary.BigEndian.PutUint32(buf[offset+12:], math.Float32bits(float32(p.location.z)))

	binary.BigEndian.PutUint32(buf[offset+16:], math.Float32bits(float32(p.rotation.x)))
	binary.BigEndian.PutUint32(buf[offset+20:], math.Float32bits(float32(p.rotation.y)))
	binary.BigEndian.PutUint32(buf[offset+24:], math.Float32bits(float32(p.rotation.z)))

}
