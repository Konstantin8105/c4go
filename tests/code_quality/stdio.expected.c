//
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package code_quality

import "github.com/Konstantin8105/c4go/noarch"
import "fmt"

// print - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/tests/code_quality/stdio.c:3
func print() {
	fmt.Printf("Hello")
	noarch.Printf([]byte("Hello, %d\x00"), 42)
}
