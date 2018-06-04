package ast

import (
	"testing"
)

func TestCXXRecordDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x2c6c2d0 <col:1, col:7> col:7 implicit class person`: &CXXRecordDecl{
			Addr:       0x2c6c2d0,
			Pos:        NewPositionFromString("col:1, col:7"),
			Prev:       "",
			Position2:  "col:7",
			Kind:       "class",
			Implicit:   true,
			Name:       "person",
			Definition: false,
			ChildNodes: []Node{},
		},
		`0x2c6c2d0 <col:1, col:7> class person`: &CXXRecordDecl{
			Addr: 0x2c6c2d0,
			Pos:  NewPositionFromString("col:1, col:7"),
			// Prev: "",
			// Position2:  "col:7",
			Kind:       "class",
			Implicit:   false,
			Name:       "person",
			Definition: false,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
