package game

import (
	"github.com/eaglesight/eaglesight-backend/world"
)

// Player : connected player's informations
type Player struct {
	conn      PlayerConn
	input     chan world.PlayerInput
	deconnect chan *Player
	profile   PlayerProfile
}

// PlayerProfile ...
type PlayerProfile struct {
	Name  string           `json:"username"`
	UUID  string           `json:"accessKey"`
	UID   uint8            `json:"uid"`
	Model world.PlaneModel `json:"planeModel"`
}

// PlayerConn represent an connection to a player
type PlayerConn interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
}

// NewPlayer returns a new player
func NewPlayer(profile PlayerProfile, conn PlayerConn, deconect chan *Player) (player *Player) {

	player = &Player{
		conn:      conn,
		input:     make(chan world.PlayerInput),
		deconnect: make(chan *Player),
		profile:   profile,
	}
	go player.listen()
	return player
}

func (p *Player) listen() {

	for {
		message, err := p.conn.Receive()

		if err != nil {
			p.Close()
			break
		}

		// Check the opcode
		switch message[0] {
		case 0x3:
			p.input <- world.PlayerInput{UID: p.profile.UID, Data: message}
		}
	}
}

// Close disconnect the player properly
func (p *Player) Close() error {
	p.deconnect <- p
	return p.conn.Close()
}

// Write
func (p *Player) Write(message []byte) (n int, err error) {
	err = p.conn.Send(message)

	if err != nil {
		p.Close()
		return 0, err
	}
	return len(message), nil
}
