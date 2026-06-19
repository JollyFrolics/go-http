package http

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/url"

	"github.com/google/uuid"
)

type ServerState struct {
	Address          string
	WorkingDirectory fs.FS
	Debug            bool
	MaxPackageSize   int
	MaxHeaderSize    int
	MaxBodySize      int
	Router           *Router
}

func (server *ServerState) HandleHTTPConnection(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	connection := ConnectionInfo{
		Conn:    conn,
		Method:  "",
		URL:     &url.URL{},
		Version: "",
		Headers: make(map[string]string),
		Id:      uuid.New(),
		Server:  server,
	}

	if server.MaxPackageSize < 0 {
		connection.Reader = bufio.NewReader(conn)
	} else {
		connection.Reader = bufio.NewReader(io.LimitReader(conn, int64(server.MaxPackageSize)))
	}

	err := connection.HTTPParseHeader(server)

	if err != nil {
		log.Println(err)
		return
	}

	connection.HTTPParseBody(server)
	fmt.Println("Connection Ended")
}
