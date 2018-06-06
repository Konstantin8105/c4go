/*
	Package main - transpiled by c4go

	If you have found any issues, please raise an issue at:
	https://github.com/Konstantin8105/c4go/
*/

package code_quality

import "fmt"
import "github.com/Konstantin8105/c4go/noarch"

type size_t uint32
type __time_t int32
type va_list int64
type __gnuc_va_list int64
type __codecvt_result int

const (
	__codecvt_ok      __codecvt_result = 0
	__codecvt_partial                  = 1
	__codecvt_error                    = 2
	__codecvt_noconv                   = 3
)

var stdin *noarch.File

var stdout *noarch.File

var stderr *noarch.File

// print - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/tests/code_quality/stdio.c:3
func print() {
	fmt.Printf("Hello")
	noarch.Printf([]byte("Hello, %d\x00"), 42)
}
func init() {
	stdin = noarch.Stdin
	stdout = noarch.Stdout
	stderr = noarch.Stderr
}
