[![Build Status](https://travis-ci.org/Konstantin8105/c4go.svg?branch=master)](https://travis-ci.org/Konstantin8105/c4go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/c4go)](https://goreportcard.com/report/github.com/Konstantin8105/c4go)
[![codecov](https://codecov.io/gh/Konstantin8105/c4go/branch/master/graph/badge.svg)](https://codecov.io/gh/Konstantin8105/c4go)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/Konstantin8105/c4go/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Konstantin8105/c4go?status.svg)](https://godoc.org/github.com/Konstantin8105/c4go)

A tool for converting C to Go.

The goals of this project are:

1. To create a generic tool that can convert C to Go.
2. To be cross platform (linux and mac) and work against as many clang versions
as possible (the clang AST API is not stable).
3. To be a repeatable and predictable tool (rather than doing most of the work
and you have to clean up the output to get it working.)
4. To deliver quick and small version increments.
5. The ultimate milestone is to be able to compile the
[SQLite3 source code](https://sqlite.org/download.html) and have it working
without modification. This will be the 1.0.0 release.

# Installation

```bash
go get -u github.com/Konstantin8105/c4go
```

# Usage

```bash
c2go transpile myfile.c
```

The `c2go` program processes a single C file and outputs the translated code
in Go. Let's use an included example,
[prime.c](https://github.com/Konstantin8105/c4go/blob/master/examples/prime.c):

```c
#include <stdio.h>
 
int main()
{
   int n, c;
 
   printf("Enter a number\n");
   scanf("%d", &n);
 
   if ( n == 2 )
      printf("Prime number.\n");
   else
   {
       for ( c = 2 ; c <= n - 1 ; c++ )
       {
           if ( n % c == 0 )
              break;
       }
       if ( c != n )
          printf("Not prime.\n");
       else
          printf("Prime number.\n");
   }
   return 0;
}
```

```bash
c2go transpile prime.c
go run prime.go
```

```
Enter a number
23
Prime number.
```

`prime.go` looks like:

```go
package main

import "unsafe"

import "github.com/Konstantin8105/c4go/noarch"

// ... lots of system types in Go removed for brevity.

var stdin *noarch.File
var stdout *noarch.File
var stderr *noarch.File

func main() {
	__init()
	var n int
	var c int
	noarch.Printf([]byte("Enter a number\n\x00"))
	noarch.Scanf([]byte("%d\x00"), (*[1]int)(unsafe.Pointer(&n))[:])
	if n == 2 {
		noarch.Printf([]byte("Prime number.\n\x00"))
	} else {
		for c = 2; c <= n-1; func() int {
			c += 1
			return c
		}() {
			if n%c == 0 {
				break
			}
		}
		if c != n {
			noarch.Printf([]byte("Not prime.\n\x00"))
		} else {
			noarch.Printf([]byte("Prime number.\n\x00"))
		}
	}
	return
}

func __init() {
	stdin = noarch.Stdin
	stdout = noarch.Stdout
	stderr = noarch.Stderr
}
```

# What Is Supported?

See the
[Project Progress](https://github.com/Konstantin8105/c4go/wiki/Project-Progress).

# How It Works

This is the process:

1. The C code is preprocessed with clang. This generates a larger file (`pp.c`),
but removes all the platform specific directives and macros.

2. `pp.c` is parsed with the clang AST and dumps it in a colourful text format
that
[looks like this](http://ehsanakhgari.org/wp-content/uploads/2015/12/Screen-Shot-2015-12-03-at-5.02.38-PM.png).
Apart from just parsing the C and dumping an AST, the AST contains all of the
resolved information that a compiler would need (such as data types). This means
that the code must compile successfully under clang for the AST to also be
usable.

3. Since we have all the types in the AST it's just a matter of traversing the
tree is a semi-intelligent way and producing Go. Easy, right!?

# Testing

By default only unit tests are run with `go test`. You can also include the
integration tests:

```bash
go test -tags=integration ./...
```

Integration tests in the form of complete C programs that can be found in the
[tests](https://github.com/Konstantin8105/c4go/tree/master/tests) directory.

Integration tests work like this:

1. Clang compiles the C to a binary as normal.
2. c2go converts the C file to Go.
3. The Go is built to produce another binary.
4. Both binaries are executed and the output is compared. All C files will
contain some output so the results can be verified.

# Contributing

Contributing is done with pull requests. There is no help that is too small! :)

If you're looking for where to start I can suggest
[finding a simple C program](http://www.programmingsimplified.com/c-program-examples)
(like the other examples) that do not successfully translate into Go.

Or, if you don't want to do that you can submit it as an issue so that it can be
picked up by someone else.
