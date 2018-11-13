package ast

import (
	"testing"
)

func TestCXXThisExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x3a896e8 <col:30> 'class Rectangle3 *' this`: &CXXThisExpr{
			Addr:       0x3a896e8,
			Pos:        NewPositionFromString("col:30"),
			Type:       "class Rectangle3 *",
			IsThis:     true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
