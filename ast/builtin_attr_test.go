package ast

import (
	"testing"
)

func TestBuiltinAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x22e4f60 <<invalid sloc>> Implicit 905`: &BuiltinAttr{
			Addr:        0x22e4f60,
			Pos:         NewPositionFromString("<invalid sloc>"),
			IsImplicit:  true,
			IsInherited: false,
			Name:        "905",
			ChildNodes:  []Node{},
		},
		`0x22e51c8 <<invalid sloc>> Inherited Implicit 905`: &BuiltinAttr{
			Addr:        0x22e51c8,
			Pos:         NewPositionFromString("<invalid sloc>"),
			IsImplicit:  true,
			IsInherited: true,
			Name:        "905",
			ChildNodes:  []Node{},
		},
	}

	runNodeTests(t, nodes)
}
