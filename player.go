package main

import (
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
func NewPlayer(profil PlayerProfile, plane *Plane, conn *websocket.Conn) *Player {

	return &Player{
		conn:  conn,
		uid:   profil.Token,
		send:  make(chan []byte),
		plane: plane,
		name:  "player" + strconv.FormatUint(uint64(profil.Token), 10), // TODO: manage the name
	}
}

func (p *Player) connect(connect chan *Player) {
	connect <- p
}

// TEST THIS! (mocking a connection to close?)
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
func (p *Player) listen(a *Arena) {

	go p.readPump(a.deconect, a.input)
	go p.writePump(a.deconect)

}
