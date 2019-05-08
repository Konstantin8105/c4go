package noarch

import (
	"math"
)

// Signbitd ...
func Signbitd(x float64) int32 {
	return BoolToInt(math.Signbit(x))
}

// Signbitl ...
func Signbitl(x float64) int32 {
	return BoolToInt(math.Signbit(x))
}

// IsNaN ...
func IsNaN(x float64) int32 {
	return BoolToInt(math.IsNaN(x))
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

// Hypotf compute the square root of the sum of the squares of x and y
func Hypotf(x, y float32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

// Log1pf compute ln(1+arg)
func Log1pf(arg float32) float32 {
	return float32(math.Log1p(float64(arg)))
}

// Copysignf copies sign of y to absolute value of x
func Copysignf(x float32, y float32) float32 {
	return float32(math.Copysign(float64(x), float64(y)))
}

// Expf : finds e^x
func Expf(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

// Erff : finds error function value of x
func Erff(x float32) float32 {
	return float32(math.Erf(float64(x)))
}

// Erfcf : finds error function value of x
func Erfcf(x float32) float32 {
	return float32(math.Erfc(float64(x)))
}

func LRound(x float32) int32 {
	return LRoundl(float64(x))
}

func LRoundl(x float64) int32 {
	return int32(math.Round(x))
}

func LLRound(x float32) int64 {
	return LLRoundl(float64(x))
}

func LLRoundl(x float64) int64 {
	return int64(math.Round(x))
}
