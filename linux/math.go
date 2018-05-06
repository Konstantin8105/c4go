package linux

import (
	"math"

	"github.com/Konstantin8105/c4go/noarch"
)

// IsNanf ...
func IsNanf(x float32) int {
	return noarch.BoolToInt(math.IsNaN(float64(x)))
}

// IsInff ...
func IsInff(x float32) int {
	return noarch.BoolToInt(math.IsInf(float64(x), 0))
}

// IsInf ...
func IsInf(x float64) int {
	return noarch.BoolToInt(math.IsInf(x, 0))
}
