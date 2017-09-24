package main

import (
	"encoding/binary"
	"log"

	"github.com/gorilla/websocket"
)

// type Vector3D struct {
// 	x, y, z float64
// }

//
// type PlayerState struct {
// 	livesLeft           int8
// 	location, rotation  Vector3D
// 	localAcceleration   Vector3D
// 	angularAcceleration Vector3D
// 	isOut, isFiring     bool
// }

// Player : connected player's informations
type Player struct {
	uid   uint32
	state []byte
	arena *Arena
	conn  *websocket.Conn
	send  chan []byte
}

func (p *Player) sendPlayersList() {

	playersCount := uint16(len(p.arena.players))

	message := make([]byte, 1+2+playersCount*4)

	binary.BigEndian.PutUint16(message[1:], playersCount)

	message[0] = 0x4

	offset := 0

	for k := range p.arena.players {
		binary.BigEndian.PutUint32(message[1+2+offset:], k.uid)
		offset += 4
	}

	p.send <- message
}

func (p *Player) readPump() {

	defer func() {
		p.arena.unregister <- p
		p.conn.Close()
	}()

	p.arena.register <- p

	// c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		if len(message) > 0 {
			if message[0] == 0x3 { // Plane's state update
				p.arena.state <- message
			}
		}

	}
}

func (p *Player) writePump() {

	defer func() {
		p.conn.Close()
	}()

	for {
		select {
		case message := <-p.send:

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
}
