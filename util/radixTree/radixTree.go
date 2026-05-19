package radixTree

import (
	"protocol/http/util/assert"
	"sync"
)

type radixNode struct {
	Mu     sync.Mutex
	Prefix string

	Child map[byte]*radixNode

	Data any
}

func (r *radixNode) Insert(s string, data any) {

	curr := r
	elementsFound := 0

	for elementsFound < len(s) {
		nextNode := curr.Child[s[elementsFound]]

		if nextNode == nil {
			curr.Mu.Lock()

			curr.Child[s[elementsFound]] = &radixNode{
				Prefix: s[elementsFound:],
				Data:   data,
				Child:  make(map[byte]*radixNode),
			}

			curr.Mu.Unlock()

			return
		}

		curr = nextNode
		for i := 0; i < len(curr.Prefix); i++ {
			if s[elementsFound+i+1] != curr.Prefix[i] {
				curr.Mu.Lock()
				curr.Data = nil

				assert.Assert(curr.Child[curr.Prefix[i]] != nil, "There should not be a child")

				newNode := &radixNode{
					Prefix: curr.Prefix[i:],
					Child:  curr.Child,
					Data:   curr.Data,
				}
				curr.Child[curr.Prefix[i]] = newNode

				assert.Assert(curr.Child[s[elementsFound+i+1]] != nil, "There should not be a child")

				curr.Child[s[elementsFound+i+1]] = &radixNode{
					Prefix: curr.Prefix[i:],
					Data:   curr.Data,
					Child:  make(map[byte]*radixNode),
				}
				curr.Mu.Unlock()
				return
			}
		}

		if len(curr.Prefix) == len(s[elementsFound+1:]) {
			curr.Data = data
			return
		}

		elementsFound += len(curr.Prefix)
	}
}
