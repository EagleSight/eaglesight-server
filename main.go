package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ws(arena *Arena, w http.ResponseWriter, r *http.Request) {

	r.Header.Del("Origin")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	player := &Player{
		arena: arena,
		conn:  conn,
	}

	log.Println("Client logged")

	go player.readPump()
}

func main() {

	log.Println("running...")

	arena := newArena()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws(arena, w, r)
	})

	go arena.run()

	http.ListenAndServe("127.0.0.1:8000", nil)

}
