package ast

import (
	"testing"
)

func TestFormatAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7fcc8d8ecee8 <col:6> Implicit printf 2 3`: &FormatAttr{
			Addr:         0x7fcc8d8ecee8,
			Pos:          NewPositionFromString("col:6"),
			IsImplicit:   true,
			IsInherited:  false,
			FunctionName: "printf",
			Unknown1:     2,
			Unknown2:     3,
			ChildNodes:   []Node{},
		},
		`0x7fcc8d8ecff8 </usr/include/sys/cdefs.h:351:18, col:61> printf 2 3`: &FormatAttr{
			Addr:         0x7fcc8d8ecff8,
			Pos:          NewPositionFromString("/usr/include/sys/cdefs.h:351:18, col:61"),
			IsImplicit:   false,
			IsInherited:  false,
			FunctionName: "printf",
			Unknown1:     2,
			Unknown2:     3,
			ChildNodes:   []Node{},
		},
		`0x273b4d0 <line:357:12> Inherited printf 2 3`: &FormatAttr{
			Addr:         0x273b4d0,
			Pos:          NewPositionFromString("line:357:12"),
			IsImplicit:   false,
			IsInherited:  true,
			FunctionName: "printf",
			Unknown1:     2,
			Unknown2:     3,
			ChildNodes:   []Node{},
		},
	}

	runNodeTests(t, nodes)
}
