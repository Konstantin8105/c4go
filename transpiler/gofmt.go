package transpiler

import (
	"fmt"
	"os"
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

func addBreaklines(src string) string {
	// parens
	// * ()
	// * []
	// * {}
	// * ""

	// separate long lines by comma
	// lines := strings.Split(src, "\n")
	// for i := range lines {
	// 	if len(lines[i]) > 100 {
	// 		fmt.Fprintf(os.Stdout, "%s\n", lines[i])
	// 	}
	// }

	return src
}
