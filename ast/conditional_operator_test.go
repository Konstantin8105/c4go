package ast

import (
	"testing"
)

func TestConditionalOperator(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc6ae0bc678 <col:6, col:89> 'void'`: &ConditionalOperator{
			Addr:       0x7fc6ae0bc678,
			Pos:        NewPositionFromString("col:6, col:89"),
			Type:       "void",
			ChildNodes: []Node{},
		},
		`0x2283ec0 <line:20693:23, col:108> 'sqlite3_destructor_type':'void (*)(void *)'`: &ConditionalOperator{
			Addr:       0x2283ec0,
			Pos:        NewPositionFromString("line:20693:23, col:108"),
			Type:       "sqlite3_destructor_type",
			Type2:      "void (*)(void *)",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
