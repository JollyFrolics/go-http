package http

import "strings"

type Router struct {
	radixTree map[string]*RadixNode

	conn *ConnectionInfo
}

type RequestHandler func(ProtocolContext)

func (r *Router) addRoute(method, path string, handler RequestHandler) {
	method = strings.ToLower(method)
	tree, ok := r.radixTree[method]

	if !ok {
		tree = &RadixNode{
			Child: make(map[byte]*RadixNode),
		}
		r.radixTree[method] = tree
	}

	tree.insert(path, handler)
}

func (r *Router) GET(path string, handler RequestHandler) {
	r.addRoute("get", path, handler)
}

func (r *Router) HEAD(path string, handler RequestHandler) {
	r.addRoute("head", path, handler)
}

func (r *Router) POST(path string, handler RequestHandler) {
	r.addRoute("POST", path, handler)
}

func (r *Router) PUT(path string, handler RequestHandler) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler RequestHandler) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) CONNECT(path string, handler RequestHandler) {
	r.addRoute("CONNECT", path, handler)
}

func (r *Router) OPTIONS(path string, handler RequestHandler) {
	r.addRoute("OPTIONS", path, handler)
}

func (r *Router) TRACE(path string, handler RequestHandler) {
	r.addRoute("TRACE", path, handler)
}

func (r *Router) PATCH(path string, handler RequestHandler) {
	r.addRoute("PATCH", path, handler)
}
