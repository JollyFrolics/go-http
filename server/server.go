package server

import (
	"bufio"
	"io/fs"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"strings"
)

type HTTPMethod int

const (
	NotSet HTTPMethod = iota
	GET
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	PATCH
)

var HTTPMethodName = map[HTTPMethod]string{
	NotSet:  "Not Set",
	GET:     "GET",
	HEAD:    "HEAD",
	POST:    "POST",
	PUT:     "PUT",
	DELETE:  "DELETE",
	CONNECT: "CONNECT",
	OPTIONS: "OPTIONS",
	TRACE:   "TRACE",
	PATCH:   "PATCH",
}

func (method HTTPMethod) String() string {
	return HTTPMethodName[method]
}

type ConnectionInfo struct {
	conn   *net.Conn
	Method HTTPMethod
	Reader *bufio.Reader
	Url    *url.URL
}

type ServerState struct {
	Address          string
	WorkingDirectory fs.FS
	Debug            bool
}

func HandleConnection(conn net.Conn, serverState ServerState) {
	defer conn.Close()
	connection := ConnectionInfo{
		conn:   &conn,
		Reader: bufio.NewReader(conn),
	}

	HTTPParseHeader(&connection, serverState)
	HTTPParseBody(&connection, &serverState)
	recover()
}

func HTTPParseBody(connectionInfo *ConnectionInfo, serverState *ServerState) {
	panic("unimplemented")
}

func HTTPParseHeader(connectionInfo *ConnectionInfo, serverState ServerState) error {
	textReader := textproto.NewReader(connectionInfo.Reader)
	requestLine, err := textReader.ReadLine()
	if err != nil {
		// TODO: Send invalid request
		return err
	}
	parseHTTPRequestLine(requestLine, connectionInfo)
	return nil
}

func parseHTTPRequestLine(line string, connectionInfo *ConnectionInfo) error {
	requestChunck := strings.SplitN(line, " ", 3)

	switch strings.ToUpper(requestChunck[0]) {
	case "GET":
		{
			connectionInfo.Method = GET
		}
	case "HEAD":
		{
			connectionInfo.Method = HEAD
		}
	case "POST":
		{
			connectionInfo.Method = POST
		}
	case "PUT":
		{
			connectionInfo.Method = PUT
		}
	case "DELETE":
		{
			connectionInfo.Method = DELETE
		}
	case "CONNECT":
		{
			connectionInfo.Method = CONNECT
		}
	case "OPTIONS":
		{
			connectionInfo.Method = OPTIONS
		}
	case "TRACE":
		{
			connectionInfo.Method = TRACE
		}
	case "PATCH":
		{
			connectionInfo.Method = PATCH
		}
	default:
		{
			//TODO: Allow for custom method
			connectionInfo.Method = NotSet
		}
	}

	url, err := url.ParseRequestURI(requestChunck[1])

	if err != nil {
		return err
	}

	connectionInfo.Url = url

	return nil
}

func HTTPStart(serverState ServerState) {
	listener, err := net.Listen("tcp", serverState.Address)

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

		go HandleConnection(conn, serverState)
	}
}
