package program

import (
	"fmt"
	"strings"
)

type stdFunction struct {
	cFunc          string
	functionBody   string
	dependPackages []string
	dependFuncStd  []string
}

func init() {
	source := `

//---
// fmax returns the larger of its arguments: either x or y.
// c function : double fmax(double , double )
// dep pkg    : 
// dep func   :
func fmax(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}

//---
// fmaxl returns the larger of its arguments: either x or y.
// c function : long double fmaxl(long double , long double )
// dep pkg    : 
// dep func   :
func fmaxl(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}

//---
// cbrt compute cube root
// c function : float cbrtf(float)
// dep pkg    : math
// dep func   :
func cbrtf(x float32) float32 {
	return float32(math.Cbrt(float64(x)))
}



//---
// __signbitf ...
// c function : int __signbitf(float)
// dep pkg    : math
// dep func   : BoolToInt
func __signbitf(x float32) int32 {
	return BoolToInt(math.Signbit(float64(x)))
}



//---
// __builtin_signbitf ...
// c function : int __builtin_signbitf(float)
// dep pkg    : math
// dep func   : BoolToInt
func __builtin_signbitf(x float32) int32 {
	return BoolToInt(math.Signbit(float64(x)))
}

//---
// __inline_signbitf ...
// c function : int __inline_signbitf(float)
// dep pkg    : math
// dep func   : BoolToInt
func __inline_signbitf(x float32) int32 {
	return BoolToInt(math.Signbit(float64(x)))
}

//---
// BoolToInt converts boolean value to an int, which is a common operation in C.
// 0 and 1 represent false and true respectively.
// c function : int BoolToInt(int)
// dep pkg    : 
// dep func   :
func BoolToInt(x bool) int32 {
	if x {
		return 1
	}

	return 0
}




//---
// fma returns x*y+z.
// c function : double fma(double, double, double)
// dep pkg    : 
// dep func   :
func fma(x, y, z float64) float64 {
	return x*y + z
}


//---
// fmal returns x*y+z.
// c function : long double fmal(long double, long double, long double)
// dep pkg    : 
// dep func   :
func fmal(x, y, z float64) float64 {
	return x*y + z
}



//---
// fmaf returns x*y+z.
// c function : float fmaf(float, float, float)
// dep pkg    : 
// dep func   :
func fmaf(x, y, z float32) float32 {
	return x*y + z
}


`
	// split source by parts
	var (
		splitter = "//---"
		cFunc    = "c function :"
		depPkg   = "dep pkg    :"
		depFunc  = "dep func   :"
	)
	cs := strings.Split(source, splitter)
	for _, c := range cs {
		if strings.TrimSpace(c) == "" {
			continue
		}
		var s stdFunction
		s.functionBody = c

		lines := strings.Split(c, "\n")
		var foundCfunc, foundDepPkg, foundDepFunc bool
		for _, line := range lines {
			if index := strings.Index(line, cFunc); index > 0 {
				line = line[index+len(cFunc):]
				s.cFunc = strings.TrimSpace(line)
				if line == "" {
					panic(fmt.Errorf("no function name : %v", lines))
				}
				foundCfunc = true
			}
			if index := strings.Index(line, depPkg); index > 0 {
				line = line[index+len(depPkg):]
				pkgs := strings.Split(line, " ")
				for _, pkg := range pkgs {
					if pkg == "" || strings.TrimSpace(pkg) == "" {
						continue
					}
					s.dependPackages = append(s.dependPackages, pkg)
				}
				foundDepPkg = true
			}
			if index := strings.Index(line, depFunc); index > 0 {
				line = line[index+len(depFunc):]
				funcs := strings.Split(line, " ")
				for _, f := range funcs {
					if f == "" || strings.TrimSpace(f) == "" {
						continue
					}
					s.dependFuncStd = append(s.dependFuncStd, f)
				}
				foundDepFunc = true
			}
		}
		if !foundCfunc || !foundDepPkg || !foundDepFunc {
			panic(c)
		}

		std = append(std, s)
	}
}

var std []stdFunction
