package program

import "strings"

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
					s.dependFuncStd = append(s.dependPackages, f)
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
