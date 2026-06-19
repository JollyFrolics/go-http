package http

type TreeData struct {
	Data   any
	Params map[string]string
}

type HTTPContext struct {
	Handler RequestHandler

	Params map[string]string

	Method string
}

type HTTPHandler func(ctx *HTTPContext)

