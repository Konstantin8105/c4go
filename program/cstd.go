package program

type stdFunction struct {
	cFunc          string
	functionBody   string
	dependPackages []string
	dependFuncStd  []string
}

var std = []stdFunction{
	{
		cFunc: "double fmax(double , double )",
		functionBody: `
// fmax returns the larger of its arguments: either x or y.
func fmax(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}
`,
		dependPackages: []string{},
		dependFuncStd:  []string{},
	},
	{
		cFunc: "long double fmaxl(long double , long double )",
		functionBody: `
// fmaxl returns the larger of its arguments: either x or y.
func fmaxl(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}
`,
		dependPackages: []string{},
		dependFuncStd:  []string{},
	},
	{
		cFunc: "float cbrtf(float)",
		functionBody: `
// cbrt compute cube root
func cbrtf(x float32) float32 {
	return float32(math.Cbrt(float64(x)))
}
`,
		dependPackages: []string{"math"},
		dependFuncStd:  []string{},
	},
}
