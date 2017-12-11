package main

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"sync"
	"time"
)

const (
	// UpdateDataLenght is uint32-uid|3*float32-position|3*float32-location
	UpdateDataLenght = 7 * 4
)

// Arena riens
type Arena struct {
	gameID          string
	players         map[*Player]bool
	connect         chan *Player
	deconect        chan *Player
	input           chan *PlayerInput
	everybody       chan []byte
	snapshotInputs  map[uint32]*PlayerInput
	tick            uint32
	mux             sync.Mutex
	terrain         *Terrain
	PlayersProfiles map[uint32]PlayerProfile
}

// TEST THIS!
// NewArena return a arena with default settings
func NewArena(params GameParameters, terrain *Terrain) *Arena {

	// Put all the registered players in a map
	playersProfiles := make(map[uint32]PlayerProfile)

	for _, rp := range params.Players {
		playersProfiles[rp.Token] = rp
	}

	return &Arena{
		gameID:          params.GameID,
		players:         make(map[*Player]bool),
		connect:         make(chan *Player),
		deconect:        make(chan *Player),
		input:           make(chan *PlayerInput),
		everybody:       make(chan []byte, 2),
		snapshotInputs:  make(map[uint32]*PlayerInput),
		tick:            0,
		terrain:         terrain,
		PlayersProfiles: playersProfiles,
	}
}

func (a *Arena) TakePlayerProfile(uid uint32) (profile PlayerProfile, err error) {

	a.mux.Lock()
	defer a.mux.Unlock()

	if _, ok := a.PlayersProfiles[uid]; ok {
		profile = a.PlayersProfiles[uid]
		delete(a.PlayersProfiles, uid)
		return profile, nil
	}

	return profile, errors.New("Unauthorized player")

}

// TEST THIS! (How ?)
func generateSnapshot(a *Arena, deltaT float64) []byte {

	// We lock the mutex because we want to make sure that nobody else append a state while the inputsPacket is made
	a.mux.Lock()

	a.tick++

	offset := 1 + 4 + 2 // unt8 + uint32 + uint16
	const playerDataLenght = 28
	snapshot := make([]byte, offset+len(a.snapshotInputs)*playerDataLenght)

	snapshot[0] = uint8(3)
	binary.BigEndian.PutUint32(snapshot[1:], uint32(a.tick))
	binary.BigEndian.PutUint16(snapshot[5:], uint16(len(a.snapshotInputs)))

	for k, v := range a.snapshotInputs {

		v.plane.Update(v.data, deltaT)

		// Dump everything into the slice
		binary.BigEndian.PutUint32(snapshot[offset:], v.plane.UID)

		// Location
		binary.BigEndian.PutUint32(snapshot[offset+4:], math.Float32bits(float32(v.plane.Location.X)))
		binary.BigEndian.PutUint32(snapshot[offset+8:], math.Float32bits(float32(v.plane.Location.Y)))
		binary.BigEndian.PutUint32(snapshot[offset+12:], math.Float32bits(float32(v.plane.Location.Z)))

		// Rotation
		binary.BigEndian.PutUint32(snapshot[offset+16:], math.Float32bits(float32(v.plane.Rotation.X)))
		binary.BigEndian.PutUint32(snapshot[offset+20:], math.Float32bits(float32(v.plane.Rotation.Y)))
		binary.BigEndian.PutUint32(snapshot[offset+24:], math.Float32bits(float32(v.plane.Rotation.Z)))

		offset += playerDataLenght

		a.snapshotInputs[k] = &PlayerInput{plane: v.plane, data: nil}

	}

	a.mux.Unlock()

	return snapshot

}

func (a *Arena) broadcastSnapshots() {

	previousTickTime := time.Now()

	c := time.Tick(time.Second / 60)

	for now := range c {

		// Calculate the time since the last time updated
		deltaT := now.Sub(previousTickTime).Seconds()

		// Save for the next time
		previousTickTime = now

		// Send inputs to all the players
		a.everybody <- generateSnapshot(a, deltaT)

	}
}

// Run start the Arena
func (a *Arena) Run() {

	go a.broadcastPump()

	go a.broadcastSnapshots()

	for {
		select {
		case player := <-a.connect:
			a.sendPlayersList(player)
			a.connectPlayer(player)
		case player := <-a.deconect:
			a.deconnectPlayer(player)
		case input := <-a.input:
			a.mux.Lock()
			a.snapshotInputs[input.plane.UID] = input
			a.mux.Unlock()
		}
	}
}

// Broadcast sent a []byte containing a payload to all the players
func (a *Arena) broadcastPump() {

	// Send the payload to all the players
	for snapshot := range a.everybody {
		for p := range a.players {
			p.send <- snapshot
		}
	}

}

// TEST THIS!
// Sends the list of all the connected players, including "player" itself
func (a *Arena) sendPlayersList(player *Player) {

	playersCount := len(a.players)

	offset := 1 + 2

	message := make([]byte, offset+playersCount*4)

	message[0] = 0x4
	binary.BigEndian.PutUint16(message[1:], uint16(playersCount))

	for k := range a.players {
		binary.BigEndian.PutUint32(message[offset:], k.uid)
		offset += 4
	}

	player.send <- message
}

func (a *Arena) connectPlayer(player *Player) {

	a.mux.Lock()

	a.snapshotInputs[player.uid] = &PlayerInput{plane: player.plane, data: nil}
	a.players[player] = true

	a.mux.Unlock()

	// 0x1 + player's uid
	message := make([]byte, 5)

	message[0] = 0x1 // Connection
	binary.BigEndian.PutUint32(message[1:], player.uid)

	a.everybody <- message

	log.Println(player.name + " connected.")
}

func (a *Arena) deconnectPlayer(player *Player) {

	a.mux.Lock()

	// Remove the player from the players list
	if _, ok := a.players[player]; ok {
		delete(a.players, player)
		delete(a.snapshotInputs, player.uid)

		message := make([]byte, 5)

		message[0] = 0x2 // Deconnection
		binary.BigEndian.PutUint32(message[1:], player.uid)

		a.everybody <- message
	}

	a.mux.Unlock()

	log.Println(player.name + " deconnected.")
}
