package main

import (
	"os"
	"protocol/http/cli"
	"protocol/http/server"
)

func main() {
	serverState, err := cli.ParseFlag()

	if err != nil {
		os.Exit(1)
	}

	server.HTTPStart(serverState)
}
