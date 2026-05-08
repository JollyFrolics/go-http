package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
)

type ServerState struct {
	address          string
	workingDirectory fs.FS
}

func ParseFlag() (state ServerState, err error) {
	portFlag := flag.Int("port", 8080, "The port to bind to");
	addressFlag := flag.String("address", "", "The address of bind to");
	workingDirectory := flag.String("dir", ".", "The directory to expose");

	flag.Parse();

	if !fs.ValidPath(*workingDirectory) {
		log.Fatalln(workingDirectory, " doesn't exist");
	}

	state = ServerState{
		address:          fmt.Sprintf("%v:%v", *portFlag, *addressFlag),
		workingDirectory: os.DirFS(*workingDirectory),
	};

	return state, nil;
}

func main() {

}
