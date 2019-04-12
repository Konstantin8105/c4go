package ast

import (
	"testing"
)

func TestOpaqueValueExpr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fa855aab838 <col:63, col:95> 'unsigned long'`: &OpaqueValueExpr{
			Addr:       0x7fa855aab838,
			Pos:        NewPositionFromString("col:63, col:95"),
			Type:       "unsigned long",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
