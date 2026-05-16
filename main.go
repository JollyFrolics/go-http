package main

import (
	"log"
	"protocol/http/cli"
	"protocol/http/server"
)

func main() {
	serverState, err := cli.ParseFlag()

	if err != nil {
		log.Fatal(err.Error())
	}

	server.HTTPStart(serverState)
}
