package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnect(t *testing.T) {

	player := Player{}
	plane := NewPlane(1, terrain, player.Model)

	playerConn := NewPlayerConn(player, plane, new(websocket.Conn)) //
	ch := make(chan *PlayerConn, 2)

	go playerConn.connect(ch)

	p := <-ch

	if p != playerConn {
		t.Fatalf("Player didn't connect")
	}
}
