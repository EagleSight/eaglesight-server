package game

import (
	"testing"

	"github.com/eaglesight/eaglesight-server/world"
)

func dummyParams() (params Parameters) {
	profiles := []PlayerProfile{}

	profiles = append(profiles, PlayerProfile{
		Name:  "pako_panda",
		UUID:  "pako",
		UID:   0,
		Model: world.PlaneModel{},
	})

	params = Parameters{
		GameID:  "abcdefg12345BBQ",
		Players: profiles,
	}

	return params
}

func dummyServer() *Server {
	params := dummyParams()

	return NewServer(params)
}

func TestNewServer(t *testing.T) {
	params := dummyParams()
	server := dummyServer()

	p1 := params.Players[0]

	if server.profiles[p1.UUID].UUID != p1.UUID {
		t.Fail()
	}

	if server.profiles[p1.UUID] != p1 {
		t.Fail()
	}

}

// --------------------------------------------------

type dummyPlayerConn struct {
	conn chan []byte
}

func dummyConn() *dummyPlayerConn {
	return &dummyPlayerConn{
		conn: make(chan []byte, 16),
	}
}

// Receive ...
func (c *dummyPlayerConn) Receive() (data []byte, err error) {
	return <-c.conn, err
}

func (c *dummyPlayerConn) Send(message []byte) error {
	c.conn <- message
	return nil
}

func (c *dummyPlayerConn) Close() error {
	return nil
}

// ------------------------------------------------------

func TestConnect(t *testing.T) {

	server := dummyServer()

	// change the connect for a buffered channel
	server.connect = make(chan *Player, 1)

	profile := server.profiles["pako"]

	conn := dummyConn()

	server.Connect(conn, profile)

	p := <-server.connect

	if p.profile != server.profiles["pako"] {
		t.Fail()
	}

}

func TestBroadcastMessage(t *testing.T) {

	server := dummyServer()
	profile := server.profiles["pako"]

	message := []byte{0x4, 0x4}

	conn := dummyConn()

	server.connectedPlayers[profile.UID] = NewPlayer(profile, conn)

	server.broadcastMessage(message)

	message2, _ := conn.Receive()

	if message[0] != message2[0] {
		t.Fail()
	}
}

func TestPlayersListMessage(t *testing.T) {

	server := dummyServer()
	conn := dummyConn()
	profile := server.profiles["pako"]
	server.connectedPlayers[profile.UID] = NewPlayer(profile, conn)

	list := server.playersListMessage(1)

	if list[1] != 1 || list[2] != profile.UID || len(list) != 3 {
		t.Fail()
	}
}

func TestConnectionMessage(t *testing.T) {

	message := connectionMessage(2)

	if message[0] != 0x1 || message[1] != 2 {
		t.Fail()
	}
}

func TestDeconnectionMessage(t *testing.T) {

	message := deconnectionMessage(3)

	if message[0] != 0x2 || message[1] != 3 {
		t.Fail()
	}
}

func TestConnectPlayer(t *testing.T) {

	server := dummyServer()
	conn := dummyConn()

	profile := server.profiles["pako"]
	player := NewPlayer(profile, conn)

	server.connectPlayer(player)

	if len(server.connectedPlayers) == 0 {
		t.Fail()
	}

}

func TestDeconnectPlayer(t *testing.T) {

	server := dummyServer()
	conn := dummyConn()

	profile := server.profiles["pako"]
	player := NewPlayer(profile, conn)

	server.connectPlayer(player)

	server.deconnectPlayer(player)

	if len(server.connectedPlayers) != 0 {
		t.Fail()
	}

}
