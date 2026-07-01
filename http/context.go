package http

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
