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
			IsImplicit: true,
			Name:       "C",
			ChildNodes: []Node{},
		},
		`0x266fad0 <line:74:1, line:94:1> line:74:8 C++`: &LinkageSpecDecl{
			Addr:       0x266fad0,
			Pos:        NewPositionFromString("line:74:1, line:94:1"),
			Position2:  "line:74:8",
			IsImplicit: false,
			Name:       "C++",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
