package ast

import (
	"testing"
)

func TestCXXMemberCallExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x2067880 <col:43, col:54> 'int'`: &CXXMemberCallExpr{
			Addr:       0x2067880,
			Pos:        NewPositionFromString("col:43, col:54"),
			Type:       "int",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
