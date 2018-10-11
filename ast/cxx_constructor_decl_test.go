package ast

import (
	"testing"
)

func TestCXXConstructorDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x2fbf910 <line:3:7> col:7 implicit used person 'void (void) throw()' inline`: &CXXConstructorDecl{
			Addr:       0x2fbf910,
			Pos:        NewPositionFromString("line:3:7"),
			Position2:  "col:7",
			IsImplicit: true,
			IsUsed:     true,
			Type:       "person",
			Type2:      "void (void) throw()",
			IsInline:   true,
			Other:      "",
			ChildNodes: []Node{},
		},
		`0x343dff8 parent 0x3475f10 prev 0x34761e0 <line:28:1, line:31:1> line:28:12 used Rectangle 'void (void)'`: &CXXConstructorDecl{
			Addr:       0x343dff8,
			Parent:     "0x3475f10",
			Prev:       "0x34761e0",
			Pos:        NewPositionFromString("line:28:1, line:31:1"),
			Position2:  "line:28:12",
			IsImplicit: false,
			IsUsed:     true,
			Type:       "Rectangle",
			Type2:      "void (void)",
			IsInline:   false,
			Other:      "",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
