package darwin

import (
	"fmt"
	"os"

	"github.com/Konstantin8105/c4go/noarch"
)

// BuiltinExpect handles __builtin_expect().
func BuiltinExpect(a, b int) int {
	return noarch.BoolToInt(a != b)
}

var isTest bool = false

// AssertRtn handles __assert_rtn().
func AssertRtn(
	functionName, filePath []byte,
	lineNumber int,
	expression []byte,
) bool {
	fmt.Fprintf(
		os.Stderr,
		"Assertion failed: (%s), function %s, file %s, line %d.\n",
		noarch.CStringToString(expression),
		noarch.CStringToString(functionName),
		noarch.CStringToString(filePath),
		lineNumber,
	)
	if !isTest {
		os.Exit(134)
	}

	return true
}
