// Package noarch contains low-level functions that apply to multiple platforms.
package noarch

// BoolToInt converts boolean value to an int, which is a common operation in C.
// 0 and 1 represent false and true respectively.
func BoolToInt(x bool) int {
	if x {
		return 1
	}

	return 0
}

// NotInt performs a logical not (!) on an integer and returns an integer.
func NotInt(x int) int {
	if x == 0 {
		return 1
	}

	return 0
}

// NotUint16 performs a logical not (!) on an integer and returns an integer.
func NotUint16(x uint16) int {
	if x == 0 {
		return 1
	}

	return 0
}

// NotUint32 performs a logical not (!) on an integer and returns an integer.
func NotUint32(x uint32) int {
	if x == 0 {
		return 1
	}

	return 0
}
