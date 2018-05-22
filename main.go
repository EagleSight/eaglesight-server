package main

import (
	"github.com/eaglesight/eaglesight-server/game"
	"github.com/eaglesight/eaglesight-server/world"
	"github.com/eaglesight/eaglesight-server/wsconnector"
)

func main() {
	params := game.LoadGameParametersFromFile()
	terrain, _ := world.LoadTerrain("./map.esmap")
	world := world.NewWorld(terrain)

	server := game.NewServer(params)
	wsconn := &wsconnector.Connector{}
	server.Run(world, wsconn)
}
