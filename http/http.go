package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type RequestHandler func(*ConnectionInfo)

type ConnectionInfo struct {
	Conn    net.Conn
	Method  string
	Reader  *bufio.Reader
	URL     *url.URL
	Version string
	Headers map[string]string
	Body    []byte
	Id      uuid.UUID
	Server  *ServerState
}

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

func (connection *ConnectionInfo) HTTPParseBody(server *ServerState) {
	if server.MaxBodySize == -1 {

		return
	}

	var data []byte = make([]byte, server.MaxBodySize)

	_, err := connection.Conn.Read(data)
	if err != nil {
		log.Println(err)
	}

	connection.Body = data
}

func (connection *ConnectionInfo) HTTPParseHeader(server *ServerState) error {
	err := connection.parseHTTPRequestLine()

	if err != nil {
		return err
	}

	headerByte := 0

	for {
		line, err := connection.Reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if line == "" {
			break
		}

		headerByte += len(line)

		if int64(headerByte) >= server.MaxHeaderSize {
			return errors.New("Max header size exceeded by" + line)
		}

		header := strings.SplitN(line, ":", 2)

		if len(header) != 2 {
			return errors.New("Invalid Header" + line)
		}

		key := strings.TrimSpace(header[0])
		value := strings.TrimSpace(header[1])

		connection.Headers[key] = value
	}
	fmt.Printf("%v\n", connection.Headers)
	return nil
}

func (connection *ConnectionInfo) parseHTTPRequestLine() error {

	line, err := connection.Reader.ReadString('\r')

	if err != nil {
		return err
	}

	log.Println(line)

	var method string
	var version string
	var path string

	count, err := fmt.Sscanf(line, "%s %s %s", &method, &path, &version)

	if err != nil {
		return err
	}

	if count != 3 {
		return errors.New("Invalid Request Line")
	}

	connection.Method = method

	url, err := url.Parse(path)
	if err != nil {
		return err
	}

	connection.URL = url

	connection.Version = version

	return nil
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
