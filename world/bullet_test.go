package world

import (
	"testing"

	"github.com/eaglesight/eaglesight-backend/mathutils"
)

func TestNewBullet(t *testing.T) {
	origin := mathutils.Vector3D{
		X: 4,
		Y: 3,
		Z: 2,
	}
	direction := mathutils.NewMatrix3()
	bullet := NewBullet(3, origin, &direction, 400, 12)

	if bullet.speed.Y != 400 {
		t.Fail()
	}
}

func TestBulletUpdate(t *testing.T) {
	origin := mathutils.Vector3D{
		X: 0,
		Y: 0,
		Z: 0,
	}
	direction := mathutils.NewMatrix3()
	bullet := NewBullet(3, origin, &direction, 400, 12)

	if bullet.location.Y != 0 {
		t.Errorf("Bullet's location is %s", bullet.location)
	}

	bullet.Update(1)

	if bullet.location.Y == 0 {
		t.Errorf("Bullet's location is %s", bullet.location)
	}
}
