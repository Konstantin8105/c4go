package noarch

import "fmt"

// BoolToInt converts boolean value to an int, which is a common operation in C.
// 0 and 1 represent false and true respectively.
func BoolToInt(x bool) int32 {
	if x {
		return 1
	}

	return 0
}

// Not performs a logical not (!) on an integer and returns an integer.
func Not(x interface{}) int32 {
	switch v := x.(type) {
	case int, int8, int16, int32, int64:
		return BoolToInt(v != 0)
	case uint, uint8, uint16, uint32, uint64:
		return BoolToInt(v != 0)
	case float32, float64:
		return BoolToInt(v != 0)
	}
	panic(fmt.Errorf("not support type %T", x))
}

// NotUint16 performs a logical not (!) on an integer and returns an integer.
func NotUint16(x uint16) int32 {
	if x == 0 {
		return 1
	}

	return 0
}

// NotUint32 performs a logical not (!) on an integer and returns an integer.
func NotUint32(x uint32) int32 {
	if x == 0 {
		return 1
	}

	return 0
}

// NotByte performs a logical not (!) on an integer and returns an integer.
func NotByte(x byte) int32 {
	if x == 0 {
		return 1
	}

	return 0
}
