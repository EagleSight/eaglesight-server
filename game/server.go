package game

import (
	"errors"
	"log"
	"time"

	"github.com/eaglesight/eaglesight-backend/world"
)

// Server ...
type Server struct {
	gameID           string
	verification     chan verificationRequest
	connect          chan *Player
	deconnect        chan *Player
	connectedPlayers map[uint8]*Player
	profiles         map[string]PlayerProfile
}

// NewServer return a arena with default settings (TEST THIS!)
func NewServer(params Parameters) *Server {
	// Put all the registered players in a map
	profiles := make(map[string]PlayerProfile)

	// Fill up the map
	for _, profile := range params.Players {
		profiles[profile.UUID] = profile
	}

	return &Server{
		gameID:           params.GameID,
		verification:     make(chan verificationRequest),
		connect:          make(chan *Player, 1),
		deconnect:        make(chan *Player, 1),
		connectedPlayers: make(map[uint8]*Player),
		profiles:         profiles,
	}
}

type verificationRequest struct {
	uuid        string
	reponseChan chan PlayerProfile
}

// Verify that a UUID is assiciated with a player.
// If it's the case, return this player's profil.
// Otherwise, an error is returned
func (s *Server) Verify(uuid string) (PlayerProfile, error) {
	resp := make(chan PlayerProfile)
	s.verification <- verificationRequest{uuid: uuid, reponseChan: resp}
	profile, hasProfile := <-resp

	if !hasProfile {
		return profile, errors.New("No profile was found with this UUID")
	}
	return profile, nil
}

func (s *Server) verify(request *verificationRequest) {
	// Check if the uuid is valid
	if profile, ok := s.profiles[request.uuid]; ok {
		// Check if it's not already connected
		if _, ok = s.connectedPlayers[profile.UID]; !ok {
			request.reponseChan <- profile
			return
		}
	}

	close(request.reponseChan)
}

// Connect add a player to the server
func (s *Server) Connect(conn PlayerConn, profile PlayerProfile) {
	s.connect <- NewPlayer(profile, conn, s.deconnect)
}

// Run start the server
func (s *Server) Run(world *world.World, connectors ...Connector) {

	if len(connectors) == 0 {
		log.Fatalln("No connectors loaded.")
	}

	go world.Run(time.Second/100, time.Second/30)

	for _, connector := range connectors {
		go connector.Start(s)
	}

	for {
		select {
		case snapshot := <-world.Snapshots:
			s.broadcastMessage(snapshot)
		case request := <-s.verification:
			s.verify(&request)
		case player := <-s.connect:
			go player.Listen(world.Input)
			world.Join(player.profile.UID, player.profile.Model)
			s.connectPlayer(player)
		case player := <-s.deconnect:
			world.Leave(player.profile.UID)
			s.deconnectPlayer(player)
		}
	}

}

// broadcast broadcasts a message to all players
func (s *Server) broadcastMessage(message []byte) {
	for _, p := range s.connectedPlayers {
		p.Write(message)
	}
}

// sendPlayersList Sends the list of all the connected players
// including "player" itself in first position
func (s *Server) playersListMessage(uid uint8) []byte {

	offset := 1 + 1
	message := make([]byte, offset+len(s.connectedPlayers))
	message[0] = 0x4 // 4 == List of players
	message[1] = uid

	for k := range s.connectedPlayers {
		message[offset] = k
		offset++
	}
	return message
}

func (s *Server) connectPlayer(player *Player) {

	log.Printf("Connecting player with UUID %s", player.profile.UUID)

	// Send the players list to a player
	player.Write(s.playersListMessage(player.profile.UID))

	s.connectedPlayers[player.profile.UID] = player

	s.broadcastMessage(connectionMessage(player.profile.UID))
	log.Println(player.profile.Name + " connected.")
}

func connectionMessage(UID uint8) []byte {
	// 0x1 + player's uid
	message := make([]byte, 2)
	message[0] = 0x1 // Connection
	message[1] = UID
	return message
}

func (s *Server) deconnectPlayer(player *Player) {

	log.Printf("Deconnection of %s in progress.\n", player.profile.Name)

	// Remove the player from the players list
	if _, ok := s.connectedPlayers[player.profile.UID]; ok {
		log.Println("ok")

		delete(s.connectedPlayers, player.profile.UID)
		// Close the connection
		player.Close()
		// Send the message to all the players left
		s.broadcastMessage(deconnectionMessage(player.profile.UID))
		log.Println(player.profile.Name + " deconnected.")
	}
}

func deconnectionMessage(UID uint8) []byte {
	message := make([]byte, 2)
	message[0] = 0x2 // Deconnection
	message[1] = UID
	return message
}
