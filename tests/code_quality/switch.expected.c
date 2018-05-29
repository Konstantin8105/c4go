/*
	Package main - transpiled by c4go

	If you have found any issues, please raise an issue at:
	https://github.com/Konstantin8105/c4go/
*/

package code_quality

// switch_function - transpiled function from  /home/lepricon/go/src/github.com/Konstantin8105/c4go/tests/code_quality/switch.c:1
func switch_function() {
	var i int = 34
	switch i {
	case (0):
		fallthrough
	case (1):
		{
			return
		}
	case (2):
		{
			_ = (i)
			return
		}
	case 3:
		{
			var c int
			return
		}
	case 4:
	case 5:
	case 6:
		fallthrough
	case 7:
		{
			var d int
			break
		}
	}
}
func init() {
}
