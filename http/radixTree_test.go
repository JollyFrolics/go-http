package http

import (
	"protocol/util/colorize"
	"testing"
)

func assertPanics(t *testing.T, f func(), msg string) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}()

	f()
}

func TestRadixNode_insert(t *testing.T) {

	tree := &RadixNode{}
	assertPanics(
		t, func() {
			_ = tree.insert("", 100)

		}, "Inserting empty string should panic didn't",
	)

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

	if tree.Child['h'] != curr {
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

	if tree.Child['h'].Child['o'].Data != 200 {
		t.Error(colorize.Error("Expected data of 200 got"), curr.Data)
	}

	if tree.Child['h'].Child['i'].Data != 300 {
		t.Error(colorize.Error("Expected data of 300 got"), curr.Data)
	}

	if tree.Child['h'].Child['i'] != curr {
		t.Error(colorize.Error("Mapping of /, h, i is incorrect"), curr.Data)
	}

	root := &RadixNode{}

	root.insert("abcd", "first")

	root.insert("ab", "second")
}


func TestRadixNode_lookup(t *testing.T) {
}

func Test_longestCommonPrefix(t *testing.T) {
	tests := []struct {
		s    string
		p    string
		want int
	}{
		{"abcd", "abcd", 4},
		{"abcd", "abef", 2},
		{"abcd", "ab", 2},
		{"ab", "abcd", 2},
		{"abcd", "xyz", 0},
		{"", "", 0},
		{"abc", "", 0},
		{"", "abc", 0},
	}

	for _, tt := range tests {
		got := longestCommonPrefix(tt.s, tt.p)

		if got != tt.want {
			t.Errorf(
				"longestCommonPrefix(%q,%q)=%d want %d",
				tt.s,
				tt.p,
				got,
				tt.want,
			)
		}
	}
}
