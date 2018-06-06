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
	}

	runNodeTests(t, nodes)
}
