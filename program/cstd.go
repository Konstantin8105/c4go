package program

import (
	"fmt"
	"strings"
)

type stdFunction struct {
	cFunc          string
	includeHeader  string
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


//---
// realloc is function from stdlib.h.
// c function : void * realloc(void* , size_t )
// dep pkg    : reflect
// dep func   : memcpy
func realloc(ptr interface{}, size uint32) interface{} {
	if ptr == nil {
		return make([]byte, size)
	}
	elemType := reflect.TypeOf(ptr).Elem()
	ptrNew := reflect.MakeSlice(reflect.SliceOf(elemType), int(size), int(size)).Interface()
	// copy elements
	memcpy(ptrNew, ptr, size)
	return ptrNew
}


//---
// memcpy is function from string.h.
// c function : void * memcpy( void * , const void * , size_t )
// dep pkg    : reflect
// dep func   :
func memcpy(dst, src interface{}, size uint32) interface{} {
	switch reflect.TypeOf(src).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(src)
		d := reflect.ValueOf(dst)
		if s.Len() == 0 {
			return dst
		}
		if s.Len() > 0 {
			size /= uint32(int(s.Index(0).Type().Size()))
		}
		var val reflect.Value
		for i := 0; i < int(size); i++ {
			if i < s.Len() {
				val = s.Index(i)
			}
			d.Index(i).Set(val)
		}
	}
	return dst
}

//---
// __assert_fail from assert.h
// c function : bool __assert_fail(const char*, const char*, unsigned int, const char*)
// dep pkg    : fmt os github.com/Konstantin8105/c4go/noarch
// dep func   :
func __assert_fail(
	expression, filePath []byte,
	lineNumber uint32,
	functionName []byte,
) bool {
	fmt.Fprintf(
		os.Stderr,
		"a.out: %s:%d: %s: Assertion %s%s' failed.\n",
		noarch.CStringToString(filePath),
		lineNumber,
		noarch.CStringToString(functionName),
		string(byte(96)),
		noarch.CStringToString(expression),
	)
	os.Exit(134)

	return true
}


//---
// tolower from ctype.h
// c function : int tolower(int)
// dep pkg    : unicode
// dep func   :
func tolower (_c int32) int32 {
	return int32(unicode.ToLower(rune(_c)))
}


//---
// toupper from ctype.h
// c function : int toupper(int)
// dep pkg    : unicode
// dep func   :
func toupper(_c int32) int32 {
	return int32(unicode.ToUpper(rune(_c)))
}



//---
// __isnanf from math.h
// c function : int __isnanf(float)
// dep pkg    : math
// dep func   : BoolToInt
func __isnanf(x float32) int32 {
	return BoolToInt(math.IsNaN(float64(x)))
}

//---
// __isinff from math.h
// c function : int __isinff(float)
// dep pkg    : math
// dep func   : BoolToInt
func __isinff(x float32) int32 {
	return BoolToInt(math.IsInf(float64(x), 0))
}



//---
// __isinf from math.h
// c function : int __isinf(double)
// dep pkg    : math
// dep func   : BoolToInt
func __isinf(x float64) int32 {
	return BoolToInt(math.IsInf(x, 0))
}


//---
// __isinfl from math.h
// c function : int __isinfl(long double)
// dep pkg    : math
// dep func   : BoolToInt
func __isinfl(x float64) int32 {
	return BoolToInt(math.IsInf(x, 0))
}

//---
// __signbit from math.h
// c function : int __signbit(double)
// dep pkg    : math
// dep func   : BoolToInt
func __signbit(x float64) int32 {
	return BoolToInt(math.Signbit(x))
}

//---
// __signbitl from math.h
// c function : int __signbitl(long double)
// dep pkg    : math
// dep func   : BoolToInt
func __signbitl(x float64) int32 {
	return BoolToInt(math.Signbit(x))
}


//---
// __isnanl from math.h
// c function : int __isnanl(long double)
// dep pkg    : math
// dep func   : BoolToInt
func __isnanl(x float64) int32 {
	return BoolToInt(math.IsNaN(x))
}


//---
// __isnan from math.h
// c function : int __isnan(double)
// dep pkg    : math
// dep func   : BoolToInt
func __isnan(x float64) int32 {
	return BoolToInt(math.IsNaN(x))
}


//---
// fmin from math.h
// c function : double fmin(double , double )
// dep pkg    : 
// dep func   : 
// fmin returns the smaller of its arguments: either x or y.
func fmin(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

//---
// fminl from math.h
// c function : double fminl(long double , long double )
// dep pkg    : 
// dep func   : 
// fmin returns the smaller of its arguments: either x or y.
func fminl(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

//---
// fminf from math.h
// c function : float fminf(float , float ) 
// dep pkg    : 
// dep func   : 
// fminf returns the smaller of its arguments: either x or y.
func fminf(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

//---
// fmaxf from math.h
// c function : float fmaxf(float , float ) 
// dep pkg    : 
// dep func   : 
// fmaxf returns the larger of its arguments: either x or y.
func fmaxf(x, y float32) float32 {
	if x < y {
		return y
	}
	return x
}


//---
// expm1f from math.h
// c function : float expm1f(float) 
// dep pkg    : math
// dep func   : 
// expm1f returns e raised to the power x minus one: e^x-1
func expm1f(x float32) float32 {
	return float32(math.Expm1(float64(x)))
}


//---
// exp2f from math.h
// c function : float exp2f(float) 
// dep pkg    : math
// dep func   : 
// exp2f Returns the base-2 exponential function of x, which is 2 raised
// to the power x: 2^x
func exp2f(x float32) float32 {
	return float32(math.Exp2(float64(x)))
}


//---
// __ctype_b_loc from ctype.h
// c function : const unsigned short int** __ctype_b_loc()
// dep pkg    : unicode
// dep func   : 
func __ctype_b_loc() [][]uint16 {
	var characterTable []uint16

	for i := 0; i < 255; i++ {
		var c uint16

		// Each of the bitwise expressions below were copied from the enum
		// values, like _ISupper, etc.

		if unicode.IsUpper(rune(i)) {
			c |= ((1 << (0)) << 8)
		}

		if unicode.IsLower(rune(i)) {
			c |= ((1 << (1)) << 8)
		}

		if unicode.IsLetter(rune(i)) {
			c |= ((1 << (2)) << 8)
		}

		if unicode.IsDigit(rune(i)) {
			c |= ((1 << (3)) << 8)
		}

		if unicode.IsDigit(rune(i)) ||
			(i >= 'a' && i <= 'f') ||
			(i >= 'A' && i <= 'F') {
			// IsXDigit. This is the same implementation as the Mac version.
			// There may be a better way to do this.
			c |= ((1 << (4)) << 8)
		}

		if unicode.IsSpace(rune(i)) {
			c |= ((1 << (5)) << 8)
		}

		if unicode.IsPrint(rune(i)) {
			c |= ((1 << (6)) << 8)
		}

		// The IsSpace check is required because Go treats spaces as graphic
		// characters, which C does not.
		if unicode.IsGraphic(rune(i)) && !unicode.IsSpace(rune(i)) {
			c |= ((1 << (7)) << 8)
		}

		// http://www.cplusplus.com/reference/cctype/isblank/
		// The standard "C" locale considers blank characters the tab
		// character ('\t') and the space character (' ').
		if i == int('\t') || i == int(' ') {
			c |= ((1 << (8)) >> 8)
		}

		if unicode.IsControl(rune(i)) {
			c |= ((1 << (9)) >> 8)
		}

		if unicode.IsPunct(rune(i)) {
			c |= ((1 << (10)) >> 8)
		}

		if unicode.IsLetter(rune(i)) || unicode.IsDigit(rune(i)) {
			c |= ((1 << (11)) >> 8)
		}

		// Yes, I know this is a hideously slow way to do it but I just want to
		// test if this works right now.
		characterTable = append(characterTable, c)
	}
	return [][]uint16{characterTable}
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

	// Examples:
	// __signbitf ...
	// __inline_signbitf ...
	// __builtin_signbitf ...
	to := []string{"__inline_", "__builtin_"}
	for i, size := 0, len(cs); i < size; i++ {
		if !strings.Contains(cs[i], "__") {
			continue
		}
		for j := range to {
			part := strings.Replace(cs[i], "__", to[j], -1)
			cs = append(cs, part)
		}
	}

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
