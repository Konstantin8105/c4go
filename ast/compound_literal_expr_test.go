package ast

import (
	"testing"
)

func TestCompoundLiteralExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x5575acce81f0 <col:21, col:40> 'struct node':'struct node' lvalue`: &CompoundLiteralExpr{
			Addr:       0x5575acce81f0,
			Pos:        NewPositionFromString("col:21, col:40"),
			Type1:      "struct node",
			Type2:      "struct node",
			Lvalue:     false,
			ChildNodes: []Node{},
		},
		`0x350b398 <col:24, col:31> '__CONST_SOCKADDR_ARG':'__CONST_SOCKADDR_ARG'`: &CompoundLiteralExpr{
			Addr:       0x350b398,
			Pos:        NewPositionFromString("col:24, col:31"),
			Type1:      "__CONST_SOCKADDR_ARG",
			Type2:      "__CONST_SOCKADDR_ARG",
			Lvalue:     true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
