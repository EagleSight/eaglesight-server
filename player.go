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
	conn  *websocket.Conn
	send  chan []byte
	plane *Plane
	name  string
}

// TEST THIS!
// NewPlayer returns a new player
func NewPlayer(uid uint32, a *Arena, conn *websocket.Conn) *Player {

	return &Player{
		conn:  conn,
		uid:   uid,
		send:  make(chan []byte, 16),
		plane: NewPlane(uid, a.terrain),
		name:  "player" + strconv.FormatUint(uint64(uid), 10), // TODO: manage the name
	}
}

// TEST THIS!
func (p *Player) sendPlayersList(players *map[*Player]bool) {

	playersCount := len(*players)

	offset := 1 + 2

	message := make([]byte, offset+playersCount*4)

	message[0] = 0x4
	binary.BigEndian.PutUint16(message[1:], uint16(playersCount))

	for k := range *players {
		binary.BigEndian.PutUint32(message[offset:], k.uid)
		offset += 4
	}

	p.send <- message
}

// TEST THIS!
func (p *Player) connect(connect chan *Player) {
	connect <- p
}

// TEST THIS!
func (p *Player) deconnect(deconnect chan *Player) {
	deconnect <- p
	p.conn.Close()
}

// TEST THIS!
func (p *Player) readPump(deconect chan *Player, input chan *PlayerInput) {

	defer p.deconnect(deconect)

	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch message[0] {
		case 0x3:
			input <- &PlayerInput{plane: p.plane, data: message}
		}

	}
}

// TEST THIS!
func (p *Player) writePump(deconect chan *Player) {

	defer p.deconnect(deconect)

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

// TEST THIS! (How ?)
// Listen starts the message pumps
func (p *Player) Listen(a *Arena) {

	go p.readPump(a.deconect, a.input)
	go p.writePump(a.deconect)

}
