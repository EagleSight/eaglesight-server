package main

import (
	"github.com/eaglesight/eaglesight-backend/game"
	"github.com/eaglesight/eaglesight-backend/world"
	"github.com/eaglesight/eaglesight-backend/wsconnector"
)

func main() {
	params := game.LoadGameParametersFromFile()
	terrain, _ := world.LoadTerrain("./map.esmap")
	world := world.NewWorld(terrain)

	server := game.NewServer(params)
	wsconn := &wsconnector.Connector{}
	server.Run(world, wsconn)
}
