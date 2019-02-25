//
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package code_quality

// if_1 - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/tests/code_quality/if.c:1
func if_1() {
	var a int = 5
	var b int = 2
	var c int = 4
	if a > b {
		return
	} else if c <= a {
		a = 0
	}
	_ = (a)
	_ = (b)
	_ = (c)
	var w int = func() int {
		if 2 > 1 {
			return -1
		}
		return 5
	}()
	var r int
	r = func() int {
		if 2 > 1 {
			return -1
		}
		return 5
	}()
	r = func() int {
		if 2 > 1 {
			return -1
		}
		return 5
	}()
	r = func() int {
		if w > 1 {
			return -1
		}
		return 5
	}()
	r = func() int {
		if w > 1 {
			return -1
		}
		return 5
	}()
	r = func() int {
		if map[bool]int{false: 0, true: 1}[(w > 1)]+map[bool]int{false: 0, true: 1}[(r == 4)] != 0 {
			return -1
		}
		return 5
	}()
	if w > 0 {
		r = 3
	}
}
