package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 2048,
}

func webSocketHandler(w http.ResponseWriter, r *http.Request, arena *Arena) {

	// ???
	r.Header.Del("Origin")

	// Retrive the "uuid" params from the URL
	uuid := r.FormValue("uuid")

	if uuid == "" {
		log.Println("No uuid provided")
		return
	}

	// Verify if the player is registered
	player, uid, err := arena.ValidatePlayer(uuid)

	// Something happened while retriving the player's infos
	if err != nil {
		log.Println(err)
		return
	}

	// Upgrade the websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	plane := NewPlane(uid, arena.terrain, player.Model)

	playerConn := NewPlayerConn(player, plane, conn)

	playerConn.listen(arena)

	playerConn.connect(arena.connect)

}

func main() {

	masterURL := flag.String("master_url", "", "Master's URL. Leave empty to run in 'local' mode")
	secret := flag.String("secret", "", "Secret to include in the request")

	master, err := NewMaster(*masterURL, *secret)

	//  Init the game params with the default (dev) settings
	gameParams := LoadGameParametersFromFile()

	// Only there is a master
	if master.IsReachable() {

	}

	// Get the terrain
	terrain, err := LoadTerrain()

	if err != nil {
		log.Fatal(err) // TODO: Handle this
	}

	arena := NewArena(gameParams, terrain)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		webSocketHandler(w, r, arena)
	})

	go arena.Run() // Start the arena

	go func() { //  We make that async

		time.Sleep(time.Second) // Wait a second, just to be sure that the server is started

		// Tell the master we are ready
		log.Println("Ready!")
		if err := master.Ready(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Running...")
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err)
	}

}
