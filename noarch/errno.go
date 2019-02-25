package noarch

var errno int

func ErrnoLocation() []int {
	return []int{errno}
}
