package radixTree

import (
	"errors"
	"sync"
)

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
	Mu     sync.RWMutex
	Prefix string

	Type NodeType

	FullPath string

	Data any

	Child map[byte]*RadixNode
}

func (r *RadixNode) Insert(s string, nodeType NodeType, data any) (*RadixNode, error) {
	if len(s) == 0 {
		return nil, errors.New("Can't insert a path of length 1")
	}

	if s[0] != '/' {
		return nil, errors.New("Path must start with /")
	}

	curr := r

	elementsFound := 0

	for {
		curr.Mu.Lock()

		if curr.Child == nil {
			curr.Child = make(map[byte]*RadixNode)
		}

		for i := range len(curr.Prefix) {

			if curr.Prefix[i] != s[elementsFound] {
				splitNode := &RadixNode{
					Prefix:   curr.Prefix[i:],
					Child:    curr.Child,
					Data:     curr.Data,
					FullPath: curr.FullPath,
					Type: curr.Type,
				}

				newNode := &RadixNode{
					Prefix:   s[elementsFound:],
					FullPath: s,
					Data:     data,
					Child:    make(map[byte]*RadixNode),
					Type: nodeType,
				}

				curr.Child = make(map[byte]*RadixNode)
				curr.Child[splitNode.Prefix[0]] = splitNode
				curr.Child[s[elementsFound]] = newNode

				curr.Prefix = curr.Prefix[:i]
				curr.Data = nil
				curr.FullPath = ""

				curr.Mu.Unlock()

				return newNode, nil
			}
			elementsFound++
			if elementsFound >= len(s) {
				curr.Data = data
				curr.Mu.Unlock()
				return curr, nil
			}
		}

		newNode := curr.Child[s[elementsFound]]
		if newNode == nil {
			curr.Child[s[elementsFound]] = &RadixNode{
				Prefix:   s[elementsFound:],
				Data:     data,
				FullPath: s,
				Child:    make(map[byte]*RadixNode),
				Type: nodeType,
			}

			curr.Mu.Unlock()
			return curr.Child[s[elementsFound]], nil
		}
		curr.Mu.Unlock()
		curr = newNode
	}
}

func (r *RadixNode) Lookup(s string) (*ProtocolContext, error) {
	curr := r

	if curr == nil {
		return nil, errors.New("The node is empty")
	}

	result := &ProtocolContext{
		Params: make(map[string]string),
	}

	prevTokenWasSlash := false
	prefixIndex := 0
	var prev *RadixNode = nil
	for i := range len(s) {
		token := s[i]

		if prefixIndex >= len(curr.Prefix) {
			prev = curr
			curr = curr.Child[s[i]]

			if curr == nil {
				return result, errors.New("The address doesn't match")
			}
			prefixIndex = 0
			continue
		}

		if curr.Prefix[prefixIndex] == ':' {
			if !prevTokenWasSlash {
				return result, errors.New(": should be after /.")
			}
			if prev.Prefix == "*" {

			}
		} else if curr.Prefix[prefixIndex] == '*' {

			if !prevTokenWasSlash {
				return result, errors.New("* should be after /.")
			}

		} else if curr.Prefix[prefixIndex] != token {
			return result, errors.New("The address doesn't match")

		} else if token == '/' {
			prevTokenWasSlash = true
			prefixIndex++

		} else {
			prefixIndex++
		}

	}

	result.Data = curr.Data

	return result, nil
}

