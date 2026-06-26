package http

import (
	"log"
	"os"
	"protocol/util/colorize"
)

type TreeData struct {
	Data   any
	Params map[string]string
}

type HTTPContext struct {
	Params map[string]string

	Method string

	ConnInfo *ConnectionInfo

	Path string

	server *ServerState
}

type HTTPHandler func(ctx *HTTPContext)

func (ctx *HTTPContext) FILE_PATH(file string) {
}
