package ast

import (
	"testing"
)

func TestEnableIfAttr(t *testing.T) {
	nodes := map[string]Node{
		`0xb8add300 <col:91, col:106> "<no message provided>"`: &EnableIfAttr{
			Addr:        0xb8add300,
			Pos:         NewPositionFromString("col:91, col:106"),
			Message1:    "<no message provided>",
			IsInherited: false,
			ChildNodes:  []Node{},
		},
	}

	runNodeTests(t, nodes)
}
