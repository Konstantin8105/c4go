package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
)

func unused(filenames ...string) {
	for _, src := range filenames {
		fset := token.NewFileSet() // positions are relative to fset

		// Parse src but stop after processing the imports.
		f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			return
		}

		// map
		m := map[string][]string{}

		// Print the imports from the file's AST.
		for i := range f.Decls {
			decls := f.Decls[i]
			if fd, ok := decls.(*ast.FuncDecl); ok {
				name := fd.Name.Name
				var cs callSearcher
				ast.Walk(&cs, fd)
				m[name] = cs.used
			}
		}

		// list of all function
		used := map[string]bool{}
		list := []string{"main"}

		for iter := len(m) * 10; iter > 0; iter-- {
		rem:
			sort.Strings(list)
			for i := range list {
				if i == 0 {
					continue
				}
				if list[i-1] == list[i] {
					// remove from list
					list = append(list[:i], list[i+1:]...)
					goto rem
				}
			}
			for i := range list {
				_, ok := used[list[i]]
				if ok {
					// remove from list
					list = append(list[:i], list[i+1:]...)
					goto rem
				}
			}
			var newlist []string
			for k, v := range m {
				for i := range list {
					used[list[i]] = true
					if k == list[i] {
						used[k] = true
						newlist = append(newlist, v...)
					}
				}
			}
			list = newlist
		}

		// full list
		full := map[string]bool{}
		for k, v := range m {
			full[k] = true
			for i := range v {
				full[v[i]] = true
			}
		}

		// unused
		for k := range full {
			_, ok := used[k]
			if ok {
				continue
			}
			fmt.Fprintf(os.Stdout, "%s\n", k)
		}
	}
}

type callSearcher struct {
	used []string
}

var goFuncs = []string{
	"len", "make", "append",
	"string",
	"float64", "float32",
	"int", "int64", "int32", "int16",
	"go", "close",
	"panic",
}

// Visit for walking by node tree
func (c *callSearcher) Visit(node ast.Node) (w ast.Visitor) {
	if call, ok := node.(*ast.CallExpr); ok {
		if f, ok := call.Fun.(*ast.Ident); ok {
			isFound := false
			for i := range goFuncs {
				if f.Name == goFuncs[i] {
					isFound = true
				}
			}
			if !isFound {
				c.used = append(c.used, f.Name)
			}
		}
	}
	return c
}
