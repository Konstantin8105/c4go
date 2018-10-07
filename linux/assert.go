package linux

import (
	"fmt"
	"os"

	"github.com/Konstantin8105/c4go/noarch"
)

var osExit func(int) = os.Exit

// AssertFail handles __assert_fail().
func AssertFail(
	expression, filePath *byte,
	lineNumber uint,
	functionName *byte,
) bool {
	fmt.Fprintf(
		os.Stderr,
		"a.out: %s:%d: %s: Assertion `%s' failed.\n",
		noarch.CStringToString(filePath),
		lineNumber,
		noarch.CStringToString(functionName),
		noarch.CStringToString(expression),
	)
	osExit(134)

	return true
}
