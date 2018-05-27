package main

import (
	"flag"

	"github.com/eaglesight/eaglesight-server/game"
	"github.com/eaglesight/eaglesight-server/world"
	"github.com/eaglesight/eaglesight-server/wsconnector"
)

func main() {

	terrainLocation := flag.String("map", "./map.esmap", ".esmap file location")
	paramsFileLocation := flag.String("conf", "./game.json", "config of the game")
	wsport := flag.Uint("wsport", 8000, "Port on which the websocket's server will listen")

	flag.Parse()

	params := game.LoadGameParametersFromFile(*paramsFileLocation)
	terrain, _ := world.LoadTerrain(*terrainLocation)
	world := world.NewWorld(terrain)

	server := game.NewServer(params)
	wsconn := wsconnector.NewConnector(uint16(*wsport))
	server.Run(world, wsconn)
}
