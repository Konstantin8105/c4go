package noarch

import (
	"math"
)

// Signbitf ...
func Signbitf(x float32) int {
	return BoolToInt(math.Signbit(float64(x)))
}

// Signbitd ...
func Signbitd(x float64) int {
	return BoolToInt(math.Signbit(x))
}

// Signbitl ...
func Signbitl(x float64) int {
	return BoolToInt(math.Signbit(x))
}

// IsNaN ...
func IsNaN(x float64) int {
	return BoolToInt(math.IsNaN(x))
}

// Fma returns x*y+z.
func Fma(x, y, z float64) float64 {
	return x*y + z
}

// Fmaf returns x*y+z.
func Fmaf(x, y, z float32) float32 {
	return x*y + z
}

// Fmin returns the smaller of its arguments: either x or y.
func Fmin(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

// Fminf returns the smaller of its arguments: either x or y.
func Fminf(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

// Fmax returns the larger of its arguments: either x or y.
func Fmax(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}

// Fmaxf returns the larger of its arguments: either x or y.
func Fmaxf(x, y float32) float32 {
	if x < y {
		return y
	}
	return x
}

// Expm1 returns e raised to the power x minus one: e^x-1
func Expm1(x float64) float64 {
	return math.Exp(x) - 1
}

// Expm1f returns e raised to the power x minus one: e^x-1
func Expm1f(x float32) float32 {
	return float32(math.Exp(float64(x))) - 1
}
