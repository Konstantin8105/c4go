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

// Expm1f returns e raised to the power x minus one: e^x-1
func Expm1f(x float32) float32 {
	return float32(math.Expm1(float64(x)))
}

// Exp2f Returns the base-2 exponential function of x, which is 2 raised
// to the power x: 2^x
func Exp2f(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}

// Fdim returns the positive difference between x and y.
func Fdim(x, y float64) float64 {
	if x > y {
		return x - y
	}
	return 0
}

// Fdimf returns the positive difference between x and y.
func Fdimf(x, y float32) float32 {
	if x > y {
		return x - y
	}
	return 0
}

// Log2f returns the binary (base-2) logarithm of x.
func Log2f(x float32) float32 {
	return float32(math.Log2(float64(x)))
}

// Sinhf compute hyperbolic sine
func Sinhf(a float32) float32 {
	return float32(math.Sinh(float64(a)))
}

// Coshf compute hyperbolic cose
func Coshf(a float32) float32 {
	return float32(math.Cosh(float64(a)))
}

// Tanhf compute hyperbolic tan
func Tanhf(a float32) float32 {
	return float32(math.Tanh(float64(a)))
}

// Cbrt compute cube root
func Cbrtf(x float32) float32 {
	return float32(math.Cbrt(float64(x)))
}
