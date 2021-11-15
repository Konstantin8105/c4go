[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/c4go)](https://goreportcard.com/report/github.com/Konstantin8105/c4go)
[![codecov](https://codecov.io/gh/Konstantin8105/c4go/branch/master/graph/badge.svg)](https://codecov.io/gh/Konstantin8105/c4go)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/Konstantin8105/c4go/master/LICENSE)
[![Go](https://github.com/Konstantin8105/c4go/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Konstantin8105/c4go/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Konstantin8105/c4go.svg)](https://pkg.go.dev/github.com/Konstantin8105/c4go)


A tool for [transpiling](https://en.wikipedia.org/wiki/Source-to-source_compiler) C code to Go code.

Milestones of the project:

1. Transpiling project [GNU GSL](https://www.gnu.org/software/gsl/).
2. Transpiling project [GTK+](https://www.gtk.org/).

Notes:
* Transpiler works on linux machines
* Need to have installed `clang`. See [llvm download page](http://releases.llvm.org/download.html)



# Installation

Installation with version generation.

```bash
# get code
go get -u github.com/Konstantin8105/c4go

# move to project source
cd $GOPATH/src/github.com/Konstantin8105/c4go

# generate version
go generate ./...

# install
go install

# testing
c4go version
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

// main - transpiled function from  C4GO/examples/prime.c:3
func main() {
	var n int32
	var c int32
	fmt.Printf("Enter a number\n")
	// get value
	noarch.Scanf([]byte("%d\x00"), c4goUnsafeConvert_int32(&n))
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

// c4goUnsafeConvert_int32 : created by c4go
func c4goUnsafeConvert_int32(c4go_name *int32) []int32 {
	return (*[1000000]int32)(unsafe.Pointer(c4go_name))[:]
}
```

# Example with binding function

C:

```c
#include <math.h>
#include <stdio.h>

int main()
{
    int n;
    double param = 8.0, result;
    result = frexp(param, &n);
    printf("result = %5.2f\n", result);
    printf("n      = %d\n", n);
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

// #include </usr/include/math.h>
import "C"

import "github.com/Konstantin8105/c4go/noarch"
import "unsafe"

// main - transpiled function from  C4GO/examples/math.c:4
func main() {
	var n int32
	var param float64 = 8
	var result float64
	result = frexp(param, c4goUnsafeConvert_int32(&n))
	noarch.Printf([]byte("result = %5.2f\n\x00"), result)
	noarch.Printf([]byte("n      = %d\n\x00"), n)
	return
}

// c4goUnsafeConvert_int32 : created by c4go
func c4goUnsafeConvert_int32(c4go_name *int32) []int32 {
	return (*[1000000]int32)(unsafe.Pointer(c4go_name))[:]
}

// frexp - add c-binding for implemention function
func frexp(arg0 float64, arg1 []int32) float64 {
 	return float64(C.frexp(C.double(arg0), (*C.int)(unsafe.Pointer(&arg1[0]))))
}
```

# Example with C-pointers and C-arrays

```c
#include <stdio.h>

// input argument - C-pointer
void a(int* v1) { printf("a: %d\n", *v1); }

// input argument - C-array
void b(int v1[], int size)
{
    for (size--; size >= 0; size--) {
        printf("b: %d %d\n", size, v1[size]);
    }
}

int main()
{
    // value
    int i1 = 42;
    a(&i1);
    b(&i1, 1);

    // C-array
    int i2[] = { 11, 22 };
    a(i2);
    b(i2, 2);

    // C-pointer from value
    int* i3 = &i1;
    a(i3);
    b(i3, 1);

    // C-pointer from array
    int* i4 = i2;
    a(i4);
    b(i4, 2);

    // C-pointer from array
    int* i5 = i2[1];
    a(i5);
    b(i5, 1);

    return 0;
}
```

```go
//
//	Package - transpiled by c4go
//
//	If you have found any issues, please raise an issue at:
//	https://github.com/Konstantin8105/c4go/
//

package main

import "unsafe"
import "github.com/Konstantin8105/c4go/noarch"

// a - transpiled function from  C4GO/examples/ap.c:4
func a(v1 []int32) {
	// input argument - C-pointer
	noarch.Printf([]byte("a: %d\n\x00"), v1[0])
}

// b - transpiled function from  C4GO/examples/ap.c:7
func b(v1 []int32, size int32) {
	{
		// input argument - C-array
		for size -= 1; size >= 0; size-- {
			noarch.Printf([]byte("b: %d %d\n\x00"), size, v1[size])
		}
	}
}

// main - transpiled function from  C4GO/examples/ap.c:14
func main() {
	var i1 int32 = 42
	// value
	a(c4goUnsafeConvert_int32(&i1))
	b(c4goUnsafeConvert_int32(&i1), 1)
	var i2 []int32 = []int32{11, 22}
	// C-array
	a(i2)
	b(i2, 2)
	var i3 []int32 = c4goUnsafeConvert_int32(&i1)
	// C-pointer from value
	a(i3)
	b(i3, 1)
	var i4 []int32 = i2
	// C-pointer from array
	a(i4)
	b(i4, 2)
	var i5 []int32 = i2[1:]
	// C-pointer from array
	a(i5)
	b(i5, 1)

	return
}

// c4goUnsafeConvert_int32 : created by c4go
func c4goUnsafeConvert_int32(c4go_name *int32) []int32 {
	return (*[1000000]int32)(unsafe.Pointer(c4go_name))[:]
}
```

# Example of transpilation [neatvi](https://github.com/aligrudi/neatvi).

```bash
# clone git repository
git clone https://github.com/aligrudi/neatvi.git

# move in folder
cd neatvi

# look in Makefile
cat Makefile
```
Makefile:
```
CC = cc
CFLAGS = -Wall -O2
LDFLAGS =

OBJS = vi.o ex.o lbuf.o mot.o sbuf.o ren.o dir.o syn.o reg.o led.o \
	uc.o term.o rset.o regex.o cmd.o conf.o

all: vi

conf.o: conf.h

%.o: %.c
	$(CC) -c $(CFLAGS) $<
vi: $(OBJS)
	$(CC) -o $@ $(OBJS) $(LDFLAGS)
clean:
rm -f *.o vi
```
As we see not need any specific, for example `-clang-flag`.

Let's try transpile:
```bash
c4go transpile *.c
```
Result:
```
panic: clang failed: exit status 1:

/temp/neatvi/reg.c:6:14: error: redefinition of 'bufs' with a different type: 'char *[256]' vs 'struct buf [8]'
static char *bufs[256];
             ^
/temp/neatvi/ex.c:35:3: note: previous definition is here
} bufs[8];
  ^
/temp/neatvi/reg.c:12:9: error: returning 'struct buf' from a function with incompatible result type 'char *'
 return bufs[c];
        ^~~~~~~
/temp/neatvi/reg.c:17:81: error: invalid operands to binary expression ('int' and 'struct buf')
 char *pre = ((*__ctype_b_loc ())[(int) ((c))] & (unsigned short int) _ISupper) && bufs[tolower(c)] ? bufs[tolower(c)] : "";
             ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ ^  ~~~~~~~~~~~~~~~~
/temp/neatvi/reg.c:21:7: error: passing 'struct buf' to parameter of incompatible type 'void *'
 free(bufs[tolower(c)]);
      ^~~~~~~~~~~~~~~~
/usr/include/stdlib.h:483:25: note: passing argument to parameter '__ptr' here
extern void free (void *__ptr) __attribute__ ((__nothrow__ ));
                        ^
/temp/neatvi/reg.c:22:19: error: assigning to 'struct buf' from incompatible type 'char *'
 bufs[tolower(c)] = buf;
                  ^ ~~~
/temp/neatvi/reg.c:43:8: error: passing 'struct buf' to parameter of incompatible type 'void *'
  free(bufs[i]);
       ^~~~~~~
/usr/include/stdlib.h:483:25: note: passing argument to parameter '__ptr' here
extern void free (void *__ptr) __attribute__ ((__nothrow__ ));
                        ^
/temp/neatvi/regex.c:98:12: error: static declaration of 'uc_len' follows non-static declaration
static int uc_len(char *s)
           ^
/temp/neatvi/vi.h:83:5: note: previous declaration is here
int uc_len(char *s);
    ^
/temp/neatvi/regex.c:200:14: error: static declaration of 'uc_beg' follows non-static declaration
static char *uc_beg(char *beg, char *s)
             ^
/temp/neatvi/vi.h:101:7: note: previous declaration is here
char *uc_beg(char *beg, char *s);
      ^
/temp/neatvi/vi.c:635:12: error: redefinition of 'linecount'
static int linecount(char *s)
           ^
/temp/neatvi/lbuf.c:124:12: note: previous definition is here
static int linecount(char *s)
           ^
9 errors generated.


goroutine 1 [running]:
main.generateAstLines(0x1, 0x0, 0xc000010140, 0x10, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/github.com/Konstantin8105/c4go/main.go:438 +0xd40
main.Start(0x1, 0x0, 0xc000010140, 0x10, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	/go/src/github.com/Konstantin8105/c4go/main.go:356 +0xa9
main.runCommand(0x0)
	/go/src/github.com/Konstantin8105/c4go/main.go:723 +0xa45
main.main()
	/go/src/github.com/Konstantin8105/c4go/main.go:556 +0x22
```
That error is good and enought information for next step of transpilation.
Let's clarify one of them - we see 2 function with same function name and
if you open files and compare.
```
/temp/neatvi/vi.c:635:12: error: redefinition of 'linecount'
static int linecount(char *s)

/temp/neatvi/lbuf.c:124:12: note: previous definition is here
static int linecount(char *s)
```
This is 2 absolute identical function, so we can remove one of them.
I choose remove function `linecount` from file `vi.c`, because in
according to error - it is duplicate.

Also remove next function in according to errors:
* `vi.c` function `linecount`
* `regex.c` function `uc_len`
* `regex.c` function `uc_bed`

Next error is about duplicate of variable names:
```
/temp/neatvi/reg.c:6:14: error: redefinition of 'bufs' with a different type: 'char *[256]' vs 'struct buf [8]'
static char *bufs[256];
             ^
/temp/neatvi/ex.c:35:3: note: previous definition is here
} bufs[8];
```
Let's replace variable name in file `reg.c` and run:
```
sed -i.bak 's/bufs/bufs_postfix/g' reg.c
```

Now, let's try again and output in file `vi.go`:
```bash
c4go transpile -o vi.go *.c
```
Now transpiled without `clang` error. Look the result:
```bash
less vi.go
```
We see warning of lose some types:
```
// Warning (*ast.BinaryOperator):  /temp/neatvi/cmd.c:20 :Cannot transpile BinaryOperator with type 'int' : result type = {}. Error: operator is `=`. Cannot casting {__pid_t -> int}. err = Cannot resolve type '__pid_t' : I couldn't find an appropriate Go type for the C type '__pid_t'.
```

For fix that add flag `-s` for adding all type from C standard library:
```bash
c4go transpile -o vi.go -s *.c
```

Now, all is Ok. Look the result:
```bash
less vi.go
```


# C standard library implementation

```
            assert.h	       1/1	         100%
             ctype.h	     13/13	         100%
             errno.h	       0/1	           0%
             float.h	          	    undefined
            iso646.h	          	    undefined
            limits.h	          	    undefined
            locale.h	       0/3	           0%
              math.h	     22/22	         100%
            setjmp.h	       0/3	           0%
            signal.h	       3/3	         100%
            stdarg.h	       4/4	         100%
            stddef.h	       4/4	         100%
             stdio.h	     37/41	        90.2%
            stdlib.h	     31/37	        83.8%
            string.h	     21/24	        87.5%
              time.h	      1/15	        6.67%
             wchar.h	      3/60	           5%
            wctype.h	      0/21	           0%
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
time ./scripts/sqlite.sh 

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

