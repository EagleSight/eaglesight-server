package main

import (
	"flag"

	"github.com/eaglesight/eaglesight-backend/game"
)

func main() {

	masterURL := flag.String("master_url", "", "Master's URL. Leave empty to run in 'local' mode")
	secret := flag.String("secret", "", "Secret to include in the request")

	master, _ := game.NewMaster(*masterURL, *secret)

	params := game.LoadGameParametersFromFile()

	server := game.NewServer(params)

	master.Ready()

	server.Run()

	// Tell the master to stop this server
	master.Stop()

}
