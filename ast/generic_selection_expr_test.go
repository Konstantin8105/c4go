package ast

import (
	"testing"
)

func TestGenericSelectionExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x3085aa0 <line:65:17, line:66:27> 'char [5]' lvalue`: &GenericSelectionExpr{
			Addr:       0x3085aa0,
			Pos:        NewPositionFromString("line:65:17, line:66:27"),
			Type:       "char [5]",
			IsLvalue:   true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
