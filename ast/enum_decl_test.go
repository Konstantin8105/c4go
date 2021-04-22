package ast

import (
	"testing"
)

func TestEnumDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x22a6c80 <line:180:1, line:186:1> __codecvt_result`: &EnumDecl{
			Addr:       0x22a6c80,
			Parent:     0,
			Pos:        NewPositionFromString("line:180:1, line:186:1"),
			Position2:  "",
			Name:       "__codecvt_result",
			ChildNodes: []Node{},
		},
		`0x32fb5a0 <enum.c:3:1, col:45> col:6 week`: &EnumDecl{
			Addr:       0x32fb5a0,
			Parent:     0,
			Pos:        NewPositionFromString("enum.c:3:1, col:45"),
			Position2:  " col:6",
			Name:       "week",
			ChildNodes: []Node{},
		},
		`0x3a0f830 parent 0x392faf0 <line:74:5, line:76:5> line:74:10 EnumTwo`: &EnumDecl{
			Addr:       0x3a0f830,
			Parent:     0x392faf0,
			Pos:        NewPositionFromString("line:74:5, line:76:5"),
			Position2:  " line:74:10",
			Name:       "EnumTwo",
			ChildNodes: []Node{},
		},
		`0x2b002c0 prev 0x2affd78 <line:28:1, line:31:1> line:28:6 efoo`: &EnumDecl{
			Addr:       0x2b002c0,
			Prev:       "0x2affd78",
			Pos:        NewPositionFromString("line:28:1, line:31:1"),
			Position2:  " line:28:6",
			Name:       "efoo",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
