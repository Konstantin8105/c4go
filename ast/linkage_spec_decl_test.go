package ast

import (
	"testing"
)

func TestLinkageSpecDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x2efe7d8 <col:146> col:146 implicit C`: &LinkageSpecDecl{
			Addr:       0x2efe7d8,
			Pos:        NewPositionFromString("col:146"),
			Position2:  "col:146",
			Implicit:   true,
			Name:       "C",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
