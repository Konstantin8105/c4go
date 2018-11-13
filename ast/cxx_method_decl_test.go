package ast

import (
	"testing"
)

func TestCXXMethodDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x38a32c0 <line:25:5, col:37> col:9 area 'int (void)'`: &CXXMethodDecl{
			Addr:       0x38a32c0,
			Parent:     "",
			Prev:       "",
			Pos:        NewPositionFromString("line:25:5, col:37"),
			Position2:  "col:9",
			IsImplicit: false,
			IsUsed:     false,
			MethodName: "area",
			Type:       "int (void)",
			IsInline:   false,
			Other:      "",
			ChildNodes: []Node{},
		},
		`0x38a32c0 <line:25:5, col:37> col:9 used area 'int (void)'`: &CXXMethodDecl{
			Addr:       0x38a32c0,
			Parent:     "",
			Prev:       "",
			Pos:        NewPositionFromString("line:25:5, col:37"),
			Position2:  "col:9",
			IsImplicit: false,
			IsUsed:     true,
			MethodName: "area",
			Type:       "int (void)",
			IsInline:   false,
			Other:      "",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
