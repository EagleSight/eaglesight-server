package main

import (
	"encoding/binary"
	"sync"
	"time"
)

// Arena riens
type Arena struct {
	players    map[*Player]bool
	register   chan *Player
	unregister chan *Player
	state      chan []byte
	updates    map[uint32][]byte
	mux        sync.Mutex
}

func (a *Arena) broadcastStates() {

	c := time.Tick(time.Second / 30)

	for range c {

		// We lock the mutex because we want to make sure that nobody else append a state in updates while they are sent
		a.mux.Lock()

		if len(a.updates) > 0 {

			for k, v := range a.updates {

				// Send updates to all the players
				for p := range a.players {
					p.send <- v
				}
				delete(a.updates, k)
			}
		}

		a.mux.Unlock()

	}
}

func newArena() *Arena {
	return &Arena{
		players:    make(map[*Player]bool),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		state:      make(chan []byte),
		updates:    make(map[uint32][]byte),
	}
}

// Run start the Arena
func (a *Arena) Run() {

	go a.broadcastStates()

	for {
		select {
		case player := <-a.register:
			a.connectPlayer(player)
		case player := <-a.unregister:
			a.deconnectPlayer(player)
		case state := <-a.state:
			a.mux.Lock()
			a.updates[binary.LittleEndian.Uint32(state[1:5])] = state
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

	// 0x1 - player's uid ----
	message := make([]byte, 5)

	message[0] = 0x1 // Connection
	binary.LittleEndian.PutUint32(message[1:], player.uid)

	go a.Broadcast(message)
}

func (a *Arena) deconnectPlayer(player *Player) {

	// Remove the player from the players list
	if _, ok := a.players[player]; ok {
		delete(a.players, player)
	}

	message := make([]byte, 5)

	message[0] = 0x2 // Deconnection
	binary.LittleEndian.PutUint32(message[1:], player.uid)

	go a.Broadcast(message)
}
