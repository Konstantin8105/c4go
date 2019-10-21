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
		// add '_ = '
		var addPrefix bool
		if impl, ok := x.(*ast.ImplicitCastExpr); ok && len(impl.Children()) == 1 {
			if _, ok := impl.Children()[0].(*ast.ArraySubscriptExpr); ok {
				addPrefix = true
			}
		}

		var result []goast.Stmt
		// Other cases
		if parent, ok := x.(*ast.ParenExpr); ok {
			x = parent.Children()[0]
		}
		result, err = transpileToStmts(x, p)
		if err != nil {
			return nil, nil, nil, err
		}
		if addPrefix && len(result) == 1 {
			// goast.Print(token.NewFileSet(), result[0])
			result[0] = &goast.AssignStmt{
				Lhs: []goast.Expr{goast.NewIdent("_")},
				Tok: token.ASSIGN,
				Rhs: []goast.Expr{result[0].(*goast.ExprStmt).X},
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
		return nil, nil, nil, fmt.Errorf("stmts inside Block cannot be nil")
	}

	return &goast.BlockStmt{
		List: stmts,
	}, nil, nil, nil
}
