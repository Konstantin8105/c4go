package ast

import (
	"testing"
)

func TestStaticAssertDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x3526a70 <line:3108:1, col:80> col:1`: &StaticAssertDecl{
			Addr:       0x3526a70,
			Pos:        NewPositionFromString("line:3108:1, col:80"),
			Position2:  "col:1",
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}
