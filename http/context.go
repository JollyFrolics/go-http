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

func (con *HTTPContext) Write(res []byte) (int, error) {
	return con.ConnInfo.Write(res)
}
