package main

import (
	"action-camera/config"
	"action-camera/routes"
)

func main() {
	config.Init()
	r := routes.NewRouter()
	r.Run(config.HttpPort)
}
