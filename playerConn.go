package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// PlayerInput : change to apply on a player's plane
type PlayerInput struct {
	plane *Plane
	data  []byte
}

// PlayerConn : connected player's informations
type PlayerConn struct {
	uid    uint8
	input  []byte
	conn   *websocket.Conn
	send   chan []byte
	plane  *Plane
	player Player
}

// NewPlayerConn returns a new player (TEST THIS?)
func NewPlayerConn(player Player, plane *Plane, conn *websocket.Conn) *PlayerConn {

	return &PlayerConn{
		conn:   conn,
		uid:    plane.UID,
		send:   make(chan []byte),
		plane:  plane,
		player: player,
	}
}

func (p *PlayerConn) connect(connect chan *PlayerConn) {
	connect <- p
}

// TEST THIS! (mocking a connection to close?)
func (p *PlayerConn) deconnect(deconnect chan *PlayerConn) {
	deconnect <- p
	p.conn.Close()
}

// TEST THIS!
func (p *PlayerConn) readPump(deconect chan *PlayerConn, input chan *PlayerInput) {

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
func (p *PlayerConn) writePump(deconect chan *PlayerConn) {

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
func (p *PlayerConn) listen(a *Arena) {

	go p.readPump(a.deconect, a.input)
	go p.writePump(a.deconect)

}
