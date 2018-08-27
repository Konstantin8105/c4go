[![Build Status](https://travis-ci.org/Konstantin8105/c4go.svg?branch=master)](https://travis-ci.org/Konstantin8105/c4go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/c4go)](https://goreportcard.com/report/github.com/Konstantin8105/c4go)
[![codecov](https://codecov.io/gh/Konstantin8105/c4go/branch/master/graph/badge.svg)](https://codecov.io/gh/Konstantin8105/c4go)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/Konstantin8105/c4go/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Konstantin8105/c4go?status.svg)](https://godoc.org/github.com/Konstantin8105/c4go)
[![Maintainability](https://api.codeclimate.com/v1/badges/b8d0bb5533207cce5ed3/maintainability)](https://codeclimate.com/github/Konstantin8105/c4go/maintainability)

A tool for [transpiling](https://en.wikipedia.org/wiki/Source-to-source_compiler) C code to Go code.

Milestone of project:

1. Transpiling project [GNU GSL](https://www.gnu.org/software/gsl/).
2. Transpiling project [GTK+](https://www.gtk.org/).

Notes:
* Transpiler works on linux and mac machines
* Need to have installed `clang`. See [llvm download page](http://releases.llvm.org/download.html)

# Installation

`c4go` requires Go 1.9 or newer.

```bash
go get -u github.com/Konstantin8105/c4go
```

# Example of using

```bash
# Change your location to folder with examples:
cd $GOPATH/src/github.com/Konstantin8105/c4go/examples/

# Transpile one file from C example folder:
c4go transpile prime.c

# Look on result
nano prime.go

# Check the result:
go run prime.go
# Enter a number
# 13
# The number is: 13
# Prime number.
```

C code of file `prime.c`:
```c
#include <stdio.h>

int main()
{
    int n, c;

    printf("Enter a number\n");
	// get value
    scanf("%d", &n);
    printf("The number is: %d\n", n);

	// ------- 
    if (n == 2)
        printf("Prime number.\n");
    else
    {
        for (c = 2; c <= n - 1; c++)
        {
            if (n % c == 0)
                break;
        }
        if (c != n)
            printf("Not prime.\n");
        else
            printf("Prime number.\n");
    }
    return 0;
}
```

Go code of file `prime.go`:
```golang
//
//	Package main - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package main

import "unsafe"
import "fmt"
import "github.com/Konstantin8105/c4go/noarch"

type size_t uint32
type __time_t int32
type va_list int64
type __gnuc_va_list int64

var stdin *noarch.File

var stdout *noarch.File

var stderr *noarch.File

// main - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/examples/prime.c:3
func main() {
	var n int
	var c int
	fmt.Printf("Enter a number\n")
	noarch.Scanf([]byte("%d\x00"), (*[100000000]int)(unsafe.Pointer(&n))[:])
	// get value
	//
	noarch.Printf([]byte("The number is: %d\n\x00"), n)
	if n == 2 {
		fmt.Printf("Prime number.\n")
		// -------
		//
	} else {
		for c = 2; c <= n-1; c++ {
			if n%c == 0 {
				break
			}
		}
		if c != n {
			fmt.Printf("Not prime.\n")
		} else {
			fmt.Printf("Prime number.\n")
		}
	}
	return
}
func init() {
	 stdin = noarch.Stdin
	 stdout = noarch.Stdout
	 stderr = noarch.Stderr
}
```

# C standart library implementation

```
            assert.h	       1/1	         100%
             ctype.h	     14/14	         100%
             errno.h	       0/1	           0%
             float.h	          	    undefined
            iso646.h	          	    undefined
            limits.h	          	    undefined
            locale.h	       0/3	           0%
              math.h	     33/58	        56.9%
            setjmp.h	       0/3	           0%
            signal.h	       0/3	           0%
            stdarg.h	       4/4	         100%
            stddef.h	       2/6	        33.3%
             stdio.h	     33/46	        71.7%
            stdlib.h	     31/47	          66%
            string.h	      9/24	        37.5%
              time.h	      8/15	        53.3%
             wchar.h	      0/68	           0%
            wctype.h	      0/22	           0%
```

# Contributing

Feel free to add PR, issues.
Main information from: [http://en.cppreference.com/w/c](http://en.cppreference.com/w/c)

## Testing

By default only unit tests are run with `go test`. You can also include the
integration tests:

```bash
go test -tags=integration ./...
```

Integration tests in the form of complete C programs that can be found in the
[tests](https://github.com/Konstantin8105/c4go/tree/master/tests) directory.

Integration tests work like this:

1. Clang compiles the C to a binary as normal.
2. c4go converts the C file to Go.
3. The Go is built to produce another binary.
4. Both binaries are executed and the output is compared. All C files will
contain some output so the results can be verified.
