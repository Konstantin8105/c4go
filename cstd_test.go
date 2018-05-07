package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Coverage of C function

// http://www.cplusplus.com/reference/
var cstd = map[string][]string{
	"stdio.h": []string{
		"remove",
		"rename",
		"tmpfile",
		"tmpnam",
		"fclose",
		"fflush",
		"fopen",
		"freopen",
		"setbuf",
		"setvbuf",
		"fprintf",
		"fscanf",
		"printf",
		"scanf",
		"snprintf",
		"sprintf",
		"sscanf",
		"vfprintf",
		"vfscanf",
		"vprintf",
		"vscanf",
		"vsnprintf",
		"vsprintf",
		"vsscanf",
		"fgetc",
		"fgets",
		"fputc",
		"fputs",
		"getc",
		"getchar",
		"gets",
		"putc",
		"putchar",
		"puts",
		"ungetc",
		"fread",
		"fwrite",
		"fgetpos",
		"fseek",
		"fsetpos",
		"ftell",
		"rewind",
		"clearerr",
		"feof",
		"ferror",
		"perror",
	},
	"assert.h": []string{"assert"},
	"ctype.h": []string{
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
	"errno.h":  []string{"errno"},
	"float.h":  []string{},
	"iso646.h": []string{},
	"limits.h": []string{},
	"locale.h": []string{
		"struct lconv",
		"setlocale",
		"localeconv",
	},
	"math.h": []string{
		"cos",
		"sin",
		"tan",
		"acos",
		"asin",
		"atan",
		"atan2",
		"cosh",
		"sinh",
		"tanh",
		"acosh",
		"asinh",
		"atanh",
		"exp",
		"frexp",
		"ldexp",
		"log",
		"log10",
		"modf",
		"exp2",
		"expm1",
		"ilogb",
		"log1p",
		"log2",
		"logb",
		"scalbn",
		"scalbln",
		"pow",
		"sqrt",
		"cbrt",
		"hypot",
		"erf",
		"erfc",
		"tgamma",
		"lgamma",
		"ceil",
		"floor",
		"fmod",
		"trunc",
		"round",
		"lround",
		"llround",
		"rint",
		"lrint",
		"llrint",
		"nearbyint",
		"remainder",
		"remquo",
		"copysign",
		"nan",
		"nextafter",
		"nexttoward",
		"fdim",
		"fmax",
		"fmin",
		"fabs",
		"abs",
		"fma",
	},
	"setjmp.h": []string{
		"longjmp",
		"setjmp",
		"jmp_buf",
	},
	"signal.h": []string{
		"signal",
		"raise",
		"sig_atomic_t",
	},
	"stdarg.h": []string{
		"va_list",
		"va_start",
		"va_arg",
		"va_end",
	},
	"stddef.h": []string{
		"ptrdiff_t",
		"size_t",
		"max_align_t",
		"nullptr_t",
		"offsetof",
		"NULL",
	},
	"stdlib.h": []string{
		"atof",
		"atoi",
		"atol",
		"atoll",
		"strtod",
		"strtof",
		"strtol",
		"strtold",
		"strtoll",
		"strtoul",
		"strtoull",
		"rand",
		"srand",
		"calloc",
		"free",
		"malloc",
		"realloc",
		"abort",
		"atexit",
		"at_quick_exit",
		"exit",
		"getenv",
		"quick_exit",
		"system",
		"_Exit",
		"bsearch",
		"qsort",
		"abs",
		"div",
		"labs",
		"ldiv",
		"llabs",
		"lldiv",
		"mblen",
		"mbtowc",
		"wctomb",
		"mbstowcs",
		"wcstombs",
		"EXIT_FAILURE",
		"EXIT_SUCCESS",
		"MB_CUR_MAX",
		"NULL",
		"RAND_MAX",
		"div_t",
		"ldiv_t",
		"lldiv_t",
		"size_t",
	},
	"string.h": []string{
		"memcpy",
		"memmove",
		"strcpy",
		"strncpy",
		"strcat",
		"strncat",
		"memcmp",
		"strcmp",
		"strcoll",
		"strncmp",
		"strxfrm",
		"memchr",
		"strchr",
		"strcspn",
		"strpbrk",
		"strrchr",
		"strspn",
		"strstr",
		"strtok",
		"memset",
		"strerror",
		"strlen",
		"NULL",
		"size_t",
	},
	"time.h": []string{
		"clock",
		"difftime",
		"mktime",
		"time",
		"asctime",
		"ctime",
		"gmtime",
		"localtime",
		"strftime",
		"CLOCKS_PER_SEC",
		"NULL",
		"clock_t",
		"size_t",
		"time_t",
		"struct tm",
	},
	"wchar.h": []string{

		"fgetwc",
		"fgetws",
		"fputwc",
		"fputws",
		"fwide",
		"fwprintf",
		"fwscanf",
		"getwc",
		"getwchar",
		"putwc",
		"putwchar",
		"swprintf",
		"swscanf",
		"ungetwc",
		"vfwprintf",
		"vfwscanf",
		"vswprintf",
		"vswscanf",
		"vwprintf",
		"vwscanf",
		"wprintf",
		"wscanf",
		"wcstod",
		"wcstof",
		"wcstol",
		"wcstold",
		"wcstoll",
		"wcstoul",
		"wcstoull",
		"btowc",
		"mbrlen",
		"mbrtowc",
		"mbsinit",
		"mbsrtowcs",
		"wcrtomb",
		"wctob",
		"wcsrtombs",
		"wcscat",
		"wcschr",
		"wcscmp",
		"wcscoll",
		"wcscpy",
		"wcscspn",
		"wcslen",
		"wcsncat",
		"wcsncmp",
		"wcsncpy",
		"wcspbrk",
		"wcsrchr",
		"wcsspn",
		"wcsstr",
		"wcstok",
		"wcsxfrm",
		"wmemchr",
		"wmemcmp",
		"wmemcpy",
		"wmemmove",
		"wmemset",
		"wcsftime",
		"mbstate_t",
		"size_t",
		"struct tm",
		"wchar_t",
		"wint_t",
		"NULL",
		"WCHAR_MAX",
		"WCHAR_MIN",
		"WEOF",
	},
	"wctype.h": []string{

		"iswalnum",
		"iswalpha",
		"iswblank",
		"iswcntrl",
		"iswdigit",
		"iswgraph",
		"iswlower",
		"iswprint",
		"iswpunct",
		"iswspace",
		"iswupper",
		"iswxdigit",
		"towlower",
		"towupper",
		"iswctype",
		"towctrans",
		"wctrans",
		"wctype",
		"wctrans_t",
		"wctype_t",
		"wint_t",
		"WEOF",
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
	testFiles, err := filepath.Glob("tests/*.c")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range testFiles {
		body, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("Cannot read file : %v\n", file)
		}
		for include := range amount {
			// check include file
			if !bytes.Contains(body, []byte("<"+include+">")) {
				continue
			}
			for _, function := range cstd[include] {
				if bytes.Contains(body, []byte(function)) {
					amount[include][function]++
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
		var a uint
		var uniq uint
		for function := range amount[include] {
			a += amount[include][function]
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
		// Detail information
		// for function := range amount[include] {
		// 	fmt.Printf("\t%20s\t%v\n", function, amount[include][function])
		// }
	}

	sort.Slice(ps, func(i, j int) bool {
		return strings.Compare(
			strings.TrimSpace(ps[i].inc),
			strings.TrimSpace(ps[j].inc)) == -1
	})

	for _, l := range ps {
		fmt.Printf("%s\n", l.line)
	}
}
