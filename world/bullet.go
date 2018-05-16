package world

import (
	"github.com/eaglesight/eaglesight-backend/mathutils"
)

// Bullet represent a bullet
type Bullet struct {
	source      uint8              // UID of the player that shot the bullet
	location    mathutils.Vector3D // Location in global space
	speed       mathutils.Vector3D // Speed in global space
	damage      uint8              // Amount of damage the bullet make on impact
	ticksToLive uint16
}

// NewBullet create a new bullet
func NewBullet(source uint8, origin mathutils.Vector3D, direction *mathutils.Matrix3, speed float64, damage uint8) *Bullet {

	const BulletLifetime = 500 // Bullet will live for 500 updates. That might be tweaked in the future
	speedVector := mathutils.Vector3D{X: 0, Y: speed, Z: 0}

	return &Bullet{
		source:      source,
		location:    origin,
		speed:       speedVector.MultiplyByMatrix3(direction),
		damage:      damage,
		ticksToLive: BulletLifetime,
	}
}

// Update updates the state of the bullet. Returns whether the bullet is still "living"
func (b *Bullet) Update(deltaT float64) bool {

	// Apply some gravity
	b.speed.Add(mathutils.Vector3D{X: 0, Y: 9.8 * deltaT, Z: 0})
	// One tick closer to death...
	b.ticksToLive--
	return b.ticksToLive > 0
}
