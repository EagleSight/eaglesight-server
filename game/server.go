package game

import (
	"encoding/binary"
	"errors"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/eaglesight/eaglesight-backend/world"
	"github.com/gorilla/websocket"
)

// Server ...
type Server struct {
	gameID           string
	connect          chan *Player
	deconect         chan *Player
	broadcast        chan []byte
	connectedPlayers map[uint8]*Player
	profilesPool     map[string]PlayerProfile
	world            *world.World
	mux              sync.Mutex
}

// NewServer return a arena with default settings (TEST THIS!)
func NewServer(params Parameters) *Server {

	// Put all the registered players in a map
	pool := make(map[string]PlayerProfile)

	// Fill up the map
	for _, profile := range params.Players {
		pool[profile.UUID] = profile
	}

	return &Server{
		gameID:           params.GameID,
		connect:          make(chan *Player, 1),
		deconect:         make(chan *Player, 1),
		connectedPlayers: make(map[uint8]*Player),
		profilesPool:     pool,
		world:            world.NewWorld(),
	}
}

// GetProfileByUUID ...
func (s *Server) getProfileByUUID(uuid string) (profile PlayerProfile, err error) {

	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.profilesPool[uuid]; ok {
		profile = s.profilesPool[uuid]
		delete(s.profilesPool, uuid)
		return profile, nil
	}

	return profile, errors.New("UUID not found in profilesPool")
}

// TEST THIS! (How ?)
func generateSnapshot(state world.State) []byte {

	offset := 1 + 1 // uint8 + uint8
	const playerDataLenght = 1 + 1 + 1 + 3*4 + 3*4
	snapshot := make([]byte, offset+len(state.Planes)*playerDataLenght)

	snapshot[0] = uint8(0x3)
	snapshot[1] = uint8(len(state.Planes))

	for _, plane := range state.Planes {

		// UID
		snapshot[offset] = plane.UID

		// TODO: Dammage
		snapshot[offset+1] = 0

		// TODO: Firing
		snapshot[offset+2] = 0

		// Location
		binary.BigEndian.PutUint32(snapshot[offset+3:], math.Float32bits(float32(plane.Location.X)))
		binary.BigEndian.PutUint32(snapshot[offset+7:], math.Float32bits(float32(plane.Location.Y)))
		binary.BigEndian.PutUint32(snapshot[offset+11:], math.Float32bits(float32(plane.Location.Z)))

		// Rotation
		binary.BigEndian.PutUint32(snapshot[offset+15:], math.Float32bits(float32(plane.Rotation.X)))
		binary.BigEndian.PutUint32(snapshot[offset+19:], math.Float32bits(float32(plane.Rotation.Y)))
		binary.BigEndian.PutUint32(snapshot[offset+23:], math.Float32bits(float32(plane.Rotation.Z)))

		offset += playerDataLenght
	}

	return snapshot

}

func (s *Server) sendToAll(message []byte) {
	s.mux.Lock()
	for _, p := range s.connectedPlayers {
		p.send <- message
	}
	s.mux.Unlock()
}

// Run start the Arena
func (s *Server) Run() {

	// Set up the websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webSocketHandler(w, r, s)
	})

	previousTickTime := time.Now()

	c := time.Tick(time.Second / 20)

	go func() {
		for {
			select {
			case now := <-c:
				// Calculate the time since the last time updated
				deltaT := now.Sub(previousTickTime).Seconds()
				// Save for the next time
				previousTickTime = now
				// Send inputs to all the players
				s.sendToAll(generateSnapshot(s.world.Tick(deltaT)))
			case player := <-s.connect:
				s.connectPlayer(player)
			case player := <-s.deconect:
				s.deconnectPlayer(player)
			}
		}
	}()

	log.Println("Running...")
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err)
	}

}

// sendPlayersList Sends the list of all the connected players, including "player" itself in first position
func (s *Server) sendPlayersList(player *Player) {

	offset := 1 + 1

	s.mux.Lock()
	message := make([]byte, offset+len(s.connectedPlayers))
	log.Println("len(s.connectedPlayers) ", len(s.connectedPlayers))

	message[0] = 0x4 // 4 == List of players
	message[1] = player.profile.UID
	for k := range s.connectedPlayers {
		message[offset] = k
		offset++
	}
	s.mux.Unlock()

	player.send <- message

}

func (s *Server) connectPlayer(player *Player) {

	log.Printf("Connecting player with UUID %s", player.profile.UUID)

	s.sendPlayersList(player)

	// TODO: Manage error
	player.input, _ = s.world.AddPlane(player.profile.UID, player.profile.Model)

	// TODO: Notifying everybody should probably be in another function
	// 0x1 + player's uid
	message := make([]byte, 2)

	message[0] = 0x1 // Connection
	message[1] = player.profile.UID

	s.sendToAll(message)

	s.mux.Lock()

	// Add the plane
	s.connectedPlayers[player.profile.UID] = player

	log.Println(player.profile.Name + " connected.")

	s.mux.Unlock()

}

func (s *Server) deconnectPlayer(player *Player) {

	log.Printf("Deconnection of %s in progress.\n", player.profile.Name)

	s.mux.Lock()

	// Remove the player from the players list
	if _, ok := s.connectedPlayers[player.profile.UID]; ok {

		// Put the profile back in the pool
		s.profilesPool[player.profile.UUID] = player.profile

		delete(s.connectedPlayers, player.profile.UID)
	}

	s.mux.Unlock()

	// Close the connection
	player.conn.Close()

	// Close the channel, which will put the plane in an "IsNoMore" state
	close(player.input)

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
