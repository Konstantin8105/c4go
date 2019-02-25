package program

import (
	"fmt"
	"strings"

	"github.com/Konstantin8105/c4go/util"
)

// DefinitionFunction contains the prototype definition for a function.
type DefinitionFunction struct {
	// The name of the function, like "printf".
	Name string

	// The C return type, like "int".
	ReturnType string

	// The C argument types, like ["bool", "int"]. There is currently no way
	// to represent a varargs.
	ArgumentTypes []string

	// If this is not empty then this function name should be used instead
	// of the Name. Many low level functions have an exact match with a Go
	// function. For example, "sin()".
	Substitution string

	// Can be overridden with the substitution to rearrange the return variables
	// and parameters. When either of these are nil the behavior is to keep the
	// single return value and parameters the same.
	ReturnParameters []int
	Parameters       []int
}

// Each of the predefined function have a syntax that allows them to be easy to
// read (and maintain). For example:
//
//     double __builtin_fabs(double) -> noarch.Fabs
//
// Declares the prototype of __builtin_fabs (a low level function implemented
// only on Mac) with a specific substitution provided. This means that it should
// replace any instance of __builtin_fabs with:
//
//     github.com/Konstantin8105/c4go/noarch.Fabs
//
// The substitution is optional.
//
// The substituted function can also move the parameters and return value
// positions. This is called a transformation. For example:
//
//     size_t fread(void*, size_t, size_t, FILE*) -> $0, $1 = noarch.Fread($2, $3, $4)
//
// Where $0 represents the C return value and $1 and above are for each of the
// parameters.
//
// Transformations can also be used to specify variable that need to be passed
// by reference by using the prefix "&" instead of "$":
//
//     size_t fread(void*, size_t, size_t, FILE*) -> $0 = noarch.Fread(&1, $2, $3, $4)
//
var builtInFunctionDefinitions = map[string][]string{
	"assert.h": {
		// linux/assert.h
		"bool __assert_fail(const char*, const char*, unsigned int, const char*) -> linux.AssertFail",
	},
	"ctype.h": {
		// linux/ctype.h
		"const unsigned short int** __ctype_b_loc() -> linux.CtypeLoc",
		"int tolower(int) -> linux.ToLower",
		"int toupper(int) -> linux.ToUpper",
	},
	"math.h": {
		// linux/math.h
		"int __signbitf(float) -> noarch.Signbitf",
		"int __signbit(double) -> noarch.Signbitd",
		"int __signbitl(long double) -> noarch.Signbitl",
		"int __builtin_signbitf(float) -> noarch.Signbitf",
		"int __builtin_signbit(double) -> noarch.Signbitd",
		"int __builtin_signbitl(long double) -> noarch.Signbitl",
		"int __isnanf(float) -> linux.IsNanf",
		"int __isnan(double) -> noarch.IsNaN",
		"int __isnanl(long double) -> noarch.IsNaN",
		"int __isinff(float) -> linux.IsInff",
		"int __isinf(double) -> linux.IsInf",
		"int __isinfl(long double) -> linux.IsInf",
		"double __builtin_nanf(const char*) -> linux.NaN",
		"float __builtin_inff() -> linux.Inff",

		// math.h
		"int __inline_signbitf(float) -> noarch.Signbitf",
		"int __inline_signbitd(double) -> noarch.Signbitd",
		"int __inline_signbitl(long double) -> noarch.Signbitl",

		// math.h
		"double acos(double) -> math.Acos",
		"double asin(double) -> math.Asin",
		"double atan(double) -> math.Atan",
		"double atan2(double, double) -> math.Atan2",
		"double ceil(double) -> math.Ceil",
		"double cos(double) -> math.Cos",
		"double cosh(double) -> math.Cosh",
		"double fabs(double) -> math.Abs",
		"double floor(double) -> math.Floor",
		"double fmod(double, double) -> math.Mod",
		"double remainder(double, double) -> math.Remainder",
		"double ldexp(double, int) -> math.Ldexp",
		"double log(double) -> math.Log",
		"double log10(double) -> math.Log10",
		"double pow(double, double) -> math.Pow",
		"double sin(double) -> math.Sin",
		"double sinh(double) -> math.Sinh",
		"double sqrt(double) -> math.Sqrt",
		"double tan(double) -> math.Tan",
		"double tanh(double) -> math.Tanh",

		"double copysign(double, double) -> math.Copysign",
		"float copysignf(float, float) -> noarch.Copysignf",
		"long double copysignl(long double, long double) -> math.Copysign",

		"double fma(double, double, double) -> noarch.Fma",
		"float fmaf(float, float, float) -> noarch.Fmaf",
		"long double fmal(long double, long double, long double) -> noarch.Fma",

		"double fmin(double , double ) -> noarch.Fmin",
		"float fminf(float , float ) -> noarch.Fminf",
		"long double fminl(long double , long double ) -> noarch.Fmin",

		"double fmax(double , double ) -> noarch.Fmax",
		"float fmaxf(float , float ) -> noarch.Fmaxf",
		"long double fmaxl(long double , long double ) -> noarch.Fmax",

		"double expm1(double) -> math.Expm1",
		"float expm1f(float) -> noarch.Expm1f",
		"long double expm1l(long double) -> math.Expm1",

		"double exp2(double) -> math.Exp2",
		"float exp2f(float) -> noarch.Exp2f",
		"long double exp2l(long double) -> math.Exp2",

		"double exp(double) -> math.Exp",
		"float expf(float) -> noarch.Expf",
		"long double expl(long double) -> math.Exp",

		"double erf(double) -> math.Erf",
		"float erff(float) -> noarch.Erff",
		"long double erfl(long double) -> math.Erf",

		"double erfc(double) -> math.Erfc",
		"float erfcf(float) -> noarch.Erfcf",
		"long double erfcl(long double) -> math.Erfc",

		"double log2(double) -> math.Log2",
		"float log2f(float) -> noarch.Log2f",
		"long double log2l(long double) -> math.Log2",

		"double log1p(double) -> math.Log1p",
		"float log1pf(float) -> noarch.Log1pf",
		"long double log1pl(long double) -> math.Log1p",

		"double fdim(double, double) -> noarch.Fdim",
		"float fdimf(float, float) -> noarch.Fdimf",
		"long double fdiml(long double, long double) -> noarch.Fdim",

		"double asinh(double) -> math.Asinh",
		"float asinhf(float) -> noarch.Asinhf",
		"long double asinhl(long double) -> math.Asinh",

		"double acosh(double) -> math.Acosh",
		"float acoshf(float) -> noarch.Acoshf",
		"long double acoshl(long double) -> math.Acosh",

		"double atanh(double) -> math.Atanh",
		"float atanhf(float) -> noarch.Atanhf",
		"long double atanhl(long double) -> math.Atanh",

		"double sinh(double) -> math.Sinh",
		"float sinhf(float) -> noarch.Sinhf",
		"long double sinhl(long double) -> math.Sinh",

		"double cosh(double) -> math.Cosh",
		"float coshf(float) -> noarch.Coshf",
		"long double coshl(long double) -> math.Cosh",

		"double tanh(double) -> math.Tanh",
		"float tanhf(float) -> noarch.Tanhf",
		"long double tanhl(long double) -> math.Tanh",

		"double cbrt(double) -> math.Cbrt",
		"float cbrtf(float) -> noarch.Cbrtf",
		"long double cbrtl(long double) -> math.Cbrt",

		"double hypot(double, double) -> math.Hypot",
		"float hypotf(float) -> noarch.Hypotf",
		"long double hypotl(long double, long double) -> math.Hypot",
	},
	"stdio.h": {

		// linux/stdio.h
		"int _IO_getc(FILE*) -> noarch.Fgetc",
		"int _IO_putc(int, FILE*) -> noarch.Fputc",

		// stdio.h
		"int printf(const char*, ...) -> noarch.Printf",
		"int scanf(const char*, ...) -> noarch.Scanf",
		"int putchar(int) -> noarch.Putchar",
		"int puts(const char *) -> noarch.Puts",
		"FILE* fopen(const char *, const char *) -> noarch.Fopen",
		"int fclose(FILE*) -> noarch.Fclose",
		"int remove(const char*) -> noarch.Remove",
		"int rename(const char*, const char*) -> noarch.Rename",
		"int fputs(const char*, FILE*) -> noarch.Fputs",
		"FILE* tmpfile() -> noarch.Tmpfile",
		"char* fgets(char*, int, FILE*) -> noarch.Fgets",
		"void rewind(FILE*) -> noarch.Rewind",
		"int feof(FILE*) -> noarch.Feof",
		"char* tmpnam(char*) -> noarch.Tmpnam",
		"int fflush(FILE*) -> noarch.Fflush",
		"int fprintf(FILE*, const char*, ...) -> noarch.Fprintf",
		"int fscanf(FILE*, const char*, ...) -> noarch.Fscanf",
		"int fgetc(FILE*) -> noarch.Fgetc",
		"int fputc(int, FILE*) -> noarch.Fputc",
		"int getc(FILE*) -> noarch.Fgetc",
		"char * gets(char*) -> noarch.Gets",
		"int getchar() -> noarch.Getchar",
		"int putc(int, FILE*) -> noarch.Fputc",
		"int fseek(FILE*, long int, int) -> noarch.Fseek",
		"long ftell(FILE*) -> noarch.Ftell",
		"int fread(void*, int, int, FILE*) -> $0 = noarch.Fread(&1, $2, $3, $4)",
		"int fwrite(char*, int, int, FILE*) -> noarch.Fwrite",
		"int fgetpos(FILE*, int*) -> noarch.Fgetpos",
		"int fsetpos(FILE*, int*) -> noarch.Fsetpos",
		"int sprintf(char*, const char *, ...) -> noarch.Sprintf",
		"int snprintf(char*, int, const char *, ...) -> noarch.Snprintf",
		"int vsprintf(char*, const char *, ...) -> noarch.Vsprintf",
		"int vsnprintf(char*, int, const char *, ...) -> noarch.Vsnprintf",
	},
	"wchar.h": {
		// wchar.h
		"wchar_t * wcscpy(wchar_t*, const wchar_t*) -> noarch.Wcscpy",
		"int wcscmp(const wchar_t*, const wchar_t*) -> noarch.Wcscmp",
		"size_t wcslen(const wchar_t*) -> noarch.Wcslen",
	},
	"string.h": {
		// string.h
		"char* strcat(char *, const char *) -> noarch.Strcat",
		"char* strncat(char *, const char *, int) -> noarch.Strncat",
		"int strcmp(const char *, const char *) -> noarch.Strcmp",
		"char * strchr(char *, int) -> noarch.Strchr",
		"char * strstr(const char *, const char *) -> noarch.Strstr",

		"char* strcpy(const char*, char*) -> noarch.Strcpy",
		// should be: "char* strncpy(const char*, char*, size_t) -> noarch.Strncpy",
		"char* strncpy(const char*, char*, int) -> noarch.Strncpy",

		// real return type is "size_t", but it is changed to "int"
		// in according to noarch.Strlen
		"int strlen(const char*) -> noarch.Strlen",

		"char* __inline_strcat_chk(char *, const char *) -> noarch.Strcat",

		"char * memset(char *, char, unsigned int) -> noarch.Memset",
		"char * memmove(char *, char *, unsigned int) -> noarch.Memmove",
		"int memcmp(const char *, const char *, unsigned int) -> noarch.Memcmp",
		"void * memcpy(void *, const void *, size_t) -> noarch.Memcpy",
		"const char * strrchr( const char *, int) -> noarch.Strrchr",
	},
	"stdlib.h": {
		// stdlib.h
		"int abs(int) -> noarch.Abs",
		"double atof(const char *) -> noarch.Atof",
		"int atoi(const char*) -> noarch.Atoi",
		"long int atol(const char*) -> noarch.Atol",
		"long long int atoll(const char*) -> noarch.Atoll",
		"div_t div(int, int) -> noarch.Div",
		"void exit(int) -> os.Exit",
		"void free(void*) -> noarch.Free",
		"char* getenv(const char *) -> noarch.Getenv",
		"long int labs(long int) -> noarch.Labs",
		"ldiv_t ldiv(long int, long int) -> noarch.Ldiv",
		"long long int llabs(long long int) -> noarch.Llabs",
		"lldiv_t lldiv(long long int, long long int) -> noarch.Lldiv",
		"int rand() -> math/rand.Int",
		// The real definition is srand(unsigned int) however the type would be
		// different. It's easier to change the definition than create a proxy
		// function in stdlib.go.
		"void srand(long long) -> math/rand.Seed",
		"double strtod(const char *, char **) -> noarch.Strtod",
		"float strtof(const char *, char **) -> noarch.Strtof",
		"long strtol(const char *, char **, int) -> noarch.Strtol",
		"long double strtold(const char *, char **) -> noarch.Strtold",
		"long long strtoll(const char *, char **, int) -> noarch.Strtoll",
		"long unsigned int strtoul(const char *, char **, int) -> noarch.Strtoul",
		"long long unsigned int strtoull(const char *, char **, int) -> noarch.Strtoull",
		"int system(const char *) -> noarch.System",
		"void free(void*) -> _",
		"int atexit(void*) -> noarch.Atexit",
	},
	"time.h": {
		// time.h
		"time_t time(time_t *) -> noarch.Time",
		"char* ctime(const time_t *) -> noarch.Ctime",
		"struct tm * localtime(const time_t *) -> noarch.LocalTime",
		"struct tm * gmtime(const time_t *) -> noarch.Gmtime",
		"time_t mktime(struct tm *) -> noarch.Mktime",
		"char * asctime(struct tm *) -> noarch.Asctime",
	},
	"termios.h": {
		// termios.h
		"int tcsetattr(int , int , const struct termios *) -> noarch.Tcsetattr",
		"int tcgetattr(int , struct termios *) -> noarch.Tcgetattr",
		"int tcsendbreak(int , int ) -> noarch.Tcsendbreak",
		"int tcdrain(int ) -> noarch.Tcdrain",
		"int tcflush(int , int ) -> noarch.Tcflush",
		"int tcflow(int , int ) -> noarch.Tcflow",
		"void cfmakeraw(struct termios *) -> noarch.Cfmakeraw",
		"speed_t cfgetispeed(const struct termios *) -> noarch.Cfgetispeed",
		"speed_t cfgetospeed(const struct termios *) -> noarch.Cfgetospeed",
		"int cfsetispeed(struct termios *, speed_t ) -> noarch.Cfsetispeed",
		"int cfsetospeed(struct termios *, speed_t ) -> noarch.Cfsetospeed",
		"int cfsetspeed(struct termios *, speed_t ) -> noarch.Cfsetspeed",
	},
	"sys/ioctl.h": {
		"int ioctl(int , int , ... ) -> noarch.Ioctl",
	},
	"sys/time.h": {
		"int gettimeofday(struct timeval *, struct timezone *) -> noarch.Gettimeofday",
	},
	"fcntl.h": {
		"int open(const char *, int , mode_t ) -> noarch.OpenM",
		"int open(const char *, int ) -> noarch.Open",
	},
	"unistd.h": {
		"int pipe(int *) -> noarch.Pipe",
		"void exit(int) -> golang.org/x/sys/unix.Exit",
		"ssize_t write(int, const void *, size_t) -> noarch.Write",
		"ssize_t read(int, void *, size_t) -> noarch.Read",
		"int close(int) -> golang.org/x/sys/unix.CloseOnExec",
		"int isatty(int) -> noarch.Isatty",
	},
}

// GetIncludeFileNameByFunctionSignature - return name of C include header
// in according to function name and type signature
func (p *Program) GetIncludeFileNameByFunctionSignature(
	functionName, cType string) (includeFileName string, err error) {

	for k, functionList := range builtInFunctionDefinitions {
		for i := range functionList {
			if !strings.Contains(functionList[i], functionName) {
				continue
			}
			// find function name
			baseFunction := strings.Split(functionList[i], " -> ")[0]

			// separate baseFunction to function name and type
			counter := 1
			var pos int
			// var err error
			for i := len(baseFunction) - 2; i >= 0; i-- {
				if baseFunction[i] == ')' {
					counter++
				}
				if baseFunction[i] == '(' {
					counter--
				}
				if counter == 0 {
					pos = i
					break
				}
			}
			leftPart := strings.TrimSpace(baseFunction[:pos])
			rightPart := strings.TrimSpace(baseFunction[pos:])
			index := strings.LastIndex(leftPart, " ")
			if index < 0 {
				err = fmt.Errorf("Cannot found space ` ` in %v", leftPart)
				return
			}
			if strings.Replace(functionName, " ", "", -1) !=
				strings.Replace(leftPart[index+1:], " ", "", -1) {
				continue
			}
			if strings.Replace(cType, " ", "", -1) !=
				strings.Replace(leftPart[:index]+rightPart, " ", "", -1) {
				continue
			}
			return k, nil
		}
	}
	return
}

// GetFunctionDefinition will return nil if the function does not exist (is not
// registered).
func (p *Program) GetFunctionDefinition(functionName string) *DefinitionFunction {
	p.loadFunctionDefinitions()

	if f, ok := p.functionDefinitions[functionName]; ok {
		return &f
	}

	return nil
}

// AddFunctionDefinition registers a function definition. If the definition
// already exists it will be replaced.
func (p *Program) AddFunctionDefinition(f DefinitionFunction) {
	p.loadFunctionDefinitions()

	p.functionDefinitions[f.Name] = f
}

// dollarArgumentsToIntSlice converts a list of dollar arguments, like "$1, &2"
// into a slice of integers; [1, -2].
//
// This function requires at least one argument in s, but only arguments upto
// $9 or &9.
func dollarArgumentsToIntSlice(s string) []int {
	r := []int{}
	multiplier := 1

	for _, c := range s {
		if c == '$' {
			multiplier = 1
		}
		if c == '&' {
			multiplier = -1
		}

		if c >= '0' && c <= '9' {
			r = append(r, multiplier*(int(c)-'0'))
		}
	}

	return r
}

func (p *Program) loadFunctionDefinitions() {
	if p.builtInFunctionDefinitionsHaveBeenLoaded {
		return
	}

	p.functionDefinitions = map[string]DefinitionFunction{}
	p.builtInFunctionDefinitionsHaveBeenLoaded = true

	for k, v := range builtInFunctionDefinitions {
		if !p.IncludeHeaderIsExists(k) {
			continue
		}

		for _, f := range v {
			match := util.GetRegex(`^(.+) ([^ ]+)\(([, a-z*A-Z_0-9.]*)\)( -> .+)?$`).
				FindStringSubmatch(f)

			// Unpack argument types.
			argumentTypes := strings.Split(match[3], ",")
			for i := range argumentTypes {
				argumentTypes[i] = strings.TrimSpace(argumentTypes[i])
			}
			if len(argumentTypes) == 1 && argumentTypes[0] == "" {
				argumentTypes = []string{}
			}

			// Defaults for transformations.
			var returnParameters, parameters []int

			// Substitution rules.
			substitution := match[4]
			if substitution != "" {
				substitution = strings.TrimLeft(substitution, " ->")

				// The substitution might also rearrange the parameters (return and
				// parameter transformation).
				subMatch := util.GetRegex(`^(.*?) = (.*)\((.*)\)$`).
					FindStringSubmatch(substitution)
				if len(subMatch) > 0 {
					returnParameters = dollarArgumentsToIntSlice(subMatch[1])
					parameters = dollarArgumentsToIntSlice(subMatch[3])
					substitution = subMatch[2]
				}
			}

			if strings.HasPrefix(substitution, "linux.") ||
				strings.HasPrefix(substitution, "noarch.") {
				substitution = "github.com/Konstantin8105/c4go/" + substitution
			}

			p.AddFunctionDefinition(DefinitionFunction{
				Name:             match[2],
				ReturnType:       match[1],
				ArgumentTypes:    argumentTypes,
				Substitution:     substitution,
				ReturnParameters: returnParameters,
				Parameters:       parameters,
			})
		}
	}
}
