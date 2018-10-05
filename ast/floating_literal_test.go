package ast

import (
	"testing"

	"github.com/Konstantin8105/c4go/preprocessor"
)

func TestFloatingLiteral(t *testing.T) {
	nodes := map[string]Node{
		`0x7febe106f5e8 <col:24> 'double' 1.230000e+00`: &FloatingLiteral{
			Addr:       0x7febe106f5e8,
			Pos:        NewPositionFromString("col:24"),
			Type:       "double",
			Value:      1.23,
			ChildNodes: []Node{},
		},
		`0x21c65b8 <col:41> 'double' 2.718282e+00`: &FloatingLiteral{
			Addr:       0x21c65b8,
			Pos:        NewPositionFromString("col:41"),
			Type:       "double",
			Value:      2.718282e+00,
			ChildNodes: []Node{},
		},
	}

	runNodeTests(t, nodes)
}

func TestFloatingLiteralRepairFL(t *testing.T) {
	var fl FloatingLiteral
	var file preprocessor.FilePP
	errs := RepairFloatingLiteralsFromSource(&fl, file)
	if len(errs) == 0 {
		t.Errorf("Error is empty")
	}
}