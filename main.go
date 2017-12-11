package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 2048,
}

func extractUID(query string) (uint32, error) {

	i := strings.Index(query, "uid=")

	if i == -1 {
		return 0, errors.New("'uid' param is not specified")
	}

	uid, err := strconv.ParseUint(strings.SplitAfter(query, "uid=")[1], 10, 32)

	if err != nil {
		return 0, errors.New("'uid' param was not a number")
	}

	return uint32(uid), nil

}

func webSocketHandler(w http.ResponseWriter, r *http.Request, arena *Arena) {

	r.Header.Del("Origin")

	uid, err := extractUID(r.URL.RawQuery)

	if err != nil {
		log.Println(err)
		return
	}

	// gameID is empty if there is no authentication needed
	if arena.gameID != "" {
		// Verify if the player is registered
		if _, ok := arena.registeredPlayers[uint32(uid)]; !ok {
			log.Println("Unauthorized player")
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	player := NewPlayer(uid, arena, conn)

	player.connect(arena.connect)

}

func main() {

	masterURL := flag.String("master_url", "", "Master's URL. Leave empty to run in 'local' mode")
	secret := flag.String("secret", "", "Secret to include in the request")

	master, err := NewMaster(*masterURL, *secret)

	//  Init the game params with the default (dev) settings
	gameParams := DefaultGameParameters()

	// Only there is a master
	if master.IsReachable() {
		paramsSrc, err := master.FetchParameters()

		if err != nil {
			log.Fatal(err) // TODO: Handle this
		}

		gameParams.DecodeAndUpdate(paramsSrc)
	}

	// Get the terrain
	terrain, err := LoadTerrain(gameParams.TerrainURL)

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
