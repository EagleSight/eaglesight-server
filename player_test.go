package main

import (
	"testing"

	"github.com/gorilla/websocket"
)

func Connect(t *testing.T) {

	profil := DefaultPlayerProfile(1)
	plane := NewPlane(profil.Token, terrain)

	player := NewPlayer(profil, plane, new(websocket.Conn)) //
	ch := make(chan *Player, 2)

	player.connect(ch)

	p := <-ch

	if p != player {
		t.Fatalf("Player didn't connect")
	}
}
