package main

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"
)

const (
	// UpdateDataLenght is uint32-uid|3*float32-position|3*float32-location
	UpdateDataLenght = 7 * 4
)

// Arena riens
type Arena struct {
	players        map[*Player]bool
	connect        chan *Player
	deconect       chan *Player
	input          chan *PlayerInput
	snapshotInputs map[uint32]*PlayerInput
	tick           uint32
	mux            sync.Mutex
}

func (a *Arena) broadcastSnapshots() {

	c := time.Tick(time.Second / 60)

	snapshotBuffer := new(bytes.Buffer)

	for now := range c {

		a.tick++

		// We lock the mutex because we want to make sure that nobody else append a state while the inputsPacket is made
		a.mux.Lock()

		binary.Write(snapshotBuffer, binary.BigEndian, uint8(0x3))
		binary.Write(snapshotBuffer, binary.BigEndian, uint32(a.tick))
		binary.Write(snapshotBuffer, binary.BigEndian, uint16(len(a.snapshotInputs)))

		for k, v := range a.snapshotInputs {

			v.plane.UpdateIntoBuffer(snapshotBuffer, v.data, now)

			a.snapshotInputs[k] = &PlayerInput{plane: v.plane, data: nil}

		}

		a.mux.Unlock()

		// Send inputs to all the players
		go a.Broadcast(snapshotBuffer.Bytes())

		// We reset the buffer, ready for the next tick
		snapshotBuffer.Reset()

	}
}

func newArena() *Arena {
	return &Arena{
		players:        make(map[*Player]bool),
		connect:        make(chan *Player),
		deconect:       make(chan *Player),
		input:          make(chan *PlayerInput),
		snapshotInputs: make(map[uint32]*PlayerInput),
		tick:           0,
	}
}

// Run start the Arena
func (a *Arena) Run() {

	go a.broadcastSnapshots()

	for {
		select {
		case player := <-a.connect:
			a.connectPlayer(player)
		case player := <-a.deconect:
			a.deconnectPlayer(player)
		case input := <-a.input:
			a.mux.Lock()
			a.snapshotInputs[input.plane.uid] = input
			a.mux.Unlock()
		}
	}
}

// Broadcast sent a []byte containing a payload to all the players
func (a *Arena) Broadcast(payload []byte) {

	// Send the payload to all the players
	for p := range a.players {
		p.send <- payload
	}

}

func (a *Arena) connectPlayer(player *Player) {

	player.sendPlayersList()

	a.players[player] = true

	a.snapshotInputs[player.uid] = &PlayerInput{plane: player.plane, data: nil}

	// 0x1 - player's uid ----
	message := make([]byte, 5)

	message[0] = 0x1 // Connection
	binary.BigEndian.PutUint32(message[1:], player.uid)

	go a.Broadcast(message)
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

		go a.Broadcast(message)
	}

	a.mux.Unlock()
}
