package game

import (
	"testing"

	"github.com/eaglesight/eaglesight-server/world"
)

func TestNewPlayer(t *testing.T) {

	profile := PlayerProfile{}
	conn := dummyConn()

	player := NewPlayer(profile, conn)

	if player == nil {
		t.Fail()
	}
}

func TestListen(t *testing.T) {

	profile := PlayerProfile{UID: 2}
	conn := dummyConn()

	player := NewPlayer(profile, conn)

	input := make(chan world.PlayerInput, 16)
	exit := make(chan *Player, 1)

	go player.Listen(input, exit)

	conn.Send([]byte{0x3, 'h', 'e', 'l', 'l', 'o'})

	receivedInput := <-input

	if receivedInput.Data[0] != 0x3 {
		t.Fail()
	}

	if receivedInput.Data[1] != 'h' {
		t.Fail()
	}
}

func TestWrite(t *testing.T) {
	profile := PlayerProfile{UID: 2}
	conn := dummyConn()

	player := NewPlayer(profile, conn)

	player.Write([]byte{'H', 'e', 'l', 'l', 'o'})

	message := <-conn.conn

	if len(message) != 5 {
		t.Fail()
	}
}

func TestClose(t *testing.T) {
	profile := PlayerProfile{}
	conn := dummyConn()

	player := NewPlayer(profile, conn)

	err := player.Close()

	if err != nil {
		t.Fail()
	}
}
