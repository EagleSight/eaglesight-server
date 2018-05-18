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
	verification     chan verificationPromise
	connect          chan *Player
	deconnect        chan *Player
	connectedPlayers map[uint8]*Player
	profiles         map[string]PlayerProfile
	world            *world.World
}

// NewServer return a arena with default settings (TEST THIS!)
func NewServer(params Parameters) *Server {
	// Put all the registered players in a map
	profiles := make(map[string]PlayerProfile)

	// Fill up the map
	for _, profile := range params.Players {
		profiles[profile.UUID] = profile
	}
	// Load the terrain
	terrain, _ := world.LoadTerrain("./map.esmap")

	return &Server{
		gameID:           params.GameID,
		verification:     make(chan verificationPromise),
		connect:          make(chan *Player, 1),
		deconnect:        make(chan *Player, 1),
		connectedPlayers: make(map[uint8]*Player),
		profiles:         profiles,
		world:            world.NewWorld(terrain),
	}
}

type verificationPromise struct {
	uuid    string
	reponse chan PlayerProfile
}

// Verify that a UUID is assiciated with a player.
// If it's the case, return this player's profil.
// Otherwise, an error is returned
func (s *Server) Verify(uuid string) (PlayerProfile, error) {
	resp := make(chan PlayerProfile)
	s.verification <- verificationPromise{uuid: uuid, reponse: resp}
	profile, noProfile := <-resp

	if noProfile {
		return profile, errors.New("No profile was found with this UUID")
	}
	return profile, nil
}

func (s *Server) verify(promise *verificationPromise) {
	// Check if the uuid is valid
	if profile, ok := s.profiles[promise.uuid]; ok {
		// Check if it's not already connected
		if _, ok = s.connectedPlayers[profile.UID]; !ok {
			promise.reponse <- profile
			return
		}
	}

	close(promise.reponse)
}

// Connect add a player to the server
func (s *Server) Connect(conn PlayerConn, profile PlayerProfile) {
	s.connect <- NewPlayer(profile, conn, s.deconnect)
}

// Run start the server
func (s *Server) Run(connectors ...Connector) {

	log.Println("Starting world")
	go s.world.Run(time.Second/100, time.Second/30)

	log.Println("Starting connectors")
	for _, connector := range connectors {
		go connector.Start(s)
	}

	log.Println("Server ready")
	for {
		select {
		case snapshot := <-s.world.Snapshots:
			s.broadcast(snapshot)
		case promise := <-s.verification:
			s.verify(&promise)
		case player := <-s.connect:
			s.connectPlayer(player)
		case player := <-s.deconnect:
			s.deconnectPlayer(player)
		}
	}

}

// broadcast broadcasts a message to all players
func (s *Server) broadcast(message []byte) {
	for _, p := range s.connectedPlayers {
		p.Write(message)
	}
}

// sendPlayersList Sends the list of all the connected players
// including "player" itself in first position
func (s *Server) sendPlayersList(player *Player) {

	offset := 1 + 1
	message := make([]byte, offset+len(s.connectedPlayers))
	log.Println("len(s.connectedPlayers) ", len(s.connectedPlayers))
	message[0] = 0x4 // 4 == List of players
	message[1] = player.profile.UID

	for k := range s.connectedPlayers {
		message[offset] = k
		offset++
	}
	player.Write(message)
}

func (s *Server) connectPlayer(player *Player) {

	log.Printf("Connecting player with UUID %s", player.profile.UUID)

	s.sendPlayersList(player)

	s.world.Join <- struct {
		UID   uint8
		Model world.PlaneModel
	}{UID: player.profile.UID, Model: player.profile.Model}

	// TODO: Notifying everybody should probably be in another function
	// 0x1 + player's uid
	message := make([]byte, 2)

	message[0] = 0x1 // Connection
	message[1] = player.profile.UID

	s.broadcast(message)

	// TODO: Put in its own method
	// Add the plane
	s.connectedPlayers[player.profile.UID] = player

	log.Println(player.profile.Name + " connected.")

}

func (s *Server) deconnectPlayer(player *Player) {

	log.Printf("Deconnection of %s in progress.\n", player.profile.Name)

	// Remove the player from the players list
	if _, ok := s.connectedPlayers[player.profile.UID]; ok {

		// Put the profile back in the pool
		s.profiles[player.profile.UUID] = player.profile

		delete(s.connectedPlayers, player.profile.UID)
	}
	// Close the connection
	player.Close()
	// Send the message to all the players left
	message := make([]byte, 2)
	message[0] = 0x2 // Deconnection
	message[1] = player.profile.UID
	s.broadcast(message)
	log.Println(player.profile.Name + " deconnected.")
}
