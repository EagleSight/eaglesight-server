package game

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/eaglesight/eaglesight-backend/world"
	"github.com/gorilla/websocket"
)

// Server ...
type Server struct {
	gameID           string
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
		connect:          make(chan *Player, 1),
		deconnect:        make(chan *Player, 1),
		connectedPlayers: make(map[uint8]*Player),
		profiles:         profiles,
		world:            world.NewWorld(terrain),
	}
}

// Run start the server
func (s *Server) Run() {

	go s.world.Run(time.Second/100, time.Second/30)

	// Set up the websocket handler
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	webSocketHandler(w, r, s)
	// })

	for {
		select {
		case snapshot := <-s.world.Snapshots:
			s.broadcast(snapshot)
		case player := <-s.connect:
			s.connectPlayer(player)
		case player := <-s.deconnect:
			s.deconnectPlayer(player)
		}
	}

	// log.Println("Running...")
	// if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
	// 	log.Fatal(err)
	// }

}

// GetProfileByUUID ...
func (s *Server) getProfileByUUID(uuid string) (profile PlayerProfile, err error) {

	if _, ok := s.profiles[uuid]; ok {
		profile = s.profiles[uuid]
		delete(s.profiles, uuid)
		return profile, nil
	}
	return profile, errors.New("UUID not found in profilesPool")
}

// broadcast broadcasts a message to all players
func (s *Server) broadcast(message []byte) {
	for _, p := range s.connectedPlayers {
		p.Write(message)
	}
}

// sendPlayersList Sends the list of all the connected players, including "player" itself in first position
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

	s.sendToAll(message)

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
		s.profilesPool[player.profile.UUID] = player.profile

		delete(s.connectedPlayers, player.profile.UID)
	}
	// Close the connection
	player.Close()
	// Send the message to all the players left
	message := make([]byte, 2)
	message[0] = 0x2 // Deconnection
	message[1] = player.profile.UID
	s.sendToAll(message)
	log.Println(player.profile.Name + " deconnected.")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 2048,
}

func webSocketHandler(w http.ResponseWriter, r *http.Request, server *Server) {
	// Remove the Origin header
	r.Header.Del("Origin")
	// Retrive the "uuid" params from the URL
	uuid := r.FormValue("uuid")

	if uuid == "" {
		log.Println("No uuid provided")
		return
	}
	// Verify if the player is registered
	profile, err := server.getProfileByUUID(uuid)

	// Something happened while retriving the profile's infos
	if err != nil {
		log.Println(err)
		// TODO: Find a way to return a 403
		return
	}
	// Upgrade the websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Connection upgraded for %s\n", uuid)
	// Connect the player
	server.connect <- NewPlayer(profile, conn, server.deconect)
}
