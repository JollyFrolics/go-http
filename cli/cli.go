package main

import (
	"flag"
	"fmt"
	"os"
	"protocol/http"
	"strings"
)

func ParseFlag() (state *http.ServerState, err error) {
	portFlag := flag.Int("port", 8080, "The port to bind to")
	addressFlag := flag.String("address", "", "The address of bind to")
	workingDirectoryFlag := flag.String("dir", ".", "The directory to expose")
	debugModeFlag := flag.Bool("debug", false, "Set debug mode")
	protocolFlag := flag.String("protocol", "http1.1", "The protocol the server will run")
	maxHeaderSizeFlag := flag.Int64("maxHeaderSize", 1<<12, "Max size of the header")
	maxPackageSizeFlag := flag.Int64("maxPackageSize", 1<<20, "Maxium size of the each connection")
	maxBodySizeFlag := flag.Int64("maxBodySize", 1<<20, "Max size of the body")

	flag.Parse()

	if strings.HasPrefix(strings.ToUpper(*protocolFlag), "HTTP") {
		state = &http.ServerState{
			Address:          fmt.Sprintf("%v:%v", *addressFlag, *portFlag),
			WorkingDirectory: *workingDirectoryFlag,
			Debug:            *debugModeFlag,
			MaxBodySize:      *maxBodySizeFlag,
			MaxPackageSize:   *maxPackageSizeFlag,
			MaxHeaderSize:    *maxHeaderSizeFlag,
			Router:           http.Router{},
		}
	}

	return state, nil
}

func main() {
	state, err := ParseFlag()
	if err != nil {
		os.Exit(-1)
	}

	state.StartHTTP()
}
