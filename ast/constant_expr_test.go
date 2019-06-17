package ast

import (
	"testing"
)

func TestConstantExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x28f4a70 <line:223:7> 'int'`: &ConstantExpr{
			Addr:       0x28f4a70,
			Pos:        NewPositionFromString("line:223:7"),
			Type:       "int",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
