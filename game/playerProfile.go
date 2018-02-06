package game

import (
	"github.com/eaglesight/eaglesight-backend/world"
)

// PlayerProfile ...
type PlayerProfile struct {
	Name  string           `json:"username"`
	UUID  string           `json:"accessKey"`
	UID   uint8            `json:"uid"`
	Model world.PlaneModel `json:"planeModel"`
}
