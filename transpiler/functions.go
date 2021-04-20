// This file contains functions for declaring function prototypes, expressions
// that call functions, returning from function and the coordination of
// processing the function bodies.

package transpiler

import (
	"fmt"
	"os"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"

	goast "go/ast"
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
	if p.GetFunctionDefinition(n.Name) == nil {
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
	}


	if n.IsExtern {
		return
	}

	// 	var haveCompound bool
	// 	for _, ch := range n.Children() {
	// 		if _, ok := ch.(*ast.CompoundStmt); ok {
	// 			haveCompound = true
	// 			break
	// 		}
	// 	}
	// 	if !haveCompound {
	// 		return
	// 	}
	//
	// 	if len(n.Children()) == 0 {
	// 		return
	// 	}

	var body *goast.BlockStmt

	// If the function has a direct substitute in Go we do not want to
	// output the C definition of it.
	f := p.GetFunctionDefinition(n.Name)

	// Test if the function has a body. This is identified by a child node that
	// is a CompoundStmt (since it is not valid to have a function body without
	// curly brackets).
	functionBody := getFunctionBody(n)
	if functionBody == nil {
		return
	}
	p.SetHaveBody(n.Name)
	var pre, post []goast.Stmt
	body, pre, post, err = transpileToBlockStmt(functionBody, p)
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
	fieldList, err = getFieldList(p, n,
		p.GetFunctionDefinition(n.Name).ArgumentTypes)
	if err != nil {
		return
	}

	t, err := types.ResolveType(p, f.ReturnType)
	p.AddMessage(p.GenerateWarningMessage(err, n))

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
		n := f.Children()[i]
		if v, ok := n.(*ast.ParmVarDecl); ok {
			t, err := types.ResolveType(p, fieldTypes[i])
			p.AddMessage(p.GenerateWarningMessage(err, f))

			if len(t) > 0 {
				r = append(r, &goast.Field{
					Names: []*goast.Ident{util.NewIdent(v.Name)},
					Type:  goast.NewIdent(t),
				})
			}
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
