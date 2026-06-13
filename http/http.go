package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"

	"github.com/google/uuid"
)



type ServerState struct {
	Address          string
	WorkingDirectory string
	Debug            bool
	MaxPackageSize   int64
	MaxHeaderSize    int64
	MaxBodySize      int64
	Router           Router
}

type Router struct {
	path string

}

type RouteContext struct {
	path        string
	isDirectory bool
	params      map[string]string
	handler		RequestHandler
}

func (r *Router) Add(method string, path string, handle RequestHandler) error {
	if method == "" {
		return errors.New("Method shouldn't be empty. Method not registered.")
	}

	if path == "" {
		return errors.New("Path shouldn't be empty. Method not registered.")
	}

	if handle == nil {
		return errors.New("Handle is nil. Method not registered.")
	}

	return nil
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

func (server *ServerState) StartHTTP() {
	listener, err := net.Listen("tcp", server.Address)

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go server.HandleHTTPConnection(conn)
	}
}
