// Package noarch contains low-level functions that apply to multiple platforms.
package noarch

// BoolToInt converts boolean value to an int, which is a common operation in C.
// 0 and 1 represent false and true respectively.
func BoolToInt(x bool) int32 {
	if x {
		return 1
	}

	return 0
}

// NotInt performs a logical not (!) on an integer and returns an integer.
func NotInt(x int32) int32 {
	if x == 0 {
		return 1
	}

	return 0
}

// NotInt32 performs a logical not (!) on an integer and returns an integer.
func NotInt32(x int32) int32 {
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

// NotByte performs a logical not (!) on an integer and returns an integer.
func NotByte(x byte) int {
	if x == 0 {
		return 1
	}

	return 0
}
