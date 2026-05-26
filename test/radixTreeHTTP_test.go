package radixTreeTest

import (
	"protocol/util/colorize"
	"protocol/http/util/radixTree"
	"testing"
)

func TestRadixTreeInsert(t *testing.T) {
	tree := radixTree.RadixNode{}
	_, err := tree.Insert("", 100)
	if err == nil {
		t.Error(colorize.Error("Passing empty string on insert should err. It didn't."))
	} else if err.Error() != "Can't insert a path of length 1" {
		t.Error(colorize.Error("Incorrect error send on passing empty string"))
	}

	curr, err := tree.Insert("/", 100)

	if curr == nil {
		t.Error(colorize.Error("Expected a tree node got"), curr)
	}

	if curr.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['/'] != curr {
		t.Error(colorize.Error("Expected the map of the \"/\" to be the one returned from insert call wasn't"))
	}

	if tree.Child['/'].Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	curr, err = tree.Insert("/home", 200)

	if curr.Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Child['/'].Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['/'].Child['h'] != curr {
		t.Error(colorize.Error("Expected maping to result in /home didn't"), curr)
	}

	curr, err = tree.Insert("/home", 100)

	if curr.Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['/'].Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['/'].Child['h'] != curr {
		t.Error(colorize.Error("Expected maping to result in /home didn't"), curr)
	}

	curr, err = tree.Insert("/home", 200)
	curr, err = tree.Insert("/high", 300)

	if curr.Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Child['/'].Data != 100 {
		t.Error(colorize.Error("Expected data of 100 got"), curr.Data)
	}

	if tree.Child['/'].Child['h'].Child['o'].Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Child['/'].Child['h'].Child['i'].Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Child['/'].Child['h'].Child['i'] != curr {
		t.Error(colorize.Error("Mapping of /, h, i is incorrect"), curr.Data)
	}
}
