package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"sort"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

const unsafeConvertFunctionName string = "c4goUnsafeConvert_"

func ConvertValueToPointer(nodes []ast.Node, p *program.Program) (expr goast.Expr, ok bool) {
	if len(nodes) != 1 {
		return nil, false
	}

	decl, ok := nodes[0].(*ast.DeclRefExpr)
	if !ok {
		return nil, false
	}

	if types.IsPointer(decl.Type, p) {
		return nil, false
	}

	// get base type if it typedef
	var td string = decl.Type
	for {
		if t, ok := p.TypedefType[td]; ok {
			td = t
			continue
		}
		break
	}

	resolvedType, err := types.ResolveType(p, td)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, decl))
		return
	}

	var acceptable bool

	if types.IsGoBaseType(resolvedType) {
		acceptable = true
	}

	if str, ok := p.Structs[decl.Type]; ok && str.IsGlobal {
		acceptable = true
	}

	if str, ok := p.Unions[decl.Type]; ok && str.IsGlobal {
		acceptable = true
	}

	if !acceptable {
		return nil, false
	}

	// can simplify
	p.UnsafeConvertValueToPointer[resolvedType] = true

	return util.NewCallExpr(fmt.Sprintf("%s%s", unsafeConvertFunctionName, resolvedType),
		&goast.UnaryExpr{
			Op: token.AND,
			X:  goast.NewIdent(decl.Name),
		}), true
}

func GetUnsafeConvertDecls(p *program.Program) {
	if len(p.UnsafeConvertValueToPointer) == 0 {
		return
	}

	p.AddImport("unsafe")

	var names []string
	for t := range p.UnsafeConvertValueToPointer {
		names = append(names, t)
	}
	sort.Sort(sort.StringSlice(names))

	for _, t := range names {
		functionName := fmt.Sprintf("%s%s", unsafeConvertFunctionName, t)
		varName := "c4go_name"
		p.File.Decls = append(p.File.Decls, &goast.FuncDecl{
			Doc: &goast.CommentGroup{
				List: []*goast.Comment{
					&goast.Comment{
						Text: fmt.Sprintf("// %s : created by c4go\n", functionName),
					},
				},
			},
			Name: goast.NewIdent(functionName),
			Type: &goast.FuncType{
				Params: &goast.FieldList{
					List: []*goast.Field{
						{
							Names: []*goast.Ident{goast.NewIdent(varName)},
							Type:  goast.NewIdent("*" + t),
						},
					},
				},
				Results: &goast.FieldList{
					List: []*goast.Field{
						{
							Type: &goast.ArrayType{
								Lbrack: 1,
								Elt:    goast.NewIdent(t),
							},
						},
					},
				},
			},
			Body: &goast.BlockStmt{
				List: []goast.Stmt{
					&goast.ReturnStmt{
						Results: []goast.Expr{
							&goast.SliceExpr{
								X: util.NewCallExpr(fmt.Sprintf("(*[1000000]%s)", t),
									util.NewCallExpr("unsafe.Pointer",
										goast.NewIdent(varName)),
								),
							},
						},
					},
				},
			},
		})
	}

	return
}

// GetUintptrForSlice - return uintptr for slice
// Example : int64(uintptr(unsafe.Pointer((*(**int)(unsafe.Pointer(& ...slice... )))))))
func GetUintptrForSlice(expr goast.Expr, sizeof int) (goast.Expr, string) {
	returnType := "long long"

	if _, ok := expr.(*goast.SelectorExpr); ok {
		expr = &goast.IndexExpr{
			X:     expr,
			Index: goast.NewIdent("0"),
		}
	}

	if sl, ok := expr.(*goast.SliceExpr); ok {
		// from :
		//
		// 88  0: *ast.SliceExpr {
		// 89  .  X: *ast.Ident {
		// 91  .  .  Name: "b"
		// 93  .  }
		// 95  .  Low: *ast.BasicLit { ... }
		// 99  .  }
		// 102  }
		//
		// to:
		//
		// 0  *ast.IndexExpr {
		// 1  .  X: *ast.Ident {
		// 3  .  .  Name: "b"
		// 4  .  }
		// 6  .  Index: *ast.BasicLit { ... }
		// 12  }
		if sl.Low == nil {
			sl.Low = goast.NewIdent("0")
		}
		util.PanicIfNil(sl.X, "slice is nil")
		util.PanicIfNil(sl.Low, "slice low is nil")
		expr = &goast.IndexExpr{
			X:     sl.X,
			Index: sl.Low,
		}
	}

	if sl, ok := expr.(*goast.SliceExpr); ok {
		if c, ok := sl.X.(*goast.CallExpr); ok {
			if fin, ok := c.Fun.(*goast.Ident); ok && strings.Contains(fin.Name, "1000000") {
				if len(c.Args) == 1 {
					if cc, ok := c.Args[0].(*goast.CallExpr); ok {
						if fin, ok := cc.Fun.(*goast.Ident); ok && strings.Contains(fin.Name, "unsafe.Pointer") {
							if len(cc.Args) == 1 {
								if un, ok := cc.Args[0].(*goast.UnaryExpr); ok && un.Op == token.AND {
									expr = un.X
								}
							}
						}
					}
				}
			}
		}
	}

	if _, ok := expr.(*goast.CallExpr); ok {
		name := "c4go_temp_name"
		expr = util.NewAnonymousFunction(
			// body
			[]goast.Stmt{
				&goast.ExprStmt{
					X: &goast.BinaryExpr{
						X:  goast.NewIdent(name),
						Op: token.DEFINE,
						Y:  expr,
					},
				},
			},
			// defer
			nil,
			// returnValue
			util.NewCallExpr("int64", util.NewCallExpr("uintptr", util.NewCallExpr("unsafe.Pointer",
				&goast.StarExpr{
					Star: 1,
					X: &goast.CallExpr{
						Fun:    goast.NewIdent("(**byte)"),
						Lparen: 1,
						Args: []goast.Expr{&goast.CallExpr{
							Fun:    goast.NewIdent("unsafe.Pointer"),
							Lparen: 1,
							Args: []goast.Expr{
								&goast.UnaryExpr{
									Op: token.AND,
									X:  goast.NewIdent(name),
								},
							},
						}},
					},
				},
			))),
			// returnType
			"int64",
		)
		return expr, returnType
	}

	return &goast.BinaryExpr{
		X: util.NewCallExpr("int64", util.NewCallExpr("uintptr", util.NewCallExpr("unsafe.Pointer",
			&goast.StarExpr{
				Star: 1,
				X: &goast.CallExpr{
					Fun:    goast.NewIdent("(**byte)"),
					Lparen: 1,
					Args: []goast.Expr{&goast.CallExpr{
						Fun:    goast.NewIdent("unsafe.Pointer"),
						Lparen: 1,
						Args: []goast.Expr{&goast.UnaryExpr{
							Op: token.AND,
							X:  expr,
						}},
					}},
				},
			},
		))),
		Op: token.QUO,
		Y:  util.NewCallExpr("int64", goast.NewIdent(fmt.Sprintf("%d", sizeof))),
	}, returnType
}

// CreateSliceFromReference - create a slice, like :
// (*[1]int)(unsafe.Pointer(&a))[:]
func CreateSliceFromReference(goType string, expr goast.Expr) goast.Expr {
	// If the Go type is blank it means that the C type is 'void'.
	if goType == "" {
		goType = "interface{}"
	}

	// TODO remove
	fmt.Printf("%#v\n", expr)
	goast.Print(token.NewFileSet(), expr)

	// This is a hack to convert a reference to a variable into a slice that
	// points to the same location. It will look similar to:
	//
	//     (*[1]int)(unsafe.Pointer(&a))[:]
	//
	// You must always call this Go before using CreateSliceFromReference:
	//
	//     p.AddImport("unsafe")
	//
	return &goast.SliceExpr{
		X: util.NewCallExpr(fmt.Sprintf("(*[1000000]%s)", goType),
			util.NewCallExpr("unsafe.Pointer",
				&goast.UnaryExpr{
					X:  expr,
					Op: token.AND,
				}),
		),
	}
}

// pointerArithmetic - operations between 'int' and pointer
// Example C code : ptr += i
// ptr = (*(*[1]int)(unsafe.Pointer(uintptr(unsafe.Pointer(&ptr[0])) + (i)*unsafe.Sizeof(ptr[0]))))[:]
// , where i  - right
//        '+' - operator
//      'ptr' - left
//      'int' - leftType transpiled in Go type
// Note:
// 1) rigthType MUST be 'int'
// 2) pointerArithmetic - implemented ONLY right part of formula
// 3) right is MUST be positive value, because impossible multiply uintptr to (-1)
func pointerArithmetic(p *program.Program,
	left goast.Expr, leftType string,
	right goast.Expr, rightType string,
	operator token.Token) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile pointerArithmetic. err = %v", err)
		}
	}()

	if bl, ok := right.(*goast.BasicLit); ok && operator == token.ADD {
		if bl.Value == "1" && bl.Kind == token.INT {
			return &goast.SliceExpr{
				X:      left,
				Lbrack: 1,
				Low: &goast.BasicLit{
					Kind:  token.INT,
					Value: "1",
				},
				Slice3: false,
			}, leftType, preStmts, postStmts, nil
		}
	}

	if !(types.IsCInteger(p, rightType) || rightType == "bool") {
		err = fmt.Errorf("right type is not C integer type : '%s'", rightType)
		return
	}
	if !types.IsPointer(leftType, p) {
		err = fmt.Errorf("left type is not a pointer : '%s'", leftType)
		return
	}
	right, err = types.CastExpr(p, right, rightType, "int")
	if err != nil {
		return
	}

	for {
		if t, ok := p.TypedefType[leftType]; ok {
			leftType = t
			continue
		}
		break
	}

	resolvedLeftType, err := types.ResolveType(p, leftType)
	if err != nil {
		return
	}

	type pA struct {
		Name      string // name of variable: 'ptr'
		Type      string // type of variable: 'int','double'
		Condition string // condition : '-1' ,'(-1+2-2)'
		Operator  string // operator : '+', '-'
	}

	var s pA

	{
		var buf bytes.Buffer
		_ = printer.Fprint(&buf, token.NewFileSet(), left)
		s.Name = buf.String()
	}
	{
		var buf bytes.Buffer
		_ = printer.Fprint(&buf, token.NewFileSet(), right)
		s.Condition = buf.String()
	}

	s.Type = resolvedLeftType
	if resolvedLeftType != "interface{}" {
		s.Type = resolvedLeftType[2:]
	}

	s.Operator = "+"
	if operator == token.SUB {
		s.Operator = "-"
	}

	src := `package main
func main(){
	a := (*(*[1000000]{{ .Type }})(unsafe.Pointer(uintptr(
			unsafe.Pointer(&{{ .Name }}[0])) {{ .Operator }}
			(uintptr)({{ .Condition }})*unsafe.Sizeof({{ .Name }}[0]))))[:]
}`
	tmpl := template.Must(template.New("").Parse(src))
	var source bytes.Buffer
	err = tmpl.Execute(&source, s)
	if err != nil {
		err = fmt.Errorf("Cannot execute template. err = %v", err)
		return
	}

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	body := strings.Replace(source.String(), "&#43;", "+", -1)
	body = strings.Replace(body, "&amp;", "&", -1)
	body = strings.Replace(body, "&#34;", "\"", -1)
	body = strings.Replace(body, "&#39;", "'", -1)
	body = strings.Replace(body, "&gt;", ">", -1)
	body = strings.Replace(body, "&lt;", "<", -1)
	// TODO: add unicode convertor
	f, err := parser.ParseFile(fset, "", body, 0)
	if err != nil {
		body = strings.Replace(body, "\n", "", -1)
		err = fmt.Errorf("Cannot parse file. err = %v. body = `%s`", err, body)
		return
	}

	p.AddImport("unsafe")

	return f.Decls[0].(*goast.FuncDecl).Body.List[0].(*goast.AssignStmt).Rhs[0],
		leftType, preStmts, postStmts, nil
}
