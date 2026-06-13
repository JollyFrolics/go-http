package http

import "strings"

type NodeType uint8

type ProtocolContext struct {
	Data   any
	Params map[string]string
}

const (
	Static NodeType = iota
	Root
	Param
	Catch_all
)

type RadixNode struct {
	Prefix string

	Type NodeType

	FullPath string

	Data any

	Child   []*RadixNode
	Indices string

	Priority int32
}

func findWildcard(s string) (path string, i int) {
	escape := false

	for start, token := range []byte(path) {

		if escape {
			escape = false
			if token == ':' || token == '*' {
				continue
			}
			panic("invalid escape string in path: " + s)
		}

		if token != ':' && token != '*' {
			continue
		}

		for end, token := range []byte(path[start+1:]) {
			switch token {
			case '/':
				return s[start : start+end+1], start
			case ':', '*':
				panic("Only one wild card is allowed per path segment. Multiple in path: " + s)
			}
		}
		return s[start:], start

	}

	return "", -1
}

func (r *RadixNode) addRoute(path string, data any) {
	curr := r

	for {
		wildcard, index := findWildcard(path)

		if index < 0 {
			curr.insert(path, data)
			curr.Type = Static
			return
		}

		if wildcard[0] == ':' {

		}
	}
}

func (r *RadixNode) insert(p string, data any) *RadixNode {
	panic("Unimplemented")
}

func (r *RadixNode) lookup(p string) (data any) {
	panic("Unimplemented")
}
