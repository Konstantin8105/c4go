// This file contains functions for declaring function prototypes, expressions
// that call functions, returning from function and the coordination of
// processing the function bodies.

package transpiler

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"

	goast "go/ast"
	"go/parser"
	"go/token"
)

// getFunctionBody returns the function body as a CompoundStmt. If the function
// is a prototype or forward declaration (meaning it has no body) then nil is
// returned.
func getFunctionBody(n *ast.FunctionDecl) *ast.CompoundStmt {
	// It's possible that the last node is the CompoundStmt (after all the
	// parameter declarations) - but I don't know this for certain so we will
	// look at all the children for now.
	for _, c := range n.Children() {
		if b, ok := c.(*ast.CompoundStmt); ok {
			return b
		}
	}

	return nil
}

// transpileFunctionDecl transpiles the function prototype.
//
// The function prototype may also have a body. If it does have a body the whole
// function will be transpiled into Go.
//
// If there is no function body we register the function internally (actually
// either way the function is registered internally) but we do not do anything
// because Go does not use or have any use for forward declarations of
// functions.
func transpileFunctionDecl(n *ast.FunctionDecl, p *program.Program) (
	decls []goast.Decl, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileFunctionDecl. %v", err)
		}
	}()

	// This is set at the start of the function declaration so when the
	// ReturnStmt comes alone it will know what the current function is, and
	// therefore be able to lookup what the real return type should be. I'm sure
	// there is a much better way of doing this.
	p.Function = n
	defer func() {
		// Reset the function name when we go out of scope.
		p.Function = nil
	}()

	n.Name = util.ConvertFunctionNameFromCtoGo(n.Name)

	// Always register the new function. Only from this point onwards will
	// we be allowed to refer to the function.
	define := func() (err error) {
		var pr string
		var f, r []string
		pr, _, f, r, err = util.ParseFunction(n.Type)
		if err != nil {
			err = fmt.Errorf("cannot get function definition : %v", err)
			return
		}
		if len(pr) != 0 {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("prefix of type '%s' is not empty", n.Type), n))
		}

		p.AddFunctionDefinition(program.DefinitionFunction{
			Name:          n.Name,
			ReturnType:    r[0],
			ArgumentTypes: f,
			Substitution:  "",
			IncludeFile:   n.Pos.File,
		})

		return
	}

	if p.GetFunctionDefinition(n.Name) == nil {
		if err = define(); err != nil {
			return
		}
	}

	if p.Binding {
		// probably a few function in result with same names
		decls, err = bindingFunctionDecl(n, p)
		return
	}

	if n.IsExtern {
		return
	}

	// Test if the function has a body. This is identified by a child node that
	// is a CompoundStmt (since it is not valid to have a function body without
	// curly brackets).
	functionBody := getFunctionBody(n)
	if functionBody == nil {
		return
	}

	if err = define(); err != nil {
		return
	}

	// If the function has a direct substitute in Go we do not want to
	// output the C definition of it.
	f := p.GetFunctionDefinition(n.Name)

	p.SetHaveBody(n.Name)
	body, pre, post, err := transpileToBlockStmt(functionBody, p)
	if err != nil || len(pre) > 0 || len(post) > 0 {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("not correct result in function %s body: err = %v",
				n.Name, err), n))
		err = nil // Error is ignored
	}

	if p.IncludeHeaderIsExists("stdlib.h") && n.Name == "main" {
		body.List = append([]goast.Stmt{&goast.DeferStmt{
			Call: util.NewCallExpr("noarch.AtexitRun"),
		}}, body.List...)
		p.AddImport("github.com/Konstantin8105/c4go/noarch")
	}

	// if functionBody != nil {

	// If verbose mode is on we print the name of the function as a comment
	// immediately to stdout. This will appear at the top of the program but
	// make it much easier to diagnose when the transpiler errors.
	if p.Verbose {
		fmt.Fprintf(os.Stdout, "// Function: %s(%s)\n", f.Name,
			strings.Join(f.ArgumentTypes, ", "))
	}

	var fieldList = &goast.FieldList{}
	fieldList, err = getFieldList(p, n, f.ArgumentTypes)
	if err != nil {
		return
	}

	// return type

	t, err := types.ResolveType(p, f.ReturnType)
	if err != nil {
		err = fmt.Errorf("ReturnType: %s. %v", f.ReturnType, err)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil
	}

	if p.Function != nil && p.Function.Name == "main" {
		// main() function does not have a return type.
		t = ""

		// This collects statements that will be placed at the top of
		// (before any other code) in main().
		prependStmtsInMain := []goast.Stmt{}

		// In Go, the main() function does not take the system arguments.
		// Instead they are accessed through the os package. We create new
		// variables in the main() function (if needed), immediately after
		// the __init() for these variables.
		if len(fieldList.List) > 0 {
			p.AddImport("os")

			prependStmtsInMain = append(
				prependStmtsInMain,
				&goast.AssignStmt{
					Lhs: []goast.Expr{fieldList.List[0].Names[0]},
					Tok: token.DEFINE,
					Rhs: []goast.Expr{
						&goast.CallExpr{
							Fun: goast.NewIdent("int32"),
							Args: []goast.Expr{
								util.NewCallExpr("len", goast.NewIdent("os.Args")),
							},
						},
					},
				},
			)
		}

		if len(fieldList.List) > 1 {
			prependStmtsInMain = append(
				prependStmtsInMain,
				&goast.AssignStmt{
					Lhs: []goast.Expr{fieldList.List[1].Names[0]},
					Tok: token.DEFINE,
					Rhs: []goast.Expr{&goast.CompositeLit{Type: util.NewTypeIdent("[][]byte")}},
				},
				&goast.RangeStmt{
					Key:   goast.NewIdent("_"),
					Value: util.NewIdent("argvSingle"),
					Tok:   token.DEFINE,
					X:     goast.NewIdent("os.Args"),
					Body: &goast.BlockStmt{
						List: []goast.Stmt{
							&goast.AssignStmt{
								Lhs: []goast.Expr{fieldList.List[1].Names[0]},
								Tok: token.ASSIGN,
								Rhs: []goast.Expr{util.NewCallExpr(
									"append",
									fieldList.List[1].Names[0],
									util.NewCallExpr("[]byte", util.NewIdent("argvSingle")),
								)},
							},
						},
					},
				})
		}

		// Prepend statements for main().
		body.List = append(prependStmtsInMain, body.List...)

		// The main() function does not have arguments or a return value.
		fieldList = &goast.FieldList{}
	}

	// Each function MUST have "ReturnStmt",
	// except function without return type
	var addReturnName bool
	if len(body.List) > 0 {
		last := body.List[len(body.List)-1]
		if _, ok := last.(*goast.ReturnStmt); !ok && t != "" {
			body.List = append(body.List, &goast.ReturnStmt{})
			addReturnName = true
		}
	}

	// For functions without return type - no need add return at
	// the end of body
	if p.GetFunctionDefinition(n.Name).ReturnType == "void" {
		if len(body.List) > 0 {
			if _, ok := (body.List[len(body.List)-1]).(*goast.ReturnStmt); ok {
				body.List = body.List[:len(body.List)-1]
			}
		}
	}

	decls = append(decls, &goast.FuncDecl{
		Name: util.NewIdent(n.Name),
		Type: util.NewFuncType(fieldList, t, addReturnName),
		Body: body,
	})
	//}

	err = nil
	return
}

// convert from:
// |-FunctionDecl 0x55677bd48ee0 <line:924:7, col:63> col:12 InitWindow 'void (int, int, const char *)'
// | |-ParmVarDecl 0x55677bd48d08 <col:23, col:27> col:27 width 'int'
// | |-ParmVarDecl 0x55677bd48d80 <col:34, col:38> col:38 height 'int'
// | `-ParmVarDecl 0x55677bd48df8 <col:46, col:58> col:58 title 'const char *'
//
// Go equal code:
//
//	// Initwindow is binding of function "InitWindow"
//	func InitWindow(width int32, height int32, title string) {
//		cwidth  := (C.int)(width)
//		cheight := (C.int)(height)
//		ctitle  := C.CString(title)
//		defer C.free(unsafe.Pointer(ctitle))
//		// run function
//		C.InitWindow(cwidth, cheight, ctitle)
//	}
func bindingFunctionDecl(n *ast.FunctionDecl, p *program.Program) (
	decls []goast.Decl, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot bindingFunctionDecl func `%s`. %v", n.Name, err)
		}
	}()

	// generate names
	n.Name = strings.TrimSpace(n.Name)
	cName := n.Name
	goName := func() string {
		rs := []rune(n.Name)
		name := strings.ToUpper(string(rs[0]))
		if 1 < len(rs) {
			name += string(rs[1:])
		}
		return name
	}()
	// parse variables
	type variable struct {
		cName, cType   string
		valid          bool
		cgoType        string
		goName, goType string
	}
	var args []variable
	for i := range n.ChildNodes {
		switch n := n.ChildNodes[i].(type) {
		case *ast.ParmVarDecl:
			var cgoType, goType string
			ok, cgoType, goType := cTypeToGoType(n.Type)
			args = append(args, variable{
				goName:  n.Name,
				cType:   n.Type,
				valid:   ok,
				cgoType: cgoType,
				goType:  goType,
			})
		default:
			err = fmt.Errorf("not valid type:%T", n)
			return
		}
	}
	for i := range args {
		args[i].cName = "c" + args[i].goName
	}

	f := p.GetFunctionDefinition(n.Name)

	var fieldList = &goast.FieldList{}
	fieldList, err = getFieldList(p, n, f.ArgumentTypes)
	if err != nil {
		return
	}

	// fix only for binding
	for i := range fieldList.List {
		name := fieldList.List[i].Type
		id, ok := name.(*goast.Ident)
		if !ok {
			continue
		}
		if id.Name == "[]byte" || id.Name == "[] byte" {
			fieldList.List[i].Type = goast.NewIdent("string")
		}
	}

	t, err := types.ResolveType(p, f.ReturnType)
	if err != nil {
		err = fmt.Errorf("ReturnType: %s. %v", f.ReturnType, err)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil
	}

	var fd goast.FuncDecl
	fd.Name = goast.NewIdent(goName)

	body := new(goast.BlockStmt)

	for i := range args {
		if cType := args[i].cType; cType == "char *" {
			// for type `char *`
			//
			// 	bs := []byte(*text)
			// 	if len(bs) == 0 {
			// 		bs = []byte{byte(0)}
			// 	}
			// 	if 0 < len(bs) && bs[len(bs)-1] != byte(0) { // minimalize allocation
			// 		bs = append(bs, byte(0)) // for next input symbols
			// 	}
			// 	ctext := (*C.char)(unsafe.Pointer(&bs[0]))
			// 	defer func() {
			// 		*text = string(bs)
			// 		// no need : C.free(unsafe.Pointer(ctext))
			// 	}()
			src := `package main
func main() {
	{{ .GoName }}_bs := []byte(* {{ .GoName }})
	if len({{ .GoName }}_bs) == 0 {
		{{ .GoName }}_bs = []byte{byte(0)}
	}
	if 0 < len({{ .GoName }}_bs) && {{ .GoName }}_bs[len({{ .GoName }}_bs)-1] != byte(0) {
		{{ .GoName }}_bs = append({{ .GoName }}_bs, byte(0))
	}
	{{ .CName }} := (*C.char)(unsafe.Pointer(&{{ .GoName }}_bs[0]))
	defer func() {
		*{{ .GoName }} = string({{ .GoName }}_bs)
		// no need : C.free(unsafe.Pointer(ctext))
	}()
}`
			tmpl := template.Must(template.New("").Parse(src))
			var source bytes.Buffer
			err = tmpl.Execute(&source, struct{ CName, GoName string }{
				CName: args[i].cName, GoName: args[i].goName,
			})
			if err != nil {
				err = fmt.Errorf("cannot execute template \"%s\" for data : %v",
					source.String(), err)
				return
			}

			// Create the AST by parsing src.
			fset := token.NewFileSet() // positions are relative to fset
			f, err := parser.ParseFile(fset, "", source.String(), 0)
			if err != nil {
				err = fmt.Errorf("cannot parse source \"%s\" : %v",
					source.String(), err)
				p.AddMessage(p.GenerateWarningMessage(err, n))
				err = nil // ignore error
				continue
			}
			if 0 < len(f.Decls) {
				if fd, ok := f.Decls[0].(*goast.FuncDecl); ok {
					body.List = append(body.List, fd.Body.List...)
				}
			}
			continue
		}
		if !args[i].valid {
			//	func Rect(r Rectangle, s int) {
			//		var cr C.struct_Rectangle
			//		cr.x = C.int(r.x)
			//		cr.y = C.int(r.y)
			//		cs = C.int(s)
			//		C.Rect(cr, cs)
			//	}
			st := p.GetStruct(args[i].cType)
			if st == nil {
				tname := p.TypedefType[args[i].cType]
				st = p.GetStruct(tname)
				if st != nil {
					args[i].cType = tname
				}
			}
			if st != nil &&
				!strings.Contains(args[i].cType, "*") &&
				!strings.Contains(args[i].cType, "[") {
				body.List = append(body.List, &goast.DeclStmt{Decl: &goast.GenDecl{
					Tok: token.VAR,
					Specs: []goast.Spec{&goast.ValueSpec{
						Names: []*goast.Ident{goast.NewIdent(args[i].cName)},
						Type: &goast.SelectorExpr{
							X:   goast.NewIdent("C"),
							Sel: goast.NewIdent("struct_" + args[i].cType),
						},
					}},
				}})
				for fname, ftype := range st.Fields {
					ft := fmt.Sprintf("%v", ftype)
					if strings.Contains(ft, "*") ||
						strings.Contains(ft, "[") {
						err = fmt.Errorf("field type is pointer: `%s`", ft)
						p.AddMessage(p.GenerateWarningMessage(err, n))
						err = nil // ignore error
					}
					_, cgot, _ := cTypeToGoType(ft)
					body.List = append(body.List, &goast.AssignStmt{
						Lhs: []goast.Expr{&goast.SelectorExpr{
							X:   goast.NewIdent(args[i].cName),
							Sel: goast.NewIdent(fname),
						}},
						Tok: token.ASSIGN,
						Rhs: []goast.Expr{&goast.CallExpr{
							Fun: goast.NewIdent(cgot),
							Args: []goast.Expr{&goast.SelectorExpr{
								X:   goast.NewIdent(args[i].goName),
								Sel: goast.NewIdent(fname),
							}},
						}},
					})
				}
				continue
			}
			if cType := args[i].cType; 2 < len(cType) &&
				cType[len(cType)-1] == '*' &&
				strings.Count(cType, "*") == 1 &&
				strings.Count(cType, "[") == 0 {
				cType = cType[:len(cType)-1]
				if ok, cgoType, goType := cTypeToGoType(cType); ok {
					//	if active == nil {
					//		active = new(int32)
					//	}
					var ifs goast.IfStmt
					ifs.Cond = &goast.BinaryExpr{
						X:  goast.NewIdent(args[i].goName),
						Op: token.EQL,
						Y:  goast.NewIdent("nil"),
					}
					ifs.Body = &goast.BlockStmt{
						List: []goast.Stmt{
							&goast.AssignStmt{
								Lhs: []goast.Expr{
									goast.NewIdent(args[i].goName),
								},
								Tok: token.ASSIGN,
								Rhs: []goast.Expr{
									&goast.CallExpr{
										Fun: goast.NewIdent("new"),
										Args: []goast.Expr{
											goast.NewIdent(goType),
										},
									},
								},
							},
						},
					}
					body.List = append(body.List, &ifs)
					//	cactive := C.int(*active)
					body.List = append(body.List, &goast.AssignStmt{
						Lhs: []goast.Expr{
							goast.NewIdent(args[i].cName),
						},
						Tok: token.DEFINE,
						Rhs: []goast.Expr{
							&goast.CallExpr{
								Fun: goast.NewIdent(cgoType),
								Args: []goast.Expr{
									&goast.StarExpr{
										X: goast.NewIdent(args[i].goName),
									},
								},
							},
						},
					})
					//	defer func() {
					//		*active = int32(cactive)
					//	}()
					body.List = append(body.List, &goast.DeferStmt{
						Call: &goast.CallExpr{
							Fun: &goast.FuncLit{
								Type: &goast.FuncType{},
								Body: &goast.BlockStmt{
									List: []goast.Stmt{
										&goast.AssignStmt{
											Lhs: []goast.Expr{
												&goast.StarExpr{
													X: goast.NewIdent(args[i].goName),
												},
											},
											Tok: token.ASSIGN,
											Rhs: []goast.Expr{
												goast.NewIdent(args[i].cName),
											},
										},
									},
								},
							},
						},
					})
					continue
				}
			}
			// fmt.Frintln(os.Stdout, "Not implemented: ", args[i].cType, args[i].cName)
			err = fmt.Errorf("cannot parse C type: `%s` and name `%s`",
				args[i].cType, args[i].cName)
			p.AddMessage(p.GenerateWarningMessage(err, n))
			err = nil // ignore error
			continue
		}
		body.List = append(body.List, &goast.AssignStmt{
			Lhs: []goast.Expr{goast.NewIdent(args[i].cName)},
			Tok: token.DEFINE,
			Rhs: []goast.Expr{&goast.CallExpr{
				Fun:  goast.NewIdent(args[i].cgoType),
				Args: []goast.Expr{goast.NewIdent(args[i].goName)}},
			},
		})
		// free memory
		if strings.Contains(args[i].cType, "*") ||
			strings.Contains(args[i].cType, "[") {
			body.List = append(body.List, &goast.DeferStmt{
				Call: &goast.CallExpr{
					Fun: &goast.SelectorExpr{
						X:   goast.NewIdent("C"),
						Sel: goast.NewIdent("free"),
					},
					Args: []goast.Expr{
						goast.NewIdent(fmt.Sprintf("unsafe.Pointer(%s)", args[i].cName)),
					},
				},
			})
		}
	}

	ce := &goast.CallExpr{
		Fun: &goast.SelectorExpr{
			X:   goast.NewIdent("C"),
			Sel: goast.NewIdent(cName),
		},
	}
	for i := range args {
		ce.Args = append(ce.Args, goast.NewIdent(args[i].cName))
	}

	runC := false
	if t == "" {
		body.List = append(body.List, &goast.ExprStmt{X: ce})
		runC = true
	}
	st := p.GetStruct(t)
	if st == nil {
		t2 := p.TypedefType[t]
		st = p.GetStruct(t2)
		if st != nil {
			t = t2
		}
	}
	if !runC && st != nil {
		//	func Rect() Rectangle {
		//		cResult := C.Rect()
		//		var goRes Rectangle
		//		goRes.x = int32(cResult.x)
		//		goRes.y = int32(cResult.y)
		//		return goRes
		//	}
		cResult := "cResult"
		body.List = append(body.List, &goast.AssignStmt{
			Lhs: []goast.Expr{goast.NewIdent(cResult)},
			Tok: token.DEFINE,
			Rhs: []goast.Expr{ce},
		})
		goRes := "goRes"
		body.List = append(body.List, &goast.DeclStmt{Decl: &goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{&goast.ValueSpec{
				Names: []*goast.Ident{goast.NewIdent(goRes)},
				Type:  goast.NewIdent(t),
			}},
		}})
		for fname, ftype := range st.Fields {
			ft := fmt.Sprintf("%v", ftype)
			if strings.Contains(ft, "*") ||
				strings.Contains(ft, "[") {
				err = fmt.Errorf("field type is pointer: `%s`", ft)
				p.AddMessage(p.GenerateWarningMessage(err, n))
				err = nil // ignore error
				return
			}
			_, _, goType := cTypeToGoType(fmt.Sprintf("%v", ft))
			body.List = append(body.List, &goast.AssignStmt{
				Lhs: []goast.Expr{&goast.SelectorExpr{
					X:   goast.NewIdent(goRes),
					Sel: goast.NewIdent(fname),
				}},
				Tok: token.ASSIGN,
				Rhs: []goast.Expr{&goast.CallExpr{
					Fun: goast.NewIdent(goType),
					Args: []goast.Expr{&goast.SelectorExpr{
						X:   goast.NewIdent(cResult),
						Sel: goast.NewIdent(fname),
					}},
				}},
			})
		}
		body.List = append(body.List, &goast.ReturnStmt{
			Results: []goast.Expr{goast.NewIdent(goRes)}})
		runC = true
	}
	if !runC {
		body.List = append(body.List, &goast.ReturnStmt{
			Results: []goast.Expr{ce}})
	}

	addReturnName := false

	decls = append(decls, &goast.FuncDecl{
		Name: util.NewIdent(goName),
		Type: util.NewFuncType(fieldList, t, addReturnName),
		Body: body,
	})

	return
}

// C/C++ type      CGO type      Go type
// C.char, C.schar (signed char), C.uchar (unsigned char), C.short, C.ushort (unsigned short), C.int, C.uint (unsigned int), C.long, C.ulong (unsigned long), C.longlong (long long), C.ulonglong (unsigned long long), C.float, C.double, C.complexfloat (complex float), and C.complexdouble (complex double)
// {"bool", "C.int32", "bool"},
var table = [][3]string{
	{"char", "C.char", "byte"},
	{"signed char", "C.schar", "int8"},
	{"unsigned char", "C.uchar", "byte"},
	{"short", "C.short", "int16"},
	{"unsigned short", "C.ushort", "uint16"},
	{"int", "C.int", "int"},
	{"unsigned int", "C.uint", "uint"},
	{"long", "C.long", "int64"},
	{"unsigned long", "C.ulong", "uint64"},
	{"long long", "C.longlong", "int64"},
	{"unsigned long long", "C.ulonglong", "uint64"},
	{"float", "C.float", "float32"},
	{"double", "C.double", "float64"},
	{"const char *", "C.CString", "string"},
	// {"char *", "C.CString", "string"},
	// {"char []", "C.CString", "string"},
	{"_Bool", "C.bool", "bool"},
}

func cTypeToGoType(cType string) (ok bool, cgoType, goType string) {
	cType = strings.TrimSpace(cType)
	for i := range table {
		if cType == table[i][0] {
			return true, table[i][1], table[i][2]
		}
	}
	return false, cType, cType
}

// getFieldList returns the parameters of a C function as a Go AST FieldList.
func getFieldList(p *program.Program, f *ast.FunctionDecl, fieldTypes []string) (
	_ *goast.FieldList, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in function field list. err = %v", err)
		}
	}()
	r := []*goast.Field{}
	for i := range fieldTypes {
		if len(f.Children()) <= i {
			err = fmt.Errorf("not correct type/children: %d, %d",
				len(f.Children()), len(fieldTypes))
			return
		}
		n := f.Children()[i]
		if v, ok := n.(*ast.ParmVarDecl); ok {
			t, err := types.ResolveType(p, fieldTypes[i])
			if err != nil {
				err = fmt.Errorf("FieldList type: %s. %v", fieldTypes[i], err)
				p.AddMessage(p.GenerateWarningMessage(err, f))
				err = nil // ignore error
				t = "C4GO_UNDEFINE_TYPE"
			}

			if t == "" {
				continue
			}
			r = append(r, &goast.Field{
				Names: []*goast.Ident{util.NewIdent(v.Name)},
				Type:  goast.NewIdent(t),
			})
		}
	}

	// for function argument: ...
	if strings.Contains(f.Type, "...") {
		r = append(r, &goast.Field{
			Names: []*goast.Ident{util.NewIdent("c4goArgs")},
			Type: &goast.Ellipsis{
				Ellipsis: 1,
				Elt: &goast.InterfaceType{
					Interface: 1,
					Methods: &goast.FieldList{
						Opening: 1,
					},
					Incomplete: false,
				},
			},
		})
	}

	return &goast.FieldList{
		List: r,
	}, nil
}

func transpileReturnStmt(n *ast.ReturnStmt, p *program.Program) (
	_ goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileReturnStmt. err = %v", err)
		}
	}()
	// There may not be a return value. Then we don't have to both ourselves
	// with all the rest of the logic below.
	if len(n.Children()) == 0 {
		return &goast.ReturnStmt{}, nil, nil, nil
	}

	var eType string
	var e goast.Expr
	e, eType, preStmts, postStmts, err = atomicOperation(n.Children()[0], p)
	if err != nil {
		return nil, nil, nil, err
	}
	if e == nil {
		return nil, nil, nil, fmt.Errorf("expr is nil")
	}

	f := p.GetFunctionDefinition(p.Function.Name)

	t, err := types.CastExpr(p, e, eType, f.ReturnType)
	if p.AddMessage(p.GenerateWarningMessage(err, n)) {
		t = util.NewNil()
	}

	results := []goast.Expr{t}

	// main() function is not allowed to return a result. Use os.Exit if
	// non-zero.
	if p.Function != nil && p.Function.Name == "main" {
		litExpr, isLiteral := e.(*goast.BasicLit)
		if !isLiteral || (isLiteral && litExpr.Value != "0") {
			p.AddImport("github.com/Konstantin8105/c4go/noarch")
			return util.NewExprStmt(&goast.CallExpr{
				Fun: goast.NewIdent("noarch.Exit"),
				Args: []goast.Expr{
					&goast.CallExpr{
						Fun:  goast.NewIdent("int32"),
						Args: results,
					},
				},
			}), preStmts, postStmts, nil
		}
		results = []goast.Expr{}
	}

	return &goast.ReturnStmt{
		Results: results,
	}, preStmts, postStmts, nil
}
