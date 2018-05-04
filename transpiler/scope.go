// This file contains functions for transpiling scopes. A scope is zero or more
// statements between a set of curly brackets.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
)

func transpileCompoundStmt(n *ast.CompoundStmt, p *program.Program) (
	_ *goast.BlockStmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	stmts := []goast.Stmt{}

	for _, x := range n.Children() {
		// Implementation va_start
		var isVaList bool
		if call, ok := x.(*ast.CallExpr); ok && call.Type == "void" {
			if impl, ok := call.Children()[0].(*ast.ImplicitCastExpr); ok {
				if impl.Type == "void (*)(struct __va_list_tag *, ...)" {
					if decl, ok := impl.Children()[0].(*ast.DeclRefExpr); ok {
						if decl.Name == "__builtin_va_start" {
							isVaList = true
						}
					}
				}
			}
		}

		var result []goast.Stmt
		if isVaList {
			// Implementation va_start
			result = []goast.Stmt{
				&goast.AssignStmt{
					Lhs: []goast.Expr{
						goast.NewIdent("c4goVaListPosition"),
					},
					Tok: token.ASSIGN,
					Rhs: []goast.Expr{
						&goast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					},
				},
			}
		} else {
			// Other cases
			if parent, ok := x.(*ast.ParenExpr); ok {
				x = parent.Children()[0]
			}
			result, err = transpileToStmts(x, p)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		if result != nil {
			stmts = append(stmts, result...)
		}
	}

	return &goast.BlockStmt{
		List: stmts,
	}, preStmts, postStmts, nil
}

func transpileToBlockStmt(node ast.Node, p *program.Program) (
	*goast.BlockStmt, []goast.Stmt, []goast.Stmt, error) {
	stmts, err := transpileToStmts(node, p)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(stmts) == 1 {
		if block, ok := stmts[0].(*goast.BlockStmt); ok {
			return block, nil, nil, nil
		}
	}

	if stmts == nil {
		return nil, nil, nil, fmt.Errorf("Stmts inside Block cannot be nil")
	}

	return &goast.BlockStmt{
		List: stmts,
	}, nil, nil, nil
}
