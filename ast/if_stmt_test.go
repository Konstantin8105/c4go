package ast

import (
	"testing"
)

func TestIfStmt(t *testing.T) {
	nodes := map[string]Node{
		`0x7fc0a69091d0 <line:11:7, line:18:7>`: &IfStmt{
			Addr:       0x7fc0a69091d0,
			Pos:        NewPositionFromString("line:11:7, line:18:7"),
			ChildNodes: []Node{},
		},
		`0x25a5f78 <line:1329:8, line:1337:26> has_else`: &IfStmt{
			Addr:       0x25a5f78,
			Pos:        NewPositionFromString("line:1329:8, line:1337:26"),
			HasElse:    true,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
