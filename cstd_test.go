package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Coverage of C function

// http://www.cplusplus.com/reference/
var cstd = map[string][]string{
	"stdio.h": {
		"clearerr",
		"fclose",
		"feof",
		"ferror",
		"fflush",
		"fgetc",
		"fgetpos",
		"fgets",
		"fopen",
		"fprintf",
		"fputc",
		"fputs",
		"fread",
		"freopen",
		"fscanf",
		"fseek",
		"fsetpos",
		"ftell",
		"fwrite",
		"getc",
		"getchar",
		"gets",
		"perror",
		"printf",
		"putc",
		"putchar",
		"puts",
		"remove",
		"rename",
		"rewind",
		"scanf",
		"setbuf",
		"setvbuf",
		"snprintf",
		"sprintf",
		"sscanf",
		"tmpfile",
		"tmpnam",
		"ungetc",
		"vfprintf",
		"vfscanf",
		"vprintf",
		"vscanf",
		"vsnprintf",
		"vsprintf",
		"vsscanf",
	},
	"assert.h": {"assert"},
	"ctype.h": {
		"isalnum",
		"isalpha",
		"isblank",
		"iscntrl",
		"isdigit",
		"isgraph",
		"islower",
		"isprint",
		"ispunct",
		"isspace",
		"isupper",
		"isxdigit",
		"tolower",
		"toupper",
	},
	"errno.h":  {"errno"},
	"float.h":  {},
	"iso646.h": {},
	"limits.h": {},
	"locale.h": {
		"lconv", // "struct lconv",
		"setlocale",
		"localeconv",
	},
	"math.h": {
		"abs",
		"acos",
		"acosh",
		"asin",
		"asinh",
		"atan",
		"atan2",
		"atanh",
		"cbrt",
		"ceil",
		"copysign",
		"cos",
		"cosh",
		"erf",
		"erfc",
		"exp",
		"exp2",
		"expm1",
		"fabs",
		"fdim",
		"floor",
		"fma",
		"fmax",
		"fmin",
		"fmod",
		"frexp",
		"hypot",
		"ilogb",
		"ldexp",
		"lgamma",
		"llrint",
		"llround",
		"log",
		"log10",
		"log1p",
		"log2",
		"logb",
		"lrint",
		"lround",
		"modf",
		"nan",
		"nearbyint",
		"nextafter",
		"nexttoward",
		"pow",
		"remainder",
		"remquo",
		"rint",
		"round",
		"scalbln",
		"scalbn",
		"sin",
		"sinh",
		"sqrt",
		"tan",
		"tanh",
		"tgamma",
		"trunc",
	},
	"setjmp.h": {
		"jmp_buf",
		"longjmp",
		"setjmp",
	},
	"signal.h": {
		"raise",
		"sig_atomic_t",
		"signal",
	},
	"stdarg.h": {
		"va_arg",
		"va_end",
		"va_list",
		"va_start",
	},
	"stddef.h": {
		"NULL",
		"max_align_t",
		"nullptr_t",
		"offsetof",
		"ptrdiff_t",
		"size_t",
	},
	"stdlib.h": {
		"EXIT_FAILURE",
		"EXIT_SUCCESS",
		"MB_CUR_MAX",
		"NULL",
		"RAND_MAX",
		"_Exit",
		"abort",
		"abs",
		"at_quick_exit",
		"atexit",
		"atof",
		"atoi",
		"atol",
		"atoll",
		"bsearch",
		"calloc",
		"div",
		"div_t",
		"exit",
		"free",
		"getenv",
		"labs",
		"ldiv",
		"ldiv_t",
		"llabs",
		"lldiv",
		"lldiv_t",
		"malloc",
		"mblen",
		"mbstowcs",
		"mbtowc",
		"qsort",
		"quick_exit",
		"rand",
		"realloc",
		"size_t",
		"srand",
		"strtod",
		"strtof",
		"strtol",
		"strtold",
		"strtoll",
		"strtoul",
		"strtoull",
		"system",
		"wcstombs",
		"wctomb",
	},
	"string.h": {
		"NULL",
		"memchr",
		"memcmp",
		"memcpy",
		"memmove",
		"memset",
		"size_t",
		"strcat",
		"strchr",
		"strcmp",
		"strcoll",
		"strcpy",
		"strcspn",
		"strerror",
		"strlen",
		"strncat",
		"strncmp",
		"strncpy",
		"strpbrk",
		"strrchr",
		"strspn",
		"strstr",
		"strtok",
		"strxfrm",
	},
	"time.h": {
		"CLOCKS_PER_SEC",
		"NULL",
		"asctime",
		"clock",
		"clock_t",
		"ctime",
		"difftime",
		"gmtime",
		"localtime",
		"mktime",
		"size_t",
		"strftime",
		"time",
		"time_t",
		"tm", // "struct tm",
	},
	"wchar.h": {
		"NULL",
		"WCHAR_MAX",
		"WCHAR_MIN",
		"WEOF",
		"btowc",
		"fgetwc",
		"fgetws",
		"fputwc",
		"fputws",
		"fwide",
		"fwprintf",
		"fwscanf",
		"getwc",
		"getwchar",
		"mbrlen",
		"mbrtowc",
		"mbsinit",
		"mbsrtowcs",
		"mbstate_t",
		"putwc",
		"putwchar",
		"size_t",
		"swprintf",
		"swscanf",
		"tm", // "struct tm",
		"ungetwc",
		"vfwprintf",
		"vfwscanf",
		"vswprintf",
		"vswscanf",
		"vwprintf",
		"vwscanf",
		"wchar_t",
		"wcrtomb",
		"wcscat",
		"wcschr",
		"wcscmp",
		"wcscoll",
		"wcscpy",
		"wcscspn",
		"wcsftime",
		"wcslen",
		"wcsncat",
		"wcsncmp",
		"wcsncpy",
		"wcspbrk",
		"wcsrchr",
		"wcsrtombs",
		"wcsspn",
		"wcsstr",
		"wcstod",
		"wcstof",
		"wcstok",
		"wcstol",
		"wcstold",
		"wcstoll",
		"wcstoul",
		"wcstoull",
		"wcsxfrm",
		"wctob",
		"wint_t",
		"wmemchr",
		"wmemcmp",
		"wmemcpy",
		"wmemmove",
		"wmemset",
		"wprintf",
		"wscanf",
	},
	"wctype.h": {
		"WEOF",
		"iswalnum",
		"iswalpha",
		"iswblank",
		"iswcntrl",
		"iswctype",
		"iswdigit",
		"iswgraph",
		"iswlower",
		"iswprint",
		"iswpunct",
		"iswspace",
		"iswupper",
		"iswxdigit",
		"towctrans",
		"towlower",
		"towupper",
		"wctrans",
		"wctrans_t",
		"wctype",
		"wctype_t",
		"wint_t",
	},
}

func TestCSTD(t *testing.T) {
	// initialization
	amount := map[string]map[string]uint{}
	for include := range cstd {
		amount[include] = map[string]uint{}
		for i := range cstd[include] {
			amount[include][cstd[include][i]] = 0
		}
	}

	// calculation
	testFiles, err := filepath.Glob("tests/" + "*.c")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range testFiles {
		body, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("Cannot read file : %v\n", file)
		}
		// separate on parts
		body = bytes.Replace(body, []byte("["), []byte(" "), -1)
		body = bytes.Replace(body, []byte("]"), []byte(" "), -1)
		body = bytes.Replace(body, []byte("("), []byte(" "), -1)
		body = bytes.Replace(body, []byte(")"), []byte(" "), -1)
		body = bytes.Replace(body, []byte("="), []byte(" "), -1)
		body = bytes.Replace(body, []byte(";"), []byte(" "), -1)
		body = bytes.Replace(body, []byte(","), []byte(" "), -1)
		body = bytes.Replace(body, []byte("+"), []byte(" "), -1)
		body = bytes.Replace(body, []byte("-"), []byte(" "), -1)
		body = bytes.Replace(body, []byte("/"), []byte(" "), -1)
		body = bytes.Replace(body, []byte("*"), []byte(" "), -1)
		body = bytes.Replace(body, []byte("\n"), []byte(" "), -1)

		lines := bytes.Split(body, []byte(" "))
		for i := range lines {
			lines[i] = bytes.TrimSpace(lines[i])
		}

		for include := range amount {
			// check include file
			if !bytes.Contains(body, []byte("<"+include+">")) {
				continue
			}
			// finding function
			for _, function := range cstd[include] {
				for i := range lines {
					if bytes.Equal(lines[i], []byte(function)) {
						amount[include][function]++
					}
				}
			}
		}
	}

	// view
	type pair struct {
		inc  string
		line string
	}
	var ps []pair
	for include := range amount {
		var uniq uint
		for function := range amount[include] {
			if amount[include][function] > 0 {
				uniq++
			}
		}
		var p pair
		p.inc = include
		if len(amount[include]) > 0 {
			length := float64(uniq) / float64(len(amount[include]))
			p.line = fmt.Sprintf("%20s\t%10s\t%12.3g%s",
				include,
				fmt.Sprintf("%v/%v", uniq, len(amount[include])),
				length*100, "%")
		} else {
			p.line = fmt.Sprintf("%20s\t%10s\t%13s", include, "", "undefined")
		}
		ps = append(ps, p)
	}

	sort.Slice(ps, func(i, j int) bool {
		return strings.Compare(
			strings.TrimSpace(ps[i].inc),
			strings.TrimSpace(ps[j].inc)) == -1
	})

	// checking with README.md
	b, err := ioutil.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Cannot read file README.md : %v", err)
	}
	for _, l := range ps {
		if !bytes.Contains(b, []byte(l.line)) {
			t.Errorf("Please update information in file `README.md` about :\n`%s`",
				l.line)
		}
	}

	if flag.CommandLine.Lookup("test.v").Value.String() == "true" {
		for _, l := range ps {
			fmt.Fprintf(os.Stdout, "%s\n", l.line)
		}

		// Detail information
		fmt.Fprintln(os.Stdout, "\nDetail information:")
		for _, l := range ps {
			fmt.Fprintf(os.Stdout, "%s\n", l.line)
			var ps []pair
			for function := range amount[l.inc] {
				ps = append(ps, pair{
					inc:  function,
					line: fmt.Sprintf("\t%20s\t%v", function, amount[l.inc][function]),
				})
			}
			sort.Slice(ps, func(i, j int) bool {
				return strings.Compare(
					strings.TrimSpace(ps[i].inc),
					strings.TrimSpace(ps[j].inc)) == -1
			})

			for _, l := range ps {
				fmt.Fprintf(os.Stdout, "%s\n", l.line)
			}
		}
	}
}
