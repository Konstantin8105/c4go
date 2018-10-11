package ast

import (
	"testing"
)

func TestCXXConstructExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x1f9ac68 <col:9> 'class person' 'void (void) throw()'`: &CXXConstructExpr{
			Addr:       0x1f9ac68,
			Pos:        NewPositionFromString("col:9"),
			Type:       "class person",
			Type2:      "void (void) throw()",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
