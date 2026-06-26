package http

import "strings"

type Router struct {
	radixTree map[string]*RadixNode

	conn *ConnectionInfo
}

type RequestHandler func(Context)

func (r *Router) addRoute(method, path string, handler RequestHandler) {
	method = strings.ToLower(method)

	if r.radixTree == nil {
		r.radixTree = make(map[string]*RadixNode)
	}

	tree, ok := r.radixTree[method]

	if !ok {
		tree = &RadixNode{
			Child: make(map[byte]*RadixNode),
		}
		r.radixTree[method] = tree
	}

	tree.insert(path, handler)
}
