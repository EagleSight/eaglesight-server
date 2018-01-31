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
	gameID         string
	connections    map[*PlayerConn]bool
	connect        chan *PlayerConn
	deconect       chan *PlayerConn
	input          chan *PlayerInput
	everybody      chan []byte
	snapshotInputs map[uint8]*PlayerInput
	mux            sync.Mutex
	terrain        *Terrain
	players        map[string]Player
	uniqueCount    uint8
}

// NewArena return a arena with default settings (TEST THIS!)
func NewArena(params GameParameters, terrain *Terrain) *Arena {

	// Put all the registered players in a map
	players := make(map[string]Player)

	for _, rp := range params.Players {
		players[rp.UUID] = rp
	}

	return &Arena{
		gameID:         params.GameID,
		connections:    make(map[*PlayerConn]bool),
		connect:        make(chan *PlayerConn),
		deconect:       make(chan *PlayerConn),
		input:          make(chan *PlayerInput),
		everybody:      make(chan []byte, 2),
		snapshotInputs: make(map[uint8]*PlayerInput),
		terrain:        terrain,
		players:        players,
		uniqueCount:    0,
	}
}

// ValidatePlayer
func (a *Arena) ValidatePlayer(uuid string) (player Player, uid uint8, err error) {

	a.mux.Lock()
	defer a.mux.Unlock()

	if _, ok := a.players[uuid]; ok {
		player = a.players[uuid]
		delete(a.players, uuid)
		a.uniqueCount++
		return player, a.uniqueCount, nil
	}

	return player, 0, errors.New("Unauthorized player: Unknown UUID")

}

// TEST THIS! (How ?)
func generateSnapshot(a *Arena, deltaT float64) []byte {

	// We lock the mutex because we want to make sure that nobody else append a state while the inputsPacket is made
	a.mux.Lock()

	offset := 1 + 2 // unt8 + uint16
	const playerDataLenght = 25
	snapshot := make([]byte, offset+len(a.snapshotInputs)*playerDataLenght)

	snapshot[0] = uint8(0x3)
	snapshot[1] = uint8(len(a.snapshotInputs))

	for k, v := range a.snapshotInputs {

		v.plane.Update(v.data, deltaT)

		// Dump everything into the slice

		// UID
		snapshot[offset] = v.plane.UID

		// Location
		binary.BigEndian.PutUint32(snapshot[offset+1:], math.Float32bits(float32(v.plane.Location.X)))
		binary.BigEndian.PutUint32(snapshot[offset+5:], math.Float32bits(float32(v.plane.Location.Y)))
		binary.BigEndian.PutUint32(snapshot[offset+9:], math.Float32bits(float32(v.plane.Location.Z)))

		// Rotation
		binary.BigEndian.PutUint32(snapshot[offset+13:], math.Float32bits(float32(v.plane.Rotation.X)))
		binary.BigEndian.PutUint32(snapshot[offset+17:], math.Float32bits(float32(v.plane.Rotation.Y)))
		binary.BigEndian.PutUint32(snapshot[offset+21:], math.Float32bits(float32(v.plane.Rotation.Z)))

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
		case playerConn := <-a.connect:
			a.sendPlayersList(playerConn)
			a.connectPlayer(playerConn)
		case playerConn := <-a.deconect:
			a.deconnectPlayer(playerConn)
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
		for p := range a.connections {
			p.send <- snapshot
		}
	}

}

// TEST THIS!
// Sends the list of all the connected players, including "player" itself
func (a *Arena) sendPlayersList(conn *PlayerConn) {

	playersCount := len(a.connections)

	offset := 1 + 1

	message := make([]byte, offset+playersCount)

	message[0] = 0x4
	message[1] = uint8(playersCount)

	for k := range a.connections {
		message[offset] = k.uid
		offset++
	}

	conn.send <- message
}

func (a *Arena) connectPlayer(conn *PlayerConn) {

	a.mux.Lock()

	a.snapshotInputs[conn.uid] = &PlayerInput{plane: conn.plane, data: nil}
	a.connections[conn] = true

	a.mux.Unlock()

	// 0x1 + player's uid
	message := make([]byte, 5)

	message[0] = 0x1 // Connection
	message[1] = conn.uid

	a.everybody <- message

	log.Println(conn.player.Name + " connected.")
}

func (a *Arena) deconnectPlayer(conn *PlayerConn) {

	a.mux.Lock()

	// Remove the player from the players list
	if _, ok := a.connections[conn]; ok {
		delete(a.connections, conn)
		delete(a.snapshotInputs, conn.uid)

		//
		a.players[conn.player.UUID] = conn.player

		message := make([]byte, 5)

		message[0] = 0x2 // Deconnection
		message[1] = conn.uid

		a.everybody <- message
	}

	a.mux.Unlock()

	log.Println(conn.player.Name + " deconnected.")
}
