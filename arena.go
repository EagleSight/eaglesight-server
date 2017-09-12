package main

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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

	c := time.Tick(time.Second / 60)

	for range c {

		if len(a.updates) > 0 {

			// We lock the mutex because we want to make sure that nobody else append a state in updates while they are sent
			a.mux.Lock()

			merge := []byte{}

			for k, v := range a.updates {
				merge = append(merge, v...)
				delete(a.updates, k)
			}

			// Now that we flushed all the state changes, we reset the updates buffer
			a.mux.Unlock()

			// Send updates to all the players
			for p := range a.players {
				p.conn.WriteMessage(websocket.BinaryMessage, merge)
			}
		}

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

func (a *Arena) run() {

	go a.broadcastStates()

	for {
		select {
		case player := <-a.register:
			a.players[player] = true
		case player := <-a.unregister:
			if _, ok := a.players[player]; ok {
				delete(a.players, player)
			}
		case state := <-a.state:
			a.updates[binary.LittleEndian.Uint32(state[1:5])] = state
		}
	}
}
