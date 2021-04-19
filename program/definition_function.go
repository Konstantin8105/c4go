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

	// Each function from some source. For example: "stdio.h"
	IncludeFile string
	// If function called, then true.
	IsCalled bool

	// If function without body, then true
	HaveBody bool

	// If function from some C standard library, then true.
	IsCstdFunction bool
	pntCstd        *stdFunction

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
	"signal.h": {
		// signal.h
		"void (*signal(int , void (*)(int)))(int) -> noarch.Signal",
		"int raise(int ) -> noarch.Raise",
	},
	"errno.h": {
		// errno.h
		"int * __errno_location(void ) -> noarch.ErrnoLocation",
	},
	"math.h": {
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
		"long double copysignl(long double, long double) -> math.Copysign",

		"double expm1(double) -> math.Expm1",
		"long double expm1l(long double) -> math.Expm1",

		"double exp2(double) -> math.Exp2",
		"long double exp2l(long double) -> math.Exp2",

		"double exp(double) -> math.Exp",
		"long double expl(long double) -> math.Exp",

		"double erf(double) -> math.Erf",
		"long double erfl(long double) -> math.Erf",

		"double erfc(double) -> math.Erfc",
		"long double erfcl(long double) -> math.Erfc",

		"double log2(double) -> math.Log2",
		"long double log2l(long double) -> math.Log2",

		"double log1p(double) -> math.Log1p",
		"long double log1pl(long double) -> math.Log1p",

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
		"long double sinhl(long double) -> math.Sinh",

		"double cosh(double) -> math.Cosh",
		"long double coshl(long double) -> math.Cosh",

		"double tanh(double) -> math.Tanh",
		"long double tanhl(long double) -> math.Tanh",

		"double cbrt(double) -> math.Cbrt",
		"long double cbrtl(long double) -> math.Cbrt",

		"double hypot(double, double) -> math.Hypot",
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
		"int vprintf(const char *, ...) -> noarch.Vprintf",
		"int vfprintf(FILE *, const char *, ...) -> noarch.Vfprintf",
		"int vsnprintf(char*, int, const char *, ...) -> noarch.Vsnprintf",
		"void perror( const char *) -> noarch.Perror",
		"ssize_t getline(char **, size_t *, FILE *) -> noarch.Getline",
		"int sscanf( const char *, const char *, ...) -> noarch.Sscanf",
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
		"const char * strrchr( const char *, int) -> noarch.Strrchr",
		"char * strdup(const char *) -> noarch.Strdup",
		"char * strerror(int ) -> noarch.Strerror",
	},
	"stdlib.h": {
		// stdlib.h
		"int abs(int) -> noarch.Abs",
		"double atof(const char *) -> noarch.Atof",
		"int atoi(const char*) -> noarch.Atoi",
		"long int atol(const char*) -> noarch.Atol",
		"long long int atoll(const char*) -> noarch.Atoll",
		"div_t div(int, int) -> noarch.Div",
		"void exit(int) -> noarch.Exit",
		"void free(void*) -> noarch.Free",
		"char* getenv(const char *) -> noarch.Getenv",
		"long int labs(long int) -> noarch.Labs",
		"ldiv_t ldiv(long int, long int) -> noarch.Ldiv",
		"long long int llabs(long long int) -> noarch.Llabs",
		"lldiv_t lldiv(long long int, long long int) -> noarch.Lldiv",
		"int rand() -> noarch.Int32",
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
		"clock_t clock(void) -> noarch.Clock",
		"double difftime(time_t , time_t ) -> noarch.Difftime",
	},
	"locale.h": {
		"struct lconv * localeconv(void) -> noarch.Localeconv",
		"char * setlocale(int , const char * ) -> noarch.Setlocale",
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
		"int open(const char *, int , ...) -> noarch.Open",
	},
	"unistd.h": {
		"int pipe(int *) -> noarch.Pipe",
		"void exit(int) -> golang.org/x/sys/unix.Exit",
		"ssize_t write(int, const void *, size_t) -> noarch.Write",
		"ssize_t read(int, void *, size_t) -> noarch.Read",
		"int close(int) -> noarch.CloseOnExec",
		"int isatty(int) -> noarch.Isatty",
		"int unlink(const char *) -> noarch.Unlink",
		"int ftruncate(int , off_t ) -> noarch.Ftruncate",
	},
	"sys/stat.h": {
		"int fstat(int , struct stat  *) -> noarch.Fstat",
		"int stat(const char * , struct stat * ) -> noarch.Stat",
		"int lstat(const char * , struct stat * ) -> noarch.Lstat",
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
				err = fmt.Errorf("cannot found space ` ` in %v", leftPart)
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
			index := strings.Index(f, "->")
			_, a, w, e, err := util.ParseFunction(f[:index])
			if err != nil {
				panic(err)
			}

			// Defaults for transformations.
			var returnParameters, parameters []int

			// Substitution rules.
			substitution := strings.TrimSpace(f[index+2:])
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

			if strings.HasPrefix(substitution, "noarch.") {
				substitution = "github.com/Konstantin8105/c4go/" + substitution
			}

			p.AddFunctionDefinition(DefinitionFunction{
				Name:             a,
				ReturnType:       e[0],
				ArgumentTypes:    w,
				Substitution:     substitution,
				ReturnParameters: returnParameters,
				Parameters:       parameters,
				IsCstdFunction:   false,
			})
		}
	}

	// initialization CSTD
	for i := range std {
		_, a, w, e, err := util.ParseFunction(std[i].cFunc)
		if err != nil {
			panic(err)
		}

		p.AddFunctionDefinition(DefinitionFunction{
			Name:           a,
			ReturnType:     e[0],
			ArgumentTypes:  w,
			IsCstdFunction: true,
			pntCstd:        &std[i],
		})
	}
}

func (p *Program) SetCalled(name string) {
	f, ok := p.functionDefinitions[name]
	if ok {
		f.IsCalled = true
		p.functionDefinitions[name] = f
	}
}

func (p *Program) SetHaveBody(name string) {
	f, ok := p.functionDefinitions[name]
	if ok {
		f.HaveBody = true
		p.functionDefinitions[name] = f
	}
}

func (p *Program) GetCstdFunction() (src string) {
	for i := range p.functionDefinitions {
		if !p.functionDefinitions[i].IsCalled {
			continue
		}
		if !p.functionDefinitions[i].IsCstdFunction {
			continue
		}
		if p.functionDefinitions[i].pntCstd == nil {
			continue
		}
		// add dependencies of packages
		p.AddImports(p.functionDefinitions[i].pntCstd.dependPackages...)

		// add dependencies of functions
		for _, funcName := range p.functionDefinitions[i].pntCstd.dependFuncStd {
			for j := range p.functionDefinitions {
				if p.functionDefinitions[j].Name != funcName {
					continue
				}
				def := p.functionDefinitions[j]
				def.IsCalled = true
				p.functionDefinitions[j] = def
				break
			}
		}
	}

	for _, v := range p.functionDefinitions {
		if !v.IsCstdFunction {
			continue
		}
		if !v.IsCalled {
			continue
		}
		src += v.pntCstd.functionBody
	}
	return
}

func (p *Program) GetOutsideCalledFunctions() (ds []DefinitionFunction) {
	for _, v := range p.functionDefinitions {
		if v.IncludeFile == "" {
			continue
		}
		if !v.IsCalled {
			continue
		}
		if !v.HaveBody {
			ds = append(ds, v)
			continue
		}
		if v.IsCstdFunction {
			continue
		}
		if p.PreprocessorFile.IsUserSource(v.IncludeFile) {
			continue
		}
		ds = append(ds, v)
	}
	return
}
