package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/format"
	"go/token"
	"sort"
	"strings"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

func generateBinding(p *program.Program) (bindHeader, bindCode string) {
	// outside called functions
	ds := p.GetOutsideCalledFunctions()
	if len(ds) == 0 {
		return
	}

	sort.Slice(ds, func(i, j int) bool {
		return ds[i].Name < ds[j].Name
	})

	// automatic binding of function
	{
		in := map[string]bool{}
		for i := range ds {
			y := ds[i].IncludeFile
			in[y] = true
			in[p.PreprocessorFile.GetBaseInclude(y)] = true
		}
		for header := range in {
			if strings.Contains(header, "bits") {
				continue
			}
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
		if i == len(d.ArgumentTypes)-1 && d.ArgumentTypes[i] == "..." {
			argResolvedType[len(argResolvedType)-1] =
				"..." + argResolvedType[len(argResolvedType)-1]
			continue
		}
		resolveType, err := types.ResolveType(p, d.ArgumentTypes[i])
		if err != nil {
			return "", fmt.Errorf("cannot generate argument binding function `%s`: %v", d.Name, err)
		}
		argResolvedType = append(argResolvedType, resolveType)
	}
	for i := range argResolvedType {
		resolveType := argResolvedType[i]
		fl.List = append(fl.List, &goast.Field{
			Names: []*goast.Ident{goast.NewIdent(fmt.Sprintf("%s%d", prefix, i))},
			Type:  goast.NewIdent(resolveType),
		})
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

	f.Body = &goast.BlockStmt{}

	stmts := bindFromCtoGo(p, d.ReturnType, returnResolvedType, util.NewCallExpr(fmt.Sprintf("C.%s", d.Name), arg...))
	f.Body.List = append(f.Body.List, stmts...)

	// add comment for function
	f.Doc = &goast.CommentGroup{
		List: []*goast.Comment{
			&goast.Comment{
				Text: fmt.Sprintf("// %s - add c-binding for implemention function", d.Name),
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
	case "int32":
		return "int", true
	case "int64":
		return "long", true
	case "float64":
		return "double", true
	case "byte":
		return "char", true
	case "uint":
		return "ulong", true
	case "noarch.Tm":
		return "struct_tm", true
	case "noarch.File":
		return "FILE", true
	case "uint32":
		return "ulong", true
	}
	return "", false
}

// TODO : add implementation
//
// Example:
// func write(arg0 int32, arg1 interface{}, arg2 uint) noarch.SsizeT {
//      a := arg1.([]byte)
//      b := string(a)
//      c := C.CString(b)
//      return noarch.SsizeT(C.write(C.int(arg0), (unsafe.Pointer(c)), C.ulong(arg2)))
// }
//
// func read(arg0 int32, arg1 interface{}, arg2 uint) noarch.SsizeT {
//      a := arg1.([]byte)
//      b := string(a)
//      c := C.CString(b)
//      S := noarch.SsizeT(C.read(C.int(arg0), unsafe.Pointer(c), C.ulong(arg2)))
//      d := C.GoString(c)
//      arg1 = []byte(d)
//      return S
// }
//
//	func read(arg0 int32, arg1 interface{}, arg2 uint) noarch.SsizeT {
//	   switch v := arg1.(type) {
//	   case []byte:
//	   	a := v
//	   	b := string(a)
//	   	c := C.CString(b)
//	   	S := noarch.SsizeT(C.read(C.int(arg0), unsafe.Pointer(c), C.ulong(arg2)))
//	   	d := C.GoString(c)
//	   	arg1 = []byte(d)
//	   	return S
//	   case *[]byte:
//	   	a := v
//	   	b := string(*a)
//	   	c := C.CString(b)
//	   	S := noarch.SsizeT(C.read(C.int(arg0), unsafe.Pointer(c), C.ulong(arg2)))
//	   	d := C.GoString(c)
//	   	arg1 = []byte(d)
//	   	return S
//	   }
//	   return noarch.SsizeT(C.read(C.int(arg0), unsafe.Pointer(&arg1), C.ulong(arg2)))
//	}
//
// 	func write(arg0 int32, arg1 interface{}, arg2 uint) noarch.SsizeT {
// 	   switch v := arg1.(type) {
// 	   case []byte: // []uint8:
// 	   	a := v
// 	   	b := string(a)
// 	   	c := C.CString(b)
// 	   	return noarch.SsizeT(C.write(C.int(arg0), (unsafe.Pointer(c)), C.ulong(arg2)))
// 	   }
// 	   return noarch.SsizeT(C.write(C.int(arg0), (unsafe.Pointer(&arg1)), C.ulong(arg2)))
// 	}

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

	t := goType

	if strings.HasPrefix(goType, "[][]") {
		t = "interface{}"
	} else if strings.HasPrefix(goType, "[") {
		// []int  -> * _Ctype_int
		t = goType[2:]
		var ok bool
		t, ok = cgoTypes(t)
		if !ok {
			// TODO: check next
			t = goType[2:]
		}
		t = "( * C." + t + " ) "
		t = strings.Replace(t, " ", "", -1)

		p.AddImport("unsafe")

		return util.NewCallExpr(t, util.NewCallExpr("unsafe.Pointer",
			&goast.UnaryExpr{
				Op: token.AND,
				X: &goast.IndexExpr{
					X:      expr,
					Lbrack: 1,
					Index:  goast.NewIdent("0"),
				},
			})), nil

	} else if strings.HasPrefix(goType, "*") {
		// *int  -> * _Ctype_int
		t = goType[1:]
		var ok bool
		t, ok = cgoTypes(t)
		if !ok {
			// TODO: check next
			t = goType[1:]
		}
		t = "( * C." + t + " ) "
		t = strings.Replace(t, " ", "", -1)

		p.AddImport("unsafe")

		return util.NewCallExpr(t, util.NewCallExpr("unsafe.Pointer",
			&goast.UnaryExpr{
				Op: token.AND,
				X:  expr,
			})), nil
	}

	if t == "interface{}" {

		p.AddImport("unsafe")

		return util.NewCallExpr("unsafe.Pointer",
			&goast.UnaryExpr{
				Op: token.AND,
				X:  expr,
			}), nil
	}

	return &goast.CallExpr{
		Fun: &goast.SelectorExpr{
			X:   goast.NewIdent("C"),
			Sel: goast.NewIdent(t),
		},
		Args: []goast.Expr{expr},
	}, nil
}

// example:
//
// returnValue := ...
// return cast_from_C_to_Go_type(returnValue)
//
func bindFromCtoGo(p *program.Program, cType string, goType string, expr goast.Expr) (stmts []goast.Stmt) {

	if expr == nil {
		expr = goast.NewIdent("C4GO_UNDEFINE_EXPR")
	}
	if goType == "" {
		goType = "C4GO_UNDEFINE_GO_TYPE"
	}

	if cType == "" || cType == "void" {
		stmts = append(stmts, &goast.ReturnStmt{Results: []goast.Expr{expr}})
		return
	}

	// from documentation : https://golang.org/cmd/cgo/
	//
	// C string to Go string
	// func C.GoString(*C.char) string
	//

	switch cType {
	case "char *":
		stmts = append(stmts, &goast.ReturnStmt{Results: []goast.Expr{
			util.NewCallExpr("[]byte",
				util.NewCallExpr("C.GoString", expr),
			),
		}})

	default:
		stmts = append(stmts, &goast.ReturnStmt{Results: []goast.Expr{
			util.NewCallExpr(goType, expr),
		}})
	}

	return
}
