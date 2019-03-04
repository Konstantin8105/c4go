package transpiler

import (
	"fmt"

	"github.com/Konstantin8105/c4go/program"
)

func generateBinding(p *program.Program) {
	// outside called functions
	ds := p.GetOutsideCalledFunctions()
	for i := range ds {
		//
		// Example:
		//
		// // #include <stdlib.h>
		// // #include <stdio.h>
		// // #include <errno.h>
		// import "C"
		//
		// func Seed(i int) {
		//   C.srandom(C.uint(i))
		// }
		//

		// input data:
		// {frexp double [double int *] true true  [] []}
		//
		// output:
		// func  frexp(arg1 float64, arg2 []int) float64 {
		//		return float64(C.frexp(C.double(arg1), unsafe.Pointer(arg2)))
		// }

		p.AddMessage(p.GenerateWarningMessage(fmt.Errorf(
			"Haven`t implementation for function : `%s`", ds[i].Name), nil))
	}
	// automatic binding of function
}
