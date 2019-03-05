[![Build Status](https://travis-ci.org/Konstantin8105/c4go.svg?branch=master)](https://travis-ci.org/Konstantin8105/c4go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/c4go)](https://goreportcard.com/report/github.com/Konstantin8105/c4go)
[![codecov](https://codecov.io/gh/Konstantin8105/c4go/branch/master/graph/badge.svg)](https://codecov.io/gh/Konstantin8105/c4go)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/Konstantin8105/c4go/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Konstantin8105/c4go?status.svg)](https://godoc.org/github.com/Konstantin8105/c4go)
[![Maintainability](https://api.codeclimate.com/v1/badges/b8d0bb5533207cce5ed3/maintainability)](https://codeclimate.com/github/Konstantin8105/c4go/maintainability)

A tool for [transpiling](https://en.wikipedia.org/wiki/Source-to-source_compiler) C code to Go code.

Milestones of the project:

1. Transpiling project [GNU GSL](https://www.gnu.org/software/gsl/).
2. Transpiling project [GTK+](https://www.gtk.org/).

Notes:
* Transpiler works on linux machines
* Need to have installed `clang`. See [llvm download page](http://releases.llvm.org/download.html)

# Installation

`c4go` requires Go 1.10 or newer.

```bash
go get -u github.com/Konstantin8105/c4go
```

# Usage example

```bash
# Change your location to the folder with examples:
cd $GOPATH/src/github.com/Konstantin8105/c4go/examples/

# Transpile one file from the C example folder:
c4go transpile prime.c

# Look at the result
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
    else {
        for (c = 2; c <= n - 1; c++) {
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
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package main

import "unsafe"
import "github.com/Konstantin8105/c4go/noarch"
import "fmt"

// main - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/examples/prime.c:3
func main() {
	var n int
	var c int
	fmt.Printf("Enter a number\n")
	// get value
	noarch.Scanf([]byte("%d\x00"), (*[100000000]int)(unsafe.Pointer(&n))[:])
	noarch.Printf([]byte("The number is: %d\n\x00"), n)
	if n == 2 {
		// -------
		fmt.Printf("Prime number.\n")
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
```

# Example with binding function

Run:

```
c4go transpile 
```


```c
#include <stdio.h>
#include <math.h>

int main()
{
	int n;
	double param = 8.0, result;
	result = frexp(param,&n);
	printf("result = %5.2f\n", result);
	printf("n      = %d\n"   , n     );
	return 0;
}
```

`c4go` add automatically C binding for function without implementation:
```golang
//
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package main

// #include </usr/include/x86_64-linux-gnu/bits/mathcalls.h>
import "C"

import "github.com/Konstantin8105/c4go/noarch"
import "unsafe"

// main - transpiled function from  $GOPATH/src/github.com/Konstantin8105/c4go/examples/math.c:4
func main() {
	var n int
	var param float64 = 8
	var result float64
	result = frexp(param, (*[100000000]int)(unsafe.Pointer(&n))[:])
	noarch.Printf([]byte("result = %5.2f\n\x00"), result)
	noarch.Printf([]byte("n      = %d\n\x00"), n)
	return
}

func frexp(arg0 float64, arg1 []int) float64 {
	return float64(C.frexp(C.double(arg0), (*_Ctype_int)(unsafe.Pointer(&arg1[0]))))
}
```




# C standard library implementation

```
            assert.h	       1/1	         100%
             ctype.h	     14/14	         100%
             errno.h	       0/1	           0%
             float.h	          	    undefined
            iso646.h	          	    undefined
            limits.h	          	    undefined
            locale.h	       3/3	         100%
              math.h	     38/58	        65.5%
            setjmp.h	       0/3	           0%
            signal.h	       3/3	         100%
            stdarg.h	       4/4	         100%
            stddef.h	       2/6	        33.3%
             stdio.h	     35/46	        76.1%
            stdlib.h	     33/47	        70.2%
            string.h	     16/24	        66.7%
              time.h	     13/15	        86.7%
             wchar.h	      3/68	        4.41%
            wctype.h	      0/22	           0%
```

# Contributing

Feel free to submit PRs or open issues.
Main information from: [en.cppreference.com](http://en.cppreference.com/w/c)

## Testing

By default, only unit tests are run with `go test`. You can also include the
integration tests:

```bash
go test -tags=integration ./...
```

Integration tests in the form of complete C programs that can be found in the
[tests](https://github.com/Konstantin8105/c4go/tree/master/tests) directory.

Integration tests work like this:

1. `clang` compiles the C to a binary as normal.
2. `c4go` converts the C file to Go.
3. Both binaries are executed and the output is compared. All C files will
contain some output so the results can be verified.

## Note

### Use lastest version of clang.

If you use `Ubuntu`, then use command like next for choose `clang` version:

```bash
sudo update-alternatives --install /usr/bin/clang clang /usr/bin/clang-6.0 1000
```

### Performance

Main time of transpilation takes `clang`, for example run:
```bash
go test -tags=integration -run=Benchmark -bench=. -benchmem
```
Result looks for example:
```
goos: linux
goarch: amd64
pkg: github.com/Konstantin8105/c4go
BenchmarkTranspile/Full-6         	       5	 274922964 ns/op	43046865 B/op	  379676 allocs/op
BenchmarkTranspile/GoCode-6       	      20	  86806808 ns/op	36577533 B/op	  308060 allocs/op
PASS
```
So, transpilation time is just 30% of full time. In my point of view
no need of performance optimization, see [Amdahl's law](https://en.wikipedia.org/wiki/Amdahl%27s_law).

### Example of performance analyse

Please run:

```bash
# Run cpuprofiling for sqlite transpilation example
time ./travis/sqlite.sh 

# Example of output:
#
# % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
#                                 Dload  Upload   Total   Spent    Left  Speed
#100 2217k  100 2217k    0     0   235k      0  0:00:09  0:00:09 --:--:--  357k
#Archive:  /tmp/SQLITE/sqlite-amalgamation-3250200.zip
#   creating: /tmp/SQLITE/sqlite-amalgamation-3250200/
#  inflating: /tmp/SQLITE/sqlite-amalgamation-3250200/sqlite3ext.h  
#  inflating: /tmp/SQLITE/sqlite-amalgamation-3250200/sqlite3.c  
#  inflating: /tmp/SQLITE/sqlite-amalgamation-3250200/sqlite3.h  
#  inflating: /tmp/SQLITE/sqlite-amalgamation-3250200/shell.c  
#After transpiling shell.c and sqlite3.c together, have summary: 695 warnings.
#In file sqlite.go summary : 3 warnings in go build.
#Amount unsafe package using: 2902
#
#real	0m18.434s
#user	0m14.212s
#sys	0m1.434s

# Run profiler
go tool pprof ./testdata/cpu.out
```

For more information, see [Profiling Go Programs](https://blog.golang.org/profiling-go-programs).
