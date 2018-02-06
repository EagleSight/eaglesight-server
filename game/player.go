package game

import (
	"log"

	"github.com/eaglesight/eaglesight-backend/world"
	"github.com/gorilla/websocket"
)

// Player : connected player's informations
type Player struct {
	conn    *websocket.Conn
	send    chan []byte
	input   chan world.PlaneInput
	profile PlayerProfile
}

// NewPlayer returns a new player
func NewPlayer(profile PlayerProfile, conn *websocket.Conn, deconect chan *Player) (player *Player) {

	player = &Player{
		conn:    conn,
		send:    make(chan []byte, 2),
		input:   make(chan world.PlaneInput),
		profile: profile,
	}

	go player.readPump(deconect)
	go player.writePump(deconect)

	return player
}

func (p *Player) readPump(deconect chan *Player) {

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
			if len(message) == 5 { // 0x3|Roll|Pitch|Yaw|Thrust
				// convert binary message to PlaneInput
				p.input <- world.PlaneInput{
					Roll:   -float64(int8(message[1])) / 127,
					Pitch:  float64(int8(message[2])) / 127,
					Yaw:    float64(int8(message[3])) / 127,
					Thrust: float64(uint8(message[4])) / 255,
				}

			}
		}

	}

	deconect <- p

}

func (p *Player) writePump(deconect chan *Player) {

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
