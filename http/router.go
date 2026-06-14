package http

import (
)

type Router struct {
	method string

	prefixTree RadixNode
}
