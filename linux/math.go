package linux

import (
	"math"

	"github.com/Konstantin8105/c4go/noarch"
)

// IsNanf ...
func IsNanf(x float32) int32 {
	return noarch.BoolToInt(math.IsNaN(float64(x)))
}

// IsInff ...
func IsInff(x float32) int32 {
	return noarch.BoolToInt(math.IsInf(float64(x), 0))
}

// IsInf ...
func IsInf(x float64) int32 {
	return noarch.BoolToInt(math.IsInf(x, 0))
}

// NaN ...
func NaN(s []byte) float64 {
	return math.NaN()
}

// Inff handles __builtin_inff().
func Inff() float32 {
	return float32(math.Inf(0))
}
