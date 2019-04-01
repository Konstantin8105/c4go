package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

func transpileVAArgExpr(n *ast.VAArgExpr, p *program.Program) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileVAArgExpr. %v", err)
		}
	}()
	// -VAArgExpr 'int'
	//  `-ImplicitCastExpr 'struct __va_list_tag *' <ArrayToPointerDecay>
	//    `-DeclRefExpr 'va_list':'...' lvalue Var 'ap' 'va_list':'...'

	p.IsHaveVaList = true

	expr, exprType, preStmts, postStmts, err = atomicOperation(n.Children()[0], p)
	if err != nil {
		return expr, exprType, preStmts, postStmts, err
	}

	goType, err := types.ResolveType(p, n.Type)
	if err != nil {
		return expr, exprType, preStmts, postStmts, err
	}

	expr = &goast.TypeAssertExpr{
		X:      util.NewCallExpr(va_arg, expr),
		Lparen: 1,
		Type:   goast.NewIdent(goType),
	}
	exprType = n.Type
	return
}

func getVaListStruct() string {
	return `

// va_list is C4GO implementation of va_list from "stdarg.h"
type va_list struct{
	position int
	slice    []interface{}
}

func create_va_list(list []interface{}) *va_list{
	return &va_list{
		position: 0,
		slice   : list,
	}
}

func va_start(v * va_list, count interface{}) {
	v.position = 0
}

func va_end(v * va_list) {
	// do nothing
}

func va_arg(v * va_list) interface{} {
	defer func(){
		 v.position++	
	}()
	value := v.slice[v.position]
	switch value.(type) {
		case int: 
			return int32(value.(int))
		default:
			return value
	}
}

`
}

const (
	create_va_list string = "create_va_list"
	va_arg                = "va_arg"
	va_start              = "va_start"
	va_end                = "va_end"
)

func VaListInit(p *program.Program, name string) []goast.Decl {
	// variable for va_list. see "variadic function"
	// header : <stdarg.h>
	// Example :
	// DeclStmt 0x2fd87e0 <line:442:2, col:14>
	// `-VarDecl 0x2fd8780 <col:2, col:10> col:10 used args 'va_list':'struct __va_list_tag [1]'
	// Result:
	// ... - convert to - c4goArgs ...interface{}
	// var args = c4goArgs

	p.IsHaveVaList = true

	return []goast.Decl{&goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names: []*goast.Ident{util.NewIdent(name)},
				Values: []goast.Expr{
					util.NewCallExpr(create_va_list, util.NewIdent("c4goArgs")),
				},
			},
		},
	}}
}

func changeVaListFuncs(functionName *string) {
	switch *functionName {
	case "__builtin_va_start":
		mod := va_start
		*functionName = mod
	case "__builtin_va_end":
		mod := va_end
		*functionName = mod
	}
	return
}
