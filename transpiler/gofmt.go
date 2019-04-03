package transpiler

import (
	"strings"
)

// Examples:
// from:
// var di [][][]byte = [][][]byte{[][]byte{[]byte("cq\x00"), []byte(";\x00")}, [][]byte{[]byte("pl\x00"), []byte("+\x00")}, [][]byte{[]byte("hy\x00"), []byte("-\x00")}, [][]byte{[]byte("sl\x00"), []byte("/\x00")}}
// to:
// var di [][][]byte = [][][]byte{
//	[][]byte{[]byte("cq\x00"), []byte(";\x00")},
//	[][]byte{[]byte("pl\x00"), []byte("+\x00")},
//	[][]byte{[]byte("hy\x00"), []byte("-\x00")},
//	[][]byte{[]byte("sl\x00"), []byte("/\x00")},
// }
//
// from:
// noarch.Printf([]byte("%d not ok - %s:%d: \x00"), current_test, []byte("/home/konstantin/go/src/github.com/Konstantin8105/c4go/tests/init.c\x00"), 13)
// to:
// noarch.Printf(
//		[]byte("%d not ok - %s:%d: \x00"),
//		current_test,
//		[]byte("/home/konstantin/go/src/github.com/Konstantin8105/c4go/tests/init.c\x00"),
//		13,
// )

const maxLineSymbol int = 100

func addBreaklines(src string) string {
	// separate long lines by comma
	lines := strings.Split(src, "\n")
	for i := range lines {
		line := lines[i]
		if len(line) > maxLineSymbol {
			line = formatting(line)
		}
	}

	return strings.Join(lines, "\n")
}

func formatting(line string) string {
	// parens
	// * ()
	// * []
	// * {}
	// * ""
	// separator :
	// * ,

	// check last byte
	if line[len(line)-1] != '}' && line[len(line)-1] != ')' {
		return line
	}
	if !strings.Contains(line, "=") {
		return line
	}

	// initialization levels
	levels := make([]int, len(line))
	for i := 0; i < len(line); i++ {
		levels[i] = -1
	}

	// add levels for parens - "
	for i, isInside := 0, false; i < len(line); i++ {
		if isInside {
			levels[i] = 0
		}
		if line[i] == '"' {
			isInside = !isInside
		}
	}

	_ = levels

	return line
}
