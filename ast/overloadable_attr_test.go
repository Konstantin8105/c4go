package ast

import (
	"testing"
)

func TestOverloadableAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa3b88bbb38 <line:4:1, line:13:1>`: &OverloadableAttr{
			Addr:       0x7fa3b88bbb38,
			Pos:        NewPositionFromString("line:4:1, line:13:1"),
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
