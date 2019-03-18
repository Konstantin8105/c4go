package ast

import (
	"testing"
)

func TestAnnotateAttr(t *testing.T) {
	nodes := map[string]Node{
		`0xb8a12b58 <col:62, col:92> "introduced_in=17"`: &AnnotateAttr{
			Addr:       0xb8a12b58,
			Pos:        NewPositionFromString("col:62, col:92"),
			Text:       "introduced_in=17",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
