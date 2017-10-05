package main

import (
	"encoding/binary"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

// PlayerInput : change to apply on a player's plane
type PlayerInput struct {
	plane *Plane
	data  []byte
}

// Player : connected player's informations
type Player struct {
	uid   uint32
	input []byte
	arena *Arena
	conn  *websocket.Conn
	send  chan []byte
	plane *Plane
}

// NewPlayer returns a new player
func NewPlayer(uid uint32, a *Arena, conn *websocket.Conn) *Player {
	return &Player{
		arena: a,
		conn:  conn,
		uid:   uid,
		send:  make(chan []byte, 16),
		plane: NewPlane(uid),
	}
}

func (p *Player) sendPlayersList() {

	playersCount := len(p.arena.players)

	offset := 1 + 2

	message := make([]byte, offset+playersCount*4)

	message[0] = 0x4
	binary.BigEndian.PutUint16(message[1:], uint16(playersCount))

	for k := range p.arena.players {
		binary.BigEndian.PutUint32(message[offset:], k.uid)
		offset += 4
	}

	p.send <- message
}

func (p *Player) deconnect() {
	log.Println(strconv.FormatUint(uint64(p.uid), 10) + " deconnected.")
	p.arena.deconect <- p
	p.conn.Close()
}

func (p *Player) readPump() {

	defer p.deconnect()

	p.arena.connect <- p

	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		if message[0] == 0x3 { // Plane's state update
			p.arena.input <- &PlayerInput{plane: p.plane, data: message}
		}

	}
}

func (p *Player) writePump() {

	defer p.deconnect()

	for message := range p.send {

		w, err := p.conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			return
		}

		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}
