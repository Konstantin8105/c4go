package transpiler

import (
	"bytes"
	"fmt"
	"strings"

	goast "go/ast"
	"go/format"
	"go/token"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
)

func generateBinding(p *program.Program) (bindHeader, bindCode string) {
	// outside called functions
	ds := p.GetOutsideCalledFunctions()
	if len(ds) == 0 {
		return
	}

	// automatic binding of function
	{
		in := map[string]bool{}
		for i := range ds {
			in[p.PreprocessorFile.GetBaseInclude(ds[i].IncludeFile)] = true
		}
		for header := range in {
			bindHeader += fmt.Sprintf("// #include <%s>\n", header)
		}
		bindHeader += "import \"C\"\n\n"
	}

	for i := range ds {
		//
		// Example:
		//
		// // #include <stdlib.h>
		// // #include <stdio.h>
		// // #include <errno.h>
		// import "C"
		//
		// func Seed(i int) {
		//   C.srandom(C.uint(i))
		// }
		//

		// input data:
		// {frexp double [double int *] true true  [] []}
		//
		// output:
		// func  frexp(arg1 float64, arg2 []int) float64 {
		//		return float64(C.frexp(C.double(arg1), unsafe.Pointer(arg2)))
		// }

		mess := p.GenerateWarningMessage(fmt.Errorf(
			"Add c-binding for implementate function : `%s`", ds[i].Name), nil)
		bindCode += mess + "\n"

		code, err := getBindFunction(p, ds[i])
		if err != nil {
			bindCode += p.GenerateWarningMessage(err, nil) + "\n"
			continue
		}
		index := strings.Index(code, "\n")
		if index < 0 {
			continue
		}
		bindCode += code[index:] + "\n"
	}

	return
}

func getBindFunction(p *program.Program, d program.DefinitionFunction) (code string, err error) {
	var f goast.FuncDecl
	f.Name = goast.NewIdent(d.Name)

	prefix := "arg"
	// arguments types
	var ft goast.FuncType
	var fl goast.FieldList
	var argResolvedType []string
	for i := range d.ArgumentTypes {
		resolveType, err := types.ResolveType(p, d.ArgumentTypes[i])
		if err != nil {
			return "", fmt.Errorf("cannot generate argument binding function `%s`: %v", d.Name, err)
		}
		fl.List = append(fl.List, &goast.Field{
			Names: []*goast.Ident{goast.NewIdent(fmt.Sprintf("%s%d", prefix, i))},
			Type:  goast.NewIdent(resolveType),
		})
		argResolvedType = append(argResolvedType, resolveType)
	}
	ft.Params = &fl
	f.Type = &ft
	// return type
	var fr goast.FieldList
	ft.Results = &fr
	var returnResolvedType string
	if d.ReturnType != "" {
		resolveType, err := types.ResolveType(p, d.ReturnType)
		if err != nil {
			return "", fmt.Errorf("cannot generate return type binding function `%s`: %v", d.Name, err)
		}
		fr.List = append(fr.List, &goast.Field{
			Type: goast.NewIdent(resolveType),
		})
		returnResolvedType = resolveType
	}

	// create body
	var arg []goast.Expr
	for i := range argResolvedType {
		// convert from Go type to Cgo type
		cgoExpr, err := ResolveCgoType(p, argResolvedType[i], goast.NewIdent(fmt.Sprintf("%s%d", prefix, i)))
		if err != nil {
			return "", fmt.Errorf("cannot resolve cgo type for function `%s`: %v", d.Name, err)
		}

		arg = append(arg, cgoExpr)
	}

	f.Body = &goast.BlockStmt{
		List: []goast.Stmt{
			&goast.ReturnStmt{
				Results: []goast.Expr{
					&goast.CallExpr{
						Fun: goast.NewIdent(returnResolvedType),
						Args: []goast.Expr{
							&goast.CallExpr{
								Fun:  goast.NewIdent(fmt.Sprintf("C.%s", d.Name)),
								Args: arg,
							},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), &goast.File{
		Name:  goast.NewIdent("main"),
		Decls: []goast.Decl{&f},
	}); err != nil {
		return "", fmt.Errorf("cannot get source of binding function : %s", d.Name)
	}

	return buf.String(), nil
}

func cgoTypes(goType string) (_ string, ok bool) {
	goType = strings.TrimSpace(goType)
	switch goType {
	case "int":
		return "int", true
	case "float64":
		return "double", true
	}
	return "", false
}

func ResolveCgoType(p *program.Program, goType string, expr goast.Expr) (a goast.Expr, err error) {
	if ct, ok := cgoTypes(goType); ok {
		return &goast.CallExpr{
			Fun: &goast.SelectorExpr{
				X:   goast.NewIdent("C"),
				Sel: goast.NewIdent(ct),
			},
			Args: []goast.Expr{expr},
		}, nil
	}

	if strings.Contains(goType, "[") {
		// []int  -> * _Ctype_int
		p.AddImport("unsafe")
		t, ok := cgoTypes(goType[2:])
		if ok {
			t = "( * _Ctype_" + t + " ) "
			return &goast.CallExpr{
				Fun: goast.NewIdent(t),
				Args: []goast.Expr{
					&goast.CallExpr{
						Fun: goast.NewIdent("unsafe.Pointer"),
						Args: []goast.Expr{
							&goast.UnaryExpr{
								Op: token.AND,
								X: &goast.IndexExpr{
									X:      expr,
									Lbrack: 1,
									Index:  goast.NewIdent("0"),
								},
							},
						},
					},
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("cannot resolve to cgo type: `%s`", goType)
}
