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
// dep pkg    : os github.com/Konstantin8105/c4go/noarch
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
