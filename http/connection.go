package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
)

type ConnectionInfo struct {
	Conn    net.Conn
	Method  string
	Reader  *bufio.Reader
	URL     string
	Version string
	Headers map[string]string
	Body    []byte
	Id      uuid.UUID
	Server  *ServerState
	Handler *HTTPHandler
}

func (connnection *ConnectionInfo) Write(res []byte) (n int, err error) {
	return connnection.Conn.Write(res)
}

func (connection *ConnectionInfo) ParseBody(server *ServerState) {
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

func (connection *ConnectionInfo) ParseHeader(server *ServerState) error {
	err := connection.parseRequestLine()

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

func (connection *ConnectionInfo) parseRequestLine() error {

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

	connection.Method = method
	connection.URL = path
	connection.Version = version

	if count != 3 {
		return errors.New("Invalid Request Line")
	}

	nodeValue, err := connection.Server.Router.radixTree[connection.Method].lookup(connection.URL)

	if err != nil {
		// Handle Page not found
		return nil
	}

	handler, ok := nodeValue.Data.(*HTTPHandler)

	if !ok {

		return errors.New("The data inserted into the radix tree is not a HTTPHandler")
	}

	connection.Handler = handler

	return nil
}
