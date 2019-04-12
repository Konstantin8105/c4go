package ast

import (
	"testing"
)

func TestBinaryConditionalOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x2cdbf90 <col:7, col:16> 'int'`: &BinaryConditionalOperator{
			Addr:       0x2cdbf90,
			Pos:        NewPositionFromString("col:7, col:16"),
			Type:       "int",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
