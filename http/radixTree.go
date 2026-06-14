package http

import "log"

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

	Child map[byte]*RadixNode
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
				panic("Only one wild card is allowed per path segment. Multiple in the path: " + s)
			}
		}
		return s[start:], start

	}

	return "", -1
}

func (r *RadixNode) addRoute(path string, data any) {
	fullPath := path
	curr := r

	for {
		wildcard, index := findWildcard(path)

		if index < 0 {
			curr.insert(path, data)
			curr.Type = Static
			curr.Child = make(map[byte]*RadixNode)
			break
		}

		if wildcard[0] == ':' {
			// It already has a param ':' in the route
			if curr.Child[':'] != nil {
				log.Panic("Can't have two param in the same place" + curr.Child[':'].Prefix + ", " + wildcard + "at path: " + fullPath)
			}

			curr = curr.insert(wildcard[1:], nil)
			curr.Type = Param

			path = path[len(wildcard):]

			continue
		}

	}
	curr.FullPath = fullPath
	curr.Data = data
}

func (r *RadixNode) insert(p string, data any) *RadixNode {
	curr := r

	if curr.Prefix == "" && curr.Child == nil {
		curr.Prefix = p
		curr.Data = data
		return curr
	}

	for {
		splitPoint := longestCommonPrefix(curr.Prefix, p)

		if splitPoint < len(curr.Prefix) {
			// TODO: SplitNode
			splitNode := &RadixNode{
				Prefix:   curr.Prefix[splitPoint:],
				Data:     curr.Data,
				FullPath: curr.FullPath,
				Child:    curr.Child,
				Type:     curr.Type,
			}

			curr.Child = make(map[byte]*RadixNode)
			curr.Child[curr.Prefix[splitPoint]] = splitNode
			curr.Prefix = curr.Prefix[:splitPoint]

			newNode := &RadixNode{
				Prefix: p[splitPoint:],
				Data:   data,
			}
			curr.Child[p[splitPoint]] = newNode
			return newNode
		}

		newNode := curr.Child[p[splitPoint]]

		if newNode == nil {
			newNode = &RadixNode{
				Prefix: p[splitPoint:],
				Data:   data,
			}

			if curr.Child == nil {
				curr.Child = make(map[byte]*RadixNode)
			}

			curr.Child[p[splitPoint]] = newNode
			return newNode
		}

		p = p[splitPoint+1:]
		curr = newNode
	}
}

func longestCommonPrefix(s, p string) int {
	min := len(p)
	if len(s) < min {
		min = len(s)
	}

	i := 0
	for ; i < min && s[i] == p[i]; i++ {
		if s[i] != p[i] {
			break
		}
	}
	return i
}

func (r *RadixNode) lookup(p string) (data any) {
	panic("Unimplemented")
}
