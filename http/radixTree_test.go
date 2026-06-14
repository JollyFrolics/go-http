package http

import (
	"protocol/util/colorize"
	"testing"
)

func TestRadixNode_addRoute(t *testing.T) {
	tree := &RadixNode{}
	_ = tree.insert("", 100)

	curr := tree.insert("/", 100)

	if curr == nil {
		t.Error(colorize.Error("Expected a tree node got"), curr)
	}

	if curr.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree != curr {
		t.Error(colorize.Error("Expected the map of the \"/\" to be the one returned from insert call wasn't"))
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	curr = tree.insert("/home", 200)

	if curr.Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['h'] != curr {
		t.Error(colorize.Error("Expected maping to result in home got: "), curr)
	}

	curr = tree.insert("/home", 100)

	if curr.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['/'] != curr {
		t.Error(colorize.Error("Expected maping to result in /home got: "), curr)
	}

	curr = tree.insert("/home", 200)
	curr = tree.insert("/high", 300)

	if curr.Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['o'].Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Child['i'].Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Child['i'] != curr {
		t.Error(colorize.Error("Mapping of /, h, i is incorrect"), curr.Data)
	}
}
