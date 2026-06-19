package http

import (
	"errors"
	"log"
	"strings"
)

type NodeType uint8

type ProtocolContext struct {
	Data   any
	Params map[string]string
}

const (
	Static NodeType = iota
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

	for start, token := range []byte(s) {

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

		for end, token := range []byte(s[start+1:]) {
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
		if len(path) == 0 {
			return
		}

		wildcard, index := findWildcard(path)

		// Static node
		if index < 0 {
			curr = curr.insert(path, data)
			curr.Type = Static

			break
		}

		// Static node
		if index > 1 {
			curr = curr.insert(path[:index-1], nil)
			curr.Type = Static

			new := curr.Child['/']
			if new == nil {
				curr.Child['/'] = &RadixNode{
					Child: make(map[byte]*RadixNode),
				}
				new = curr.Child['/']
			}

			curr = new
		}

		// Param Node
		if wildcard[0] == ':' {
			// It already has a param ':' in the route
			if curr.Type != Static {
				if curr.Prefix != wildcard[1:] {
					log.Panic("Can't have two param in the same place" + curr.Child[':'].Prefix + ", " + wildcard + "at path: " + fullPath)
				}
			}

			curr.insert(wildcard[1:], nil)
			curr.Type = Param

			if index+len(wildcard) == len(path) {
				break
			}
			newNode := &RadixNode{
				Child: make(map[byte]*RadixNode),
			}

			curr.Child['/'] = newNode

			path = path[index+len(wildcard):]

			curr = newNode

			continue
		}

		// Catch-all Node
		if curr.Type != Static {
			if curr.Prefix != wildcard[1:] {
				log.Panic("Can't have two catch all in the same place" + curr.Child[':'].Prefix + ", " + wildcard + "at path: " + fullPath)
			}
		}

		curr = curr.insert(wildcard[1:], nil)
		curr.Type = Catch_all

		if index+len(wildcard) != len(path) {
			log.Panicln("Wildcard should be the last element of the path")
		}

		break
	}
	curr.FullPath = fullPath
	curr.Data = data
}

func (r *RadixNode) insert(path string, data any) *RadixNode {
	curr := r
	p := path

	if p == "" {
		log.Panicln("Error: empty string can't be inserted")
	}

	if curr.Prefix == "" {
		curr.Prefix = p
		curr.Data = data
		if curr.Child == nil {
			curr.Child = make(map[byte]*RadixNode)
		}
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
			curr.Data = nil
			curr.FullPath = ""

			if splitPoint < len(p) {
				newNode := &RadixNode{
					Prefix: p[splitPoint:],
					Data:   data,
					Child:  make(map[byte]*RadixNode),
				}

				curr.Child[p[splitPoint]] = newNode
				return newNode
			}

		}

		if splitPoint == len(p) {
			if data != nil {
				curr.Data = data
				log.Println("Data at path: " + path + " is being overwritten.")
			}
			return curr
		}

		newNode := curr.Child[p[splitPoint]]

		if newNode == nil {
			newNode = &RadixNode{
				Prefix: p[splitPoint:],
				Child:  make(map[byte]*RadixNode),
			}

			if data != nil {
				newNode.Data = data
			}

			if curr.Child == nil {
				curr.Child = make(map[byte]*RadixNode)
			}

			curr.Child[p[splitPoint]] = newNode
			return newNode
		}

		p = p[splitPoint:]
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

func (r *RadixNode) lookup(p string) (*ProtocolContext, error) {
	curr := r
	result := &ProtocolContext{
		Params: make(map[string]string),
	}
	for {
		switch curr.Type {
		case Static:
			{
				prefixLength := len(curr.Prefix)

				if p[:prefixLength] != curr.Prefix {
					return nil, errors.New("Path not found")
				}

				if prefixLength == len(p) {
					goto exit
				}

				p = p[prefixLength:]
			}
		case Param:
			{
				// I want catch-all to return with the '/'
				if p[0] == '/' {
					p = p[1:]
				}

				parameter := curr.Prefix

				valueSize := strings.IndexByte(p, '/')

				if valueSize == -1 {
					valueSize = len(p)
				}
				value := p[:valueSize]

				result.Params[parameter] = value

				if valueSize == len(p) {
					goto exit
				}

				p = p[valueSize:]
			}
		case Catch_all:
			{
				parameter := curr.Prefix

				result.Params[parameter] = p

				result.Data = curr.Data

				return result, nil
			}
		default:
			{
				// NOTE: This is because an exhaustive matching should be done of node type
				panic("This state should be unreachable")
			}
		}
		new := curr.Child[p[0]]

		if new == nil {
			new = curr.Child['/']
		}

		if new == nil {
			return nil, errors.New("Path not found")
		}

		curr = new

	}

exit:
	result.Data = curr.Data
	return result, nil
}
