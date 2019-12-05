package ast

import (
	"testing"
)

func TestGNUInlineAttr(t *testing.T) {
	nodes := map[string]Node{
		`0xb8a12b58 <col:62, col:92>`: &GNUInlineAttr{
			Addr:       0xb8a12b58,
			Pos:        NewPositionFromString("col:62, col:92"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
