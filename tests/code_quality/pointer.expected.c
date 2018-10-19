//
//	Package main - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package code_quality

// f - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/tests/code_quality/pointer.c:1
func f(s *int) *int {
	var a int = *s
	var b *int = s
	var c *int = &a
	s = f(&a)
	s = f(b)
	s = f(c)
	return s
}

type A struct {
	a  int
	ap *int
	b  float64
	bp *float64
}
type B struct {
	a  A
	ap *A
}

func init() {
}
