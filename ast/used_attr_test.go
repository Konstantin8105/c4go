package ast

import (
	"testing"
)

func TestUsedAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x3be4e70 <col:44>`: &UsedAttr{
			Addr:       0x3be4e70,
			Pos:        NewPositionFromString("col:44"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
