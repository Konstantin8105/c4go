[![codecov](https://codecov.io/gh/Konstantin8105/tree/branch/master/graph/badge.svg)](https://codecov.io/gh/Konstantin8105/tree)
[![Build Status](https://travis-ci.org/Konstantin8105/tree.svg?branch=master)](https://travis-ci.org/Konstantin8105/tree)
[![Go Report Card](https://goreportcard.com/badge/github.com/Konstantin8105/tree)](https://goreportcard.com/report/github.com/Konstantin8105/tree)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/Konstantin8105/tree/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/Konstantin8105/tree?status.svg)](https://godoc.org/github.com/Konstantin8105/tree)

# tree

`tree` is tree viewer.

minimal example:
```
func ExampleTree() {
	artist := tree.New("Pantera")
	album := tree.New("Far Beyond Driven")
	album.Add("5 minutes Alone")
	album.Add("Some another")
	artist.Add(album)
	artist.Add("Power Metal")
	fmt.Fprintf(os.Stdout, "%s\n", artist)

	// Output:
	// Pantera
	// ├──Far Beyond Driven
	// │  ├──5 minutes Alone
	// │  └──Some another
	// └──Power Metal
}
```
