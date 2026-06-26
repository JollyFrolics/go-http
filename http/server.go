package http

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"

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
		URL:     "",
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
}

func (s *ServerState) GET(path string, handler RequestHandler) {
	s.Router.addRoute("get", path, handler)
}

func (s *ServerState) HEAD(path string, handler RequestHandler) {
	s.Router.addRoute("head", path, handler)
}

func (s *ServerState) POST(path string, handler RequestHandler) {
	s.Router.addRoute("POST", path, handler)
}

func (s *ServerState) PUT(path string, handler RequestHandler) {
	s.Router.addRoute("PUT", path, handler)
}

func (s *ServerState) DELETE(path string, handler RequestHandler) {
	s.Router.addRoute("DELETE", path, handler)
}

func (s *ServerState) CONNECT(path string, handler RequestHandler) {
	s.Router.addRoute("CONNECT", path, handler)
}

func (s *ServerState) OPTIONS(path string, handler RequestHandler) {
	s.Router.addRoute("OPTIONS", path, handler)
}

func (s *ServerState) TRACE(path string, handler RequestHandler) {
	s.Router.addRoute("TRACE", path, handler)
}

func (s *ServerState) PATCH(path string, handler RequestHandler) {
	s.Router.addRoute("PATCH", path, handler)
}

func (s *ServerState) addRoute(method, path string, handler RequestHandler) {
	s.Router.addRoute(method, path, handler)
}
