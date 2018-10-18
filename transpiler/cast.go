package transpiler

import (
	"fmt"
	goast "go/ast"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
)

func transpileImplicitCastExpr(n *ast.ImplicitCastExpr, p *program.Program, exprIsStmt bool) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileImplicitCastExpr. err = %v", err)
			p.GenerateWarningMessage(err, n)
		}
	}()

	if n.Kind == ast.CStyleCastExprNullToPointer {
		expr = goast.NewIdent("nil")
		exprType = types.NullPointer
		return
	}

	expr, exprType, preStmts, postStmts, err = transpileToExpr(
		n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if exprType == types.NullPointer {
		expr = goast.NewIdent("nil")
		exprType = types.NullPointer
		return
	}

	n.Type = types.GenerateCorrectType(n.Type)
	exprType = types.GenerateCorrectType(exprType)

	expr, err = types.CastExpr(p, expr, exprType, n.Type)
	if err != nil {
		return nil, "", nil, nil, err
	}
	exprType = n.Type

	return
}

func transpileCStyleCastExpr(n *ast.CStyleCastExpr, p *program.Program, exprIsStmt bool) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileImplicitCastExpr. err = %v", err)
		}
	}()

	n.Type = types.GenerateCorrectType(n.Type)
	n.Type2 = types.GenerateCorrectType(n.Type2)

	// Char overflow
	// example for byte(-1)
	// CStyleCastExpr 0x365f628 <col:12, col:23> 'char' <IntegralCast>
	// `-ParenExpr 0x365f608 <col:18, col:23> 'int'
	//   `-ParenExpr 0x365f5a8 <col:19, col:22> 'int'
	//     `-UnaryOperator 0x365f588 <col:20, col:21> 'int' prefix '-'
	//       `-IntegerLiteral 0x365f568 <col:21> 'int' 1
	if n.Type == "char" {
		if par, ok := n.Children()[0].(*ast.ParenExpr); ok {
			if par2, ok := par.Children()[0].(*ast.ParenExpr); ok {
				if u, ok := par2.Children()[0].(*ast.UnaryOperator); ok && u.IsPrefix {
					if _, ok := u.Children()[0].(*ast.IntegerLiteral); ok {
						return transpileToExpr(&ast.BinaryOperator{
							Type:     "int",
							Type2:    "int",
							Operator: "+",
							ChildNodes: []ast.Node{
								u,
								&ast.IntegerLiteral{
									Type:  "int",
									Value: "256",
								},
							},
						}, p, false)
					}
				}
			}
		}
	}

	if n.Kind == ast.CStyleCastExprNullToPointer {
		expr = goast.NewIdent("nil")
		exprType = types.NullPointer
		return
	}

	expr, exprType, preStmts, postStmts, err = transpileToExpr(
		n.Children()[0], p, exprIsStmt)
	if err != nil {
		return nil, "", nil, nil, err
	}

	if exprType == types.NullPointer {
		expr = goast.NewIdent("nil")
		return
	}

	if len(n.Type) != 0 && len(n.Type2) != 0 && n.Type != n.Type2 {
		var tt string
		tt, err = types.ResolveType(p, n.Type)
		expr = &goast.CallExpr{
			Fun:    goast.NewIdent(tt),
			Lparen: 1,
			Args:   []goast.Expr{expr},
		}
		exprType = n.Type
		return
	}

	if n.Kind == ast.CStyleCastExprToVoid {
		exprType = types.ToVoid
		return
	}

	if !types.IsFunction(exprType) &&
		n.Kind != ast.ImplicitCastExprArrayToPointerDecay &&
		n.Kind != "PointerToIntegral" {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		exprType = n.Type
	}
	return
}
