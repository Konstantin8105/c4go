package ast

import (
	"testing"
)

func TestCXXRecordDecl(t *testing.T) {
	nodes := map[string]Node{
		`0x2c6c2d0 <col:1, col:7> col:7 implicit class person`: &RecordDecl{
			Addr:         0x2c6c2d0,
			Pos:          NewPositionFromString("col:1, col:7"),
			Prev:         "",
			Position2:    "col:7",
			Kind:         "class",
			IsImplicit:   true,
			Name:         "person",
			IsDefinition: false,
			ChildNodes:   []Node{},
		},
		`0x2c6c2d0 <col:1, col:7> class person`: &RecordDecl{
			Addr: 0x2c6c2d0,
			Pos:  NewPositionFromString("col:1, col:7"),
			// Prev: "",
			// Position2:  "col:7",
			Kind:         "class",
			IsImplicit:   false,
			Name:         "person",
			IsDefinition: false,
			ChildNodes:   []Node{},
		},
		`0x23ac438 <line:9:1, line:16:1> line:9:7 referenced class Rectangle definition`: &RecordDecl{
			Addr:         0x23ac438,
			Pos:          NewPositionFromString("line:9:1, line:16:1"),
			Prev:         "",
			Position2:    "line:9:7",
			IsReferenced: true,
			Kind:         "class",
			Name:         "Rectangle",
			IsImplicit:   false,
			IsDefinition: true,
			ChildNodes:   []Node{},
		},
		`0x23ac438 <line:9:1, line:16:1> line:9:7 class Rectangle definition`: &RecordDecl{
			Addr:         0x23ac438,
			Pos:          NewPositionFromString("line:9:1, line:16:1"),
			Prev:         "",
			Position2:    "line:9:7",
			IsReferenced: false,
			Kind:         "class",
			Name:         "Rectangle",
			IsImplicit:   false,
			IsDefinition: true,
			ChildNodes:   []Node{},
		},
		`0x38f33c0 <col:1, col:7> col:7 implicit referenced class Circle`: &RecordDecl{
			Addr:         0x38f33c0,
			Pos:          NewPositionFromString("col:1, col:7"),
			Prev:         "",
			Position2:    "col:7",
			IsImplicit:   true,
			IsReferenced: true,
			Kind:         "class",
			Name:         "Circle",
			IsDefinition: false,
			ChildNodes:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
