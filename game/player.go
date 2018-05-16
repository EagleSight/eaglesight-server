package game

import (
	"log"

	"github.com/eaglesight/eaglesight-backend/world"
	"github.com/gorilla/websocket"
)

// Player : connected player's informations
type Player struct {
	conn      *websocket.Conn
	send      chan []byte
	input     chan world.PlayerInput
	deconnect chan *Player
	profile   PlayerProfile
}

// NewPlayer returns a new player
func NewPlayer(profile PlayerProfile, conn *websocket.Conn, deconect chan *Player) (player *Player) {

	player = &Player{
		conn:    conn,
		send:    make(chan []byte, 2),
		input:   make(chan world.PlayerInput),
		profile: profile,
	}
	go player.readPump()
	return player
}

func (p *Player) readPump() {

	for {
		_, message, err := p.conn.ReadMessage()

		if err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Check the opcode
		switch message[0] {
		case 0x3:
			p.input <- world.PlayerInput{UID: p.profile.UID, Data: message}
		}
	}
	p.deconnect <- p
}

// Close disconnect the player properly
func (p *Player) Close() {
	p.conn.Close()
}

// Write
func (p *Player) Write(message []byte) (n int, err error) {
	w, err := p.conn.NextWriter(websocket.BinaryMessage)

	if err != nil {
		p.deconnect <- p
		return 0, err
	}
	w.Write(message)

	if err := w.Close(); err != nil {
		p.deconnect <- p
		return 0, err
	}

	return len(message), nil
}
