package cli

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"protocol/http/server"
)

func ParseFlag() (state server.ServerState, err error) {
	portFlag := flag.Int("port", 8080, "The port to bind to")
	addressFlag := flag.String("address", "", "The address of bind to")
	workingDirectoryFlag := flag.String("dir", ".", "The directory to expose")
	debugModeFlag := flag.Bool("debug", false, "Set debug mode")

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
		HeaderSize: 1000,
	}

	return state, nil
}
