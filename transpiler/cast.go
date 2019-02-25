package transpiler

import (
	"fmt"
	goast "go/ast"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
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
		}
	}()

	n.Type = util.GenerateCorrectType(n.Type)
	n.Type2 = util.GenerateCorrectType(n.Type2)

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
	if n.Kind == "IntegralToPointer" {
		// ImplicitCastExpr 'double *' <IntegralToPointer>
		// `-ImplicitCastExpr 'long' <LValueToRValue>
		//   `-DeclRefExpr 'long' lvalue Var 0x30e91d8 'pnt' 'long'
		if util.IsCPointer(n.Type) {
			if t, ok := ast.GetTypeIfExist(n.Children()[0]); ok {
				if types.IsCInteger(p, *t) {
					resolveType := n.Type
					resolveType, err = types.ResolveType(p, n.Type)
					if err != nil {
						return nil, "", nil, nil, err
					}
					expr = &goast.StarExpr{
						X: &goast.ParenExpr{
							X: &goast.CallExpr{
								Fun: &goast.ParenExpr{X: goast.NewIdent("*" + resolveType)},
								Args: []goast.Expr{
									&goast.CallExpr{
										Fun: goast.NewIdent("unsafe.Pointer"),
										Args: []goast.Expr{
											&goast.CallExpr{
												Fun:  goast.NewIdent("uintptr"),
												Args: []goast.Expr{expr},
											},
										},
									},
								},
							},
						},
					}
					p.GenerateWarningMessage(
						fmt.Errorf("used unsafe convert from integer to pointer"), n)
					exprType = n.Type
					return
				} else {
					//
					// ImplicitCastExpr 'char *' <IntegralToPointer>
					// `-ImplicitCastExpr 'char' <LValueToRValue>
					//   `-ArraySubscriptExpr 'char' lvalue
					//     |-ImplicitCastExpr 'char *' <LValueToRValue>
					//     | `-DeclRefExpr 'char *' lvalue Var 0x413c8a8 'b' 'char *'
					//     `-IntegerLiteral 'int' 3
					//
					// n.Type = 'char *'
					// *t     = 'char'
					//

					if ind, ok := expr.(*goast.IndexExpr); ok {
						// from :
						//
						// 0  *ast.IndexExpr {
						// 1  .  X: *ast.Ident {
						// 3  .  .  Name: "b"
						// 4  .  }
						// 6  .  Index: *ast.BasicLit { ... }
						// 12  }
						//
						// to:
						//
						// 88  0: *ast.SliceExpr {
						// 89  .  X: *ast.Ident {
						// 91  .  .  Name: "b"
						// 93  .  }
						// 95  .  Low: *ast.BasicLit { ... }
						// 99  .  }
						// 102  }
						expr = &goast.SliceExpr{
							X:   ind.X,
							Low: ind.Index,
						}
					}
				}
			}
		}
	}

	var cast bool = true
	if in, ok := n.Children()[0].(*ast.IntegerLiteral); ok && in.Type == "int" {
		if types.IsCInteger(p, n.Type) || types.IsCFloat(p, n.Type) {
			cast = false
			exprType = n.Type
		}
	}

	if len(n.Type) != 0 && len(n.Type2) != 0 && n.Type != n.Type2 && cast {
		var tt string
		tt, err = types.ResolveType(p, n.Type)
		expr = util.NewCallExpr(tt, expr)
		exprType = n.Type
		return
	}

	if util.IsFunction(exprType) {
		cast = false
	}
	if n.Kind == ast.ImplicitCastExprArrayToPointerDecay {
		cast = false
	}
	if n.Kind == "PointerToIntegral" {
		cast = false
	}

	if cast {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		exprType = n.Type
	}

	// Convert from struct member array to slice
	// ImplicitCastExpr 'char *' <ArrayToPointerDecay>
	// `-MemberExpr 'char [20]' lvalue .input_str 0x3662ba0
	//   `-DeclRefExpr 'struct s_inp':'struct s_inp' lvalue Var 0x3662c50 's' 'struct s_inp':'struct s_inp'
	if util.IsCPointer(n.Type) {
		if len(n.Children()) > 0 {
			if memb, ok := n.Children()[0].(*ast.MemberExpr); ok && util.IsCArray(memb.Type) {
				expr = &goast.SliceExpr{
					X:      expr,
					Lbrack: 1,
					Slice3: false,
				}
			}
		}
	}

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

	n.Type = util.GenerateCorrectType(n.Type)
	n.Type2 = util.GenerateCorrectType(n.Type2)

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
		expr = util.NewCallExpr(tt, expr)
		exprType = n.Type
		return
	}

	if n.Kind == ast.CStyleCastExprToVoid {
		exprType = types.ToVoid
		return
	}

	if !util.IsFunction(exprType) &&
		n.Kind != ast.ImplicitCastExprArrayToPointerDecay &&
		n.Kind != "PointerToIntegral" {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		exprType = n.Type
	}

	// CStyleCastExpr 'int' <PointerToIntegral>
	// `-UnaryOperator 'long *' prefix '&'
	//   `-DeclRefExpr 'long' lvalue Var 0x42b5268 'l' 'long'
	//
	// CStyleCastExpr 'int' <PointerToIntegral>
	// `-ParenExpr 'long *'
	//   `-UnaryOperator 'long *' prefix '&'
	//     `-DeclRefExpr 'long' lvalue Var 0x38cb568 'l' 'long'
	if len(n.Children()) > 0 {
		if types.IsCInteger(p, n.Type) {
			if t, ok := ast.GetTypeIfExist(n.Children()[0]); ok {
				if util.IsPointer(*t) {
					// main information	: https://go101.org/article/unsafe.html
					var retType string
					expr, retType = util.GetUintptrForSlice(expr)

					expr, err = types.CastExpr(p, expr, retType, n.Type)
					if err != nil {
						return nil, "", nil, nil, err
					}

					exprType = n.Type
				}
			}
		}
	}

	return
}
