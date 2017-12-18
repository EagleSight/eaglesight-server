package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnect(t *testing.T) {

	profil := DefaultPlayerProfile(1)
	plane := NewPlane(profil.Token, terrain)

	player := NewPlayer(profil, plane, new(websocket.Conn)) //
	ch := make(chan *Player)

	go player.connect(ch)

	p := <-ch

	if p != player {
		t.Fatalf("Player didn't connect")
	}
}
