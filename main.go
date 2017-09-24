package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func extractUID(query string) (uid uint64, err error) {

	i := strings.Index(query, "uid=")

	if i == -1 {
		return 0, errors.New("'uid' param is not specified")
	}

	return strconv.ParseUint(query[i+4:len(query)], 10, 32)

}

func ws(arena *Arena, w http.ResponseWriter, r *http.Request) {

	r.Header.Del("Origin")

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	uid, err := extractUID(r.URL.RawQuery)

	if err != nil {
		log.Println(err)
		return
	}

	player := &Player{
		arena: arena,
		conn:  conn,
		uid:   uint32(uid),
		send:  make(chan []byte, 16),
	}

	go player.readPump()
	go player.writePump()
}

func main() {

	log.Println("running...")

	arena := newArena()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws(arena, w, r)
	})

	go arena.Run()

	http.ListenAndServe("127.0.0.1:8000", nil)

}
