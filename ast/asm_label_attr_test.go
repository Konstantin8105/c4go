package ast

import (
	"testing"
)

func TestAsmLabelAttr(t *testing.T) {
	nodes := map[string]Node{
		`0x7ff26d8224e8 </usr/include/sys/cdefs.h:569:36> "_fopen"`: &AsmLabelAttr{
			Addr:         0x7ff26d8224e8,
			Pos:          NewPositionFromString("/usr/include/sys/cdefs.h:569:36"),
			IsInherited:  false,
			FunctionName: "_fopen",
			ChildNodes:   []Node{},
		},
		`0x7fd55a169318 </usr/include/stdio.h:325:47> Inherited "_popen"`: &AsmLabelAttr{
			Addr:         0x7fd55a169318,
			Pos:          NewPositionFromString("/usr/include/stdio.h:325:47"),
			IsInherited:  true,
			FunctionName: "_popen",
			ChildNodes:   []Node{},
		},
		`0x1271fd0 <col:42> "__xpg_sigpause" IsLiteralLabel`: &AsmLabelAttr{
			Addr:           0x1271fd0,
			Pos:            NewPositionFromString("col:42"),
			IsInherited:    false,
			FunctionName:   "__xpg_sigpause",
			IsLiteralLabel: true,
			ChildNodes:     []Node{},
		},
	}

	runNodeTests(t, nodes)
}
