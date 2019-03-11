package noarch

var errno int32

func ErrnoLocation() []int32 {
	return []int32{errno}
}
