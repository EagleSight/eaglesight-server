package main

import (
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
	username string
	token    uint64
	state    []byte
	arena    *Arena
	conn     *websocket.Conn
	//team *Team
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
			if message[0] == 0x3 {
				p.arena.state <- message
			}
		}

	}
}
