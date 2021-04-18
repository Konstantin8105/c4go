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
		`0x1adecf0 <line:327:10, col:15> 'int' 0`: &ConstantExpr{
			Addr:       0x1adecf0,
			Pos:        NewPositionFromString("line:327:10, col:15"),
			Type:       "int",
			Value:      "0",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
