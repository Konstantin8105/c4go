package ast

import (
	"testing"
)

func TestAccessSpecDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x2eff360 <line:4:1, col:7> col:1 public`: &AccessSpecDecl{
			Addr:       0x2eff360,
			Pos:        NewPositionFromString("line:4:1, col:7"),
			Position2:  "col:1",
			Name:       "public",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
