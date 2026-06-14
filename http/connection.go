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

		if headerByte >= server.MaxHeaderSize {
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
