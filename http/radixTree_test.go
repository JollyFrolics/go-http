package http

import (
	"testing"
)

func TestRadixNode_addRoute(t *testing.T) {
	tree := RadixNode{}
	tree.addRoute('/name/idea', 1)
}

