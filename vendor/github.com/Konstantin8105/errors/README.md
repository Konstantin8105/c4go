[![Coverage Status](https://coveralls.io/repos/github/Konstantin8105/errors/badge.svg?branch=master)](https://coveralls.io/github/Konstantin8105/errors?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/errors)](https://goreportcard.com/report/github.com/Konstantin8105/errors)
[![GoDoc](https://godoc.org/github.com/Konstantin8105/errors?status.svg)](https://godoc.org/github.com/Konstantin8105/errors)
![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)

# errors

Create error tree.

### Installation

```cmd
go get -u github.com/Konstantin8105/errors
```

### Example

```go
type ErrorValue struct {
	ValueName string
	Reason    error
}

func (e ErrorValue) Error() string {
	return fmt.Sprintf("Value `%s`: %v", e.ValueName, e.Reason)
}

func Example() {
	// some input data
	f := math.NaN()
	i := -32
	var s string

	// checking
	var et Tree
	et.Name = "Check input data"
	if math.IsNaN(f) {
		et.Add(ErrorValue{
			ValueName: "f",
			Reason:    fmt.Errorf("is NaN"),
		})
	}
	if f < 0 {
		et.Add(fmt.Errorf("Parameter `f` is negative"))
	}
	if i < 0 {
		et.Add(fmt.Errorf("Parameter `i` is less zero"))
	}
	if s == "" {
		et.Add(fmt.Errorf("Parameter `s` is empty"))
	}

	if et.IsError() {
		fmt.Println(et.Error())
	}

	// walk
	Walk(&et, func(e error) {
		fmt.Fprintf(os.Stdout, "%-25s %v\n", fmt.Sprintf("%T", e), e)
	})

	// Output:
	// Check input data
	// ├──Value `f`: is NaN
	// ├──Parameter `i` is less zero
	// └──Parameter `s` is empty
	//
	// errors.ErrorValue         Value `f`: is NaN
	// *errors.errorString       Parameter `i` is less zero
	// *errors.errorString       Parameter `s` is empty
}
```

Acceptable add in error tree another error tree and possibly look like that:

```
+
├── Error 0
├── +
│   ├── Inside error 0
│   └── Some deep deep errors
│       └── Deep error 0
├── Error 1
├── Error 2
├── Error 3
├── +
│   ├── Inside error 0
│   ├── Some deep deep errors
│   │   └── Deep error 0
│   └── Inside error 1
├── Error 4
├── Error 5
├── Error 6
├── +
│   ├── Inside error 0
│   ├── Some deep deep errors
│   │   └── Deep error 0
│   ├── Inside error 1
│   ├── Inside error 2
│   └── Some deep deep errors
│       ├── Deep error 0
│       └── Deep error 1
├── Error 7
├── Error 8
├── Error 9
└── +
    ├── Inside error 0
    ├── Some deep deep errors
    │   └── Deep error 0
    ├── Inside error 1
    ├── Inside error 2
    ├── Some deep deep errors
    │   ├── Deep error 0
    │   └── Deep error 1
    └── Inside error 3
```
