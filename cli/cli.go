package cli

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"protocol/http"
)

func ParseFlag() (state server.ServerState, err error) {
	portFlag := flag.Int("port", 8080, "The port to bind to")
	addressFlag := flag.String("address", "", "The address of bind to")
	workingDirectoryFlag := flag.String("dir", ".", "The directory to expose")
	debugModeFlag := flag.Bool("debug", false, "Set debug mode")
	protocolFlag := flag.String("protocol", "http1.1", "The protocol the server will run")

	flag.Parse()

	if !fs.ValidPath(*workingDirectoryFlag) {
		log.Println(workingDirectoryFlag, " doesn't exist")
		return server.ServerState{}, errors.New("Can't find directory to make public")
	}
	state = server.ServerState{
		Address:          fmt.Sprintf("%v:%v", *addressFlag, *portFlag),
		WorkingDirectory: os.DirFS(*workingDirectoryFlag),
		Debug:            *debugModeFlag,
		MaxPackageSize: 1 << 12,
		MaxHeaderSize: 1000,
	}

	return state, nil
}
