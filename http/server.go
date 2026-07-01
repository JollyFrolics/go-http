package http

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"maps"
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
	NotFoundHandler  HTTPHandler
}

func Default() (server *ServerState) {
	return &ServerState{
		NotFoundHandler: func(ctx *HTTPContext) {
			notFoundMsg := "HTTP/1.0 404 Not Found\r\nRequest URL: " + ctx.Path
			ctx.Write([]byte(notFoundMsg))
		},
		Address: ":8000",
		Router:  &Router{},
	}
}

func (server *ServerState) Start() {
	ln, err := net.Listen("tcp", server.Address)

	defer ln.Close()

	if err != nil {
		log.Println("Couldn't get port: ", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting request: ", err)
		}

		go server.Handle(conn)
	}
}

func (server *ServerState) Handle(conn net.Conn) {
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

	err := connection.ParseHeader(server)

	if err != nil {
		log.Println(err)
		return
	}

	connection.ParseBody(server)

	connectionContext := &HTTPContext{
		ConnInfo: &connection,
		server:   server,
		Params:   connection.Headers,
		Method:   connection.Method,
		Path:     connection.URL,
	}

	tree, ok := server.Router.radixTree[connection.Method]
	if !ok {
		server.NotFoundHandler(connectionContext)
		return
	}

	reqContext, err := tree.lookup(connectionContext.Path)
	if err != nil {
		server.NotFoundHandler(connectionContext)
	}

	handler := reqContext.Data.(HTTPHandler)

	maps.Copy(connectionContext.Params, reqContext.Params)

	handler(connectionContext)
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

func (s *ServerState) AddRoute(method, path string, handler RequestHandler) {
	s.Router.addRoute(method, path, handler)
}
