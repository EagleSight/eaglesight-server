package main

import (
	"flag"

	"github.com/eaglesight/eaglesight-server/game"
	"github.com/eaglesight/eaglesight-server/world"
	"github.com/eaglesight/eaglesight-server/wsconnector"
)

func main() {

	terrainLocation := flag.String("map", "./map.esmap", ".esmap file location")
	configFileLocation := flag.String("conf", "./game.json", "config of the game")

	params := game.LoadGameParametersFromFile(*configFileLocation)
	terrain, _ := world.LoadTerrain(*terrainLocation)
	world := world.NewWorld(terrain)

	server := game.NewServer(params)
	wsconn := &wsconnector.Connector{}
	server.Run(world, wsconn)
}
