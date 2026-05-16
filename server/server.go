package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/url"
	"strings"
)

type ConnectionInfo struct {
	conn    net.Conn
	Method  string
	Reader  *bufio.Reader
	URL     *url.URL
	Version string
	Headers map[string]string
}

type ServerState struct {
	Address          string
	WorkingDirectory fs.FS
	Debug            bool
	MaxPackageSize   int
	HeaderSize       int
}

func HandleConnection(conn net.Conn, serverState *ServerState) {
	defer conn.Close()

	connection := ConnectionInfo{
		conn:    conn,
		Reader:  bufio.NewReader(io.LimitReader(conn, int64(serverState.MaxPackageSize))),
		Method:  "",
		URL:     &url.URL{},
		Version: "",
		Headers: make(map[string]string),
	}

	err := HTTPParseHeader(&connection, serverState)

	if err != nil {
		log.Println(err)
		return
	}

	HTTPParseBody(&connection, serverState)
	recover()
	fmt.Println("Connection Ended")
}

func HTTPParseBody(connectionInfo *ConnectionInfo, serverState *ServerState) {

}

func HTTPParseHeader(connectionInfo *ConnectionInfo, serverState *ServerState) error {
	line, err := connectionInfo.Reader.ReadString('\r')

	if err != nil {
		return err
	}

	err = parseHTTPRequestLine(line, connectionInfo)

	if err != nil {
		return err
	}

	for {
		line, err = connectionInfo.Reader.ReadString('\r')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			break
		}

		fmt.Println(line)

		if line == "" {
			break
		}

		header := strings.SplitN(line, ":", 2)

		if len(header) != 2 {
			return errors.New("Invalid Header" + line)
		}

		key := strings.TrimSpace(header[0])
		value := strings.TrimSpace(header[1])

		connectionInfo.Headers[key] = value
	}
	fmt.Printf("%v\n", connectionInfo.Headers)
	return nil
}

func parseHTTPRequestLine(line string, connectionInfo *ConnectionInfo) error {
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

	connectionInfo.Method = method
	url, err := url.Parse(path)
	if err != nil {
		return err
	}

	connectionInfo.URL = url

	connectionInfo.Version = version

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

		go HandleConnection(conn, &serverState)
	}
}
