package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"
)

type vector3DF32 struct {
	x, y, z float32
}

type vector3DF64 struct {
	x, y, z float64
}

// Plane describe a plane with all its properties
type Plane struct {
	uid             uint32 // Plane owner's player uid
	yaw             int8   // from -127 to 127
	location        vector3DF32
	rotation        vector3DF64
	thrust          uint8   // from 0 to 255
	speed           float64 // unit / seconds
	maxSpeed        float64
	maxAngularSpeed float64 // radian / seconds
	updatedLast     time.Time
}

// NewPlane fill the plane with its default properties
func NewPlane(uid uint32) *Plane {
	return &Plane{
		uid: uid,
		yaw: 0,
		location: vector3DF32{
			x: 0,
			y: 500,
			z: 0,
		},
		rotation: vector3DF64{
			x: 0,
			y: 3.1416,
			z: 0,
		},
		thrust:          0,
		speed:           0,
		maxSpeed:        20000,
		maxAngularSpeed: 1.5,
		updatedLast:     time.Now(),
	}
}

// UpdateIntoBuffer updates the plane's properties from new parameters
// and puts them into a buffer (first arg)
func (p *Plane) UpdateIntoBuffer(buf *bytes.Buffer, params []byte, tick time.Time) {

	// Calculate the time since the last time updated
	deltaT := tick.Sub(p.updatedLast).Seconds()

	// Set updatedLast to the current tick
	p.updatedLast = tick

	if len(params) > 0 { // We update those only if we have data

		p.yaw = int8(params[1])
		p.thrust = uint8(params[2])

	}
	// HACK: speed multiplied from thrust
	p.speed = float64(p.thrust) / 255 * p.maxSpeed

	mov := p.speed * deltaT

	rot := p.maxAngularSpeed * float64(p.yaw) / 127 * deltaT
	p.rotation.y += rot

	p.location.x += float32(math.Sin(p.rotation.y) * mov)
	p.location.z += float32(math.Cos(p.rotation.y) * mov)

	binary.Write(buf, binary.BigEndian, p.uid)

	binary.Write(buf, binary.BigEndian, p.location.x)
	binary.Write(buf, binary.BigEndian, p.location.y)
	binary.Write(buf, binary.BigEndian, p.location.z)

	binary.Write(buf, binary.BigEndian, float32(p.rotation.x))
	binary.Write(buf, binary.BigEndian, float32(p.rotation.y))
	binary.Write(buf, binary.BigEndian, float32(p.rotation.z))

}
