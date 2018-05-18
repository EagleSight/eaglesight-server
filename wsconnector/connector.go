package wsconnector

import (
	"log"
	"net/http"

	"github.com/eaglesight/eaglesight-backend/game"
	"github.com/gorilla/websocket"
)

// Connector is a websocket connector for the "game" package
type Connector struct {
}

// Start and initialize the connector
func (c *Connector) Start(server *game.Server) error {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webSocketHandler(w, r, server)
	})

	err := http.ListenAndServe("0.0.0.0:8000", nil)
	return err
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 2048,
}

func webSocketHandler(w http.ResponseWriter, r *http.Request, server *game.Server) {
	// Remove the Origin header
	r.Header.Del("Origin")
	// Retrive the "uuid" params from the URL
	uuid := r.FormValue("uuid")

	profile, err := server.Verify(uuid)

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
	server.Connect(WsPlayerConn{conn: conn}, profile)
}
