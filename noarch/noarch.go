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
func Not(x interface{}) bool {
	switch v := x.(type) {
	case bool:
		return !v
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	}
	panic(fmt.Errorf("not support type %T", x))
}
