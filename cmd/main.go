package main

import (
	"gotickets/internel/config"
	"gotickets/internel/server"
)

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)

	// start the server
	server.Start(db, cfg)

}
