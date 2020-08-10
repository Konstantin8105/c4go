// This file contains functions transpiling some general operator expressions.
// See binary.go and unary.go.

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

// ternary without middle operation
//
// Example:
//
// BinaryConditionalOperator  'int'
// |-BinaryOperator 'int' '>'
// | |-IntegerLiteral 'int' 19
// | `-UnaryOperator 'int' prefix '-'
// |   `-IntegerLiteral 'int' 9
// |-OpaqueValueExpr 'int'
// | `-BinaryOperator 'int' '>'
// |   |-IntegerLiteral 'int' 19
// |   `-UnaryOperator 'int' prefix '-'
// |     `-IntegerLiteral 'int' 9
// |-OpaqueValueExpr  'int'
// | `-BinaryOperator  'int' '>'
// |   |-IntegerLiteral 'int' 19
// |   `-UnaryOperator 'int' prefix '-'
// |     `-IntegerLiteral 'int' 9
// `-IntegerLiteral 0x3646f70 <col:18> 'int' 23
func transpileBinaryConditionalOperator(n *ast.BinaryConditionalOperator, p *program.Program) (
	_ *goast.CallExpr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpile BinaryConditionalOperator : err = %v", err)
		}
	}()

	var co ast.ConditionalOperator
	co.Type = n.Type
	co.AddChild(n.Children()[0])
	co.AddChild(&ast.IntegerLiteral{
		Type:  co.Type,
		Value: "1",
	})
	co.AddChild(n.Children()[3])

	return transpileConditionalOperator(&co, p)
}

// transpileConditionalOperator transpiles a conditional (also known as a
// ternary) operator:
//
//     a ? b : c
//
// We cannot simply convert these to an "if" statement because they by inside
// another expression.
//
// Since Go does not support the ternary operator or inline "if" statements we
// use a closure to work the same way.
//
// It is also important to note that C only evaluates the "b" or "c" condition
// based on the result of "a" (from the above example).
//
// Example AST:
// ConditionalOperator 'int'
// |-ImplicitCastExpr 'int (*)(int)' <LValueToRValue>
// | `-DeclRefExpr 'int (*)(int)' lvalue Var 'v' 'int (*)(int)'
// |-IntegerLiteral 'int' 1
// `-CallExpr 'int'
//   |-...
//
// ConditionalOperator 'int'
// |-BinaryOperator 'int' '!='
// | |-...
// |-BinaryOperator 'int' '-'
// | |-...
// `-BinaryOperator 'int' '-'
//   |-...
func transpileConditionalOperator(n *ast.ConditionalOperator, p *program.Program) (
	_ *goast.CallExpr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpile ConditionalOperator : err = %v", err)
		}
	}()

	// a - condition
	a, aType, newPre, newPost, err := atomicOperation(n.Children()[0], p)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// null in C is zero
	if aType == types.NullPointer {
		a = &goast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		}
		aType = "int"
	}

	a, err = types.CastExpr(p, a, aType, "bool")
	if err != nil {
		err = fmt.Errorf("parameter `a` : %v", err)
		return
	}

	// b - body
	b, bType, newPre, newPost, err := atomicOperation(n.Children()[1], p)
	if err != nil {
		err = fmt.Errorf("parameter `b` : %v", err)
		return
	}
	// Theorephly, length is must be zero
	if len(newPre) > 0 || len(newPost) > 0 {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("length of pre or post in body must be zero. {%d,%d}", len(newPre), len(newPost)), n))
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if n.Type != "void" {
		b, err = types.CastExpr(p, b, bType, n.Type)
		if err != nil {
			return
		}
		bType = n.Type
	}

	// c - else body
	c, cType, newPre, newPost, err := atomicOperation(n.Children()[2], p)
	if err != nil {
		err = fmt.Errorf("parameter `c` : %v", err)
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if n.Type != "void" {
		c, err = types.CastExpr(p, c, cType, n.Type)
		if err != nil {
			err = fmt.Errorf("parameter `c` : %v", err)
			return
		}
		cType = n.Type
	}

	// rightType - generate return type
	var returnType string
	if n.Type != "void" {
		returnType, err = types.ResolveType(p, n.Type)
		if err != nil {
			return
		}
	}

	var bod, els goast.BlockStmt

	bod.Lbrace = 1
	if bType != types.ToVoid {
		if n.Type != "void" {
			bod.List = []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{b},
				},
			}
		} else {
			bod.List = []goast.Stmt{
				&goast.ExprStmt{
					X: b,
				},
			}
		}
	}

	els.Lbrace = 1
	if cType != types.ToVoid {
		if n.Type != "void" {
			els.List = []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{c},
				},
			}
		} else {
			els.List = []goast.Stmt{
				&goast.ExprStmt{
					X: c,
				},
			}
		}
	}

	stmts := append([]goast.Stmt{}, &goast.IfStmt{
		Cond: a,
		Body: &bod,
		Else: &els,
	})
	if len(bod.List) > 0 {
		if _, ok := bod.List[len(bod.List)-1].(*goast.ReturnStmt); ok {
			stmts = append([]goast.Stmt{}, &goast.IfStmt{
				Cond: a,
				Body: &bod,
			})
			stmts = append(stmts, els.List...)
		}
	}

	if bType == cType && bType == "long double" && returnType == "float64" {
		// typically use case
		//
		//		double a = 54;
		//		double b = -4;
		//		double c;
		//		c = a > b ? a : b;
		//
		// now :
		// 	c = func() float64 {
		// 		if b > a {
		// 			return a
		// 		}
		// 		return b
		// 	}()
		//
		// want:
		// 	c = math.Max(a,b)
	}

	return util.NewFuncClosure(
		returnType,
		stmts...), n.Type, preStmts, postStmts, nil
}

// transpileParenExpr transpiles an expression that is wrapped in parentheses.
// There is a special case where "(0)" is treated as a NULL (since that's what
// the macro expands to). We have to return the type as "null" since we don't
// know at this point what the NULL expression will be used in conjunction with.
func transpileParenExpr(n *ast.ParenExpr, p *program.Program) (
	r *goast.ParenExpr, exprType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpile ParenExpr. err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	n.Type = util.GenerateCorrectType(n.Type)
	n.Type2 = util.GenerateCorrectType(n.Type2)

	expr, exprType, preStmts, postStmts, err := atomicOperation(n.Children()[0], p)
	if err != nil {
		return
	}
	if expr == nil {
		err = fmt.Errorf("expr is nil")
		return
	}

	if exprType == types.NullPointer {
		r = &goast.ParenExpr{X: expr}
		return
	}

	if !util.IsFunction(exprType) &&
		exprType != "void" &&
		exprType != "bool" &&
		exprType != types.ToVoid {
		expr, err = types.CastExpr(p, expr, exprType, n.Type)
		if err != nil {
			return
		}
		exprType = n.Type
	}

	var ok bool
	r, ok = expr.(*goast.ParenExpr)
	if !ok {
		r = &goast.ParenExpr{X: expr}
	}

	return
}

func transpileCompoundAssignOperator(
	n *ast.CompoundAssignOperator, p *program.Program, exprIsStmt bool) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileCompoundAssignOperator. err = %v", err)
		}
	}()

	operator := n.Opcode[:len(n.Opcode)-1]

	if len(n.ChildNodes) != 2 {
		err = fmt.Errorf("not enought ChildNodes: %d", len(n.ChildNodes))
		return
	}

	if !types.IsCPointer(n.Type, p) && !types.IsCArray(n.Type, p) {
		return transpileBinaryOperator(&ast.BinaryOperator{
			Type:       n.Type,
			Operator:   n.Opcode,
			ChildNodes: n.ChildNodes,
		}, p, false)
	}

	return transpileBinaryOperator(&ast.BinaryOperator{
		Type:     n.Type,
		Operator: "=",
		ChildNodes: []ast.Node{
			n.ChildNodes[0],
			&ast.BinaryOperator{
				Type:       n.Type,
				Operator:   operator,
				ChildNodes: n.ChildNodes,
			},
		},
	}, p, false)
}

// getTokenForOperator returns the Go operator token for the provided C
// operator.
func getTokenForOperator(operator string) token.Token {
	switch operator {
	// Arithmetic
	case "--":
		return token.DEC
	case "++":
		return token.INC
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	case "%":
		return token.REM

	// Assignment
	case "=":
		return token.ASSIGN
	case "+=":
		return token.ADD_ASSIGN
	case "-=":
		return token.SUB_ASSIGN
	case "*=":
		return token.MUL_ASSIGN
	case "/=":
		return token.QUO_ASSIGN
	case "%=":
		return token.REM_ASSIGN
	case "&=":
		return token.AND_ASSIGN
	case "|=":
		return token.OR_ASSIGN
	case "^=":
		return token.XOR_ASSIGN
	case "<<=":
		return token.SHL_ASSIGN
	case ">>=":
		return token.SHR_ASSIGN

	// Bitwise
	case "&":
		return token.AND
	case "|":
		return token.OR
	case "~":
		return token.XOR
	case ">>":
		return token.SHR
	case "<<":
		return token.SHL
	case "^":
		return token.XOR

	// Comparison
	case ">=":
		return token.GEQ
	case "<=":
		return token.LEQ
	case "<":
		return token.LSS
	case ">":
		return token.GTR
	case "!=":
		return token.NEQ
	case "==":
		return token.EQL

	// Logical
	case "!":
		return token.NOT
	case "&&":
		return token.LAND
	case "||":
		return token.LOR

	// Other
	case ",":
		return token.COMMA
	}

	panic(fmt.Sprintf("unknown operator: %s", operator))
}

func convertToWithoutAssign(operator token.Token) token.Token {
	switch operator {
	case token.ADD_ASSIGN: // "+="
		return token.ADD
	case token.SUB_ASSIGN: // "-="
		return token.SUB
	case token.MUL_ASSIGN: // "*="
		return token.MUL
	case token.QUO_ASSIGN: // "/="
		return token.QUO
	}
	panic(fmt.Sprintf("not support operator: %v", operator))
}

func findUnaryWithInteger(node ast.Node) (*ast.UnaryOperator, bool) {
	switch n := node.(type) {
	case *ast.UnaryOperator:
		return n, true
	case *ast.ParenExpr:
		return findUnaryWithInteger(n.Children()[0])
	}
	return nil, false
}

func atomicOperation(n ast.Node, p *program.Program) (
	expr goast.Expr, exprType string, preStmts, postStmts []goast.Stmt, err error) {

	expr, exprType, preStmts, postStmts, err = transpileToExpr(n, p, false)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot create atomicOperation |%T|. err = %v", n, err)
		}
		if exprType == "" {
			p.AddMessage(p.GenerateWarningMessage(fmt.Errorf("exprType is empty"), n))
		}
	}()

	switch v := n.(type) {
	case *ast.UnaryOperator:
		switch v.Operator {
		case "&", "*", "!", "-", "+", "~":
			return
		}
		// UnaryOperator 0x252d798 <col:17, col:18> 'double' prefix '-'
		// `-FloatingLiteral 0x252d778 <col:18> 'double' 0.000000e+00
		if _, ok := v.Children()[0].(*ast.IntegerLiteral); ok {
			return
		}
		if _, ok := v.Children()[0].(*ast.FloatingLiteral); ok {
			return
		}

		// UnaryOperator 0x3001768 <col:204, col:206> 'int' prefix '++'
		// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
		// OR
		// UnaryOperator 0x3001768 <col:204, col:206> 'int' postfix '++'
		// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
		var varName string
		if vv, ok := v.Children()[0].(*ast.DeclRefExpr); ok {
			varName = vv.Name

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			// operators: ++, --
			if v.IsPrefix {
				// Example:
				// UnaryOperator 0x3001768 <col:204, col:206> 'int' prefix '++'
				// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
				expr = util.NewAnonymousFunction(
					append(preStmts, &goast.ExprStmt{X: expr}),
					nil,
					util.NewIdent(varName),
					exprResolveType)
				preStmts = nil
				break
			}
			// Example:
			// UnaryOperator 0x3001768 <col:204, col:206> 'int' postfix '++'
			// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
			expr = util.NewAnonymousFunction(preStmts,
				[]goast.Stmt{&goast.ExprStmt{X: expr}},
				util.NewIdent(varName),
				exprResolveType)
			preStmts = nil

			break
		}

		// UnaryOperator 'char *' postfix '++'
		// `-ParenExpr 'char *' lvalue
		//   `-UnaryOperator 'char *' lvalue prefix '*'
		//     `-ImplicitCastExpr 'char **' <LValueToRValue>
		//       `-DeclRefExpr 'char **' lvalue Var 0x2699168 'bpp' 'char **'
		//
		// UnaryOperator 'int' postfix '++'
		// `-MemberExpr 'int' lvalue .pos 0x358b538
		//   `-ArraySubscriptExpr 'struct struct_I_A':'struct struct_I_A' lvalue
		//     |-ImplicitCastExpr 'struct struct_I_A *' <ArrayToPointerDecay>
		//     | `-DeclRefExpr 'struct struct_I_A [2]' lvalue Var 0x358b6e8 'siia' 'struct struct_I_A [2]'
		//     `-IntegerLiteral 'int' 0
		varName = "tempVar"

		nextNode := v.Children()[0]
		for {
			if par, ok := nextNode.(*ast.ParenExpr); ok {
				nextNode = par.ChildNodes[0]
				continue
			}
			break
		}
		expr, exprType, preStmts, postStmts, err = transpileToExpr(nextNode, p, false)
		if err != nil {
			return
		}

		var exprResolveType string
		exprResolveType, err = types.ResolveType(p, v.Type)
		if err != nil {
			return
		}

		if types.IsPointer(v.Type, p) {
			switch e := expr.(type) {
			case *goast.IndexExpr:
				if v.Operator == "++" {
					// expr = 'bpp[0]'
					// example of snippet:
					//	func  () []byte{
					//		tempVar = bpp[0]
					//		defer func(){
					//			bpp = bpp[1:]
					//		}()
					//		return tempVar
					//	}
					expr = util.NewAnonymousFunction(
						// body :
						append(preStmts, &goast.AssignStmt{
							Lhs: []goast.Expr{util.NewIdent(varName)},
							Tok: token.DEFINE,
							Rhs: []goast.Expr{expr},
						}),
						// defer :
						append([]goast.Stmt{
							&goast.AssignStmt{
								Lhs: []goast.Expr{
									e,
								},
								Tok: token.ASSIGN,
								Rhs: []goast.Expr{
									&goast.SliceExpr{
										X:      e,
										Low:    goast.NewIdent("1"),
										Slice3: false,
									},
								},
							},
						}, postStmts...),
						// return :
						util.NewIdent(varName),
						exprResolveType)
					preStmts = nil
					postStmts = nil
					return
				}

			case *goast.Ident, *goast.SelectorExpr:
				if v.Operator == "++" {
					// expr = 'p'
					// example of snippet:
					//	func  () [][]byte{
					//		tempVar = p
					//		defer func(){
					//			p = p[1:]
					//		}()
					//		return tempVar
					//	}
					expr = util.NewAnonymousFunction(
						// body :
						append(preStmts, &goast.AssignStmt{
							Lhs: []goast.Expr{util.NewIdent(varName)},
							Tok: token.DEFINE,
							Rhs: []goast.Expr{expr},
						}),
						// defer :
						append([]goast.Stmt{
							&goast.AssignStmt{
								Lhs: []goast.Expr{
									e,
								},
								Tok: token.ASSIGN,
								Rhs: []goast.Expr{
									&goast.SliceExpr{
										X:      e,
										Low:    goast.NewIdent("1"),
										Slice3: false,
									},
								},
							},
						}, postStmts...),
						// return :
						util.NewIdent(varName),
						exprResolveType)
					preStmts = nil
					postStmts = nil
					return
				}

			default:
				// TODO add here
				p.AddMessage(p.GenerateWarningMessage(
					fmt.Errorf("transpilation pointer is not support: %T", e), v))
			}
		}

		body := append(preStmts, &goast.AssignStmt{
			Lhs: []goast.Expr{util.NewIdent(varName)},
			Tok: token.DEFINE,
			Rhs: []goast.Expr{util.NewUnaryExpr(
				expr,
				token.AND,
			)},
		})

		deferBody := postStmts
		postStmts = nil
		preStmts = nil

		switch v.Operator {
		case "++":
			expr = &goast.BinaryExpr{
				X:  &goast.StarExpr{X: util.NewIdent(varName)},
				Op: token.ADD_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}
		case "--":
			expr = &goast.BinaryExpr{
				X:  &goast.StarExpr{X: util.NewIdent(varName)},
				Op: token.SUB_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}
		}

		body = append(body, preStmts...)
		deferBody = append(deferBody, postStmts...)

		// operators: ++, --
		if v.IsPrefix {
			// Example:
			// UnaryOperator 0x3001768 <col:204, col:206> 'int' prefix '++'
			// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
			expr = util.NewAnonymousFunction(
				append(body, &goast.ExprStmt{X: expr}),
				deferBody,
				&goast.StarExpr{
					X: util.NewIdent(varName),
				},
				exprResolveType)
			preStmts = nil
			postStmts = nil
			break
		}
		// Example:
		// UnaryOperator 0x3001768 <col:204, col:206> 'int' postfix '++'
		// `-DeclRefExpr 0x3001740 <col:206> 'int' lvalue Var 0x303e888 'current_test' 'int'
		expr = util.NewAnonymousFunction(body,
			append(deferBody, &goast.ExprStmt{X: expr}),
			&goast.StarExpr{
				X: util.NewIdent(varName),
			},
			exprResolveType)
		preStmts = nil
		postStmts = nil

	case *ast.CompoundAssignOperator:
		// CompoundAssignOperator 0x32911c0 <col:18, col:28> 'int' '-=' ComputeLHSTy='int' ComputeResultTy='int'
		// |-DeclRefExpr 0x3291178 <col:18> 'int' lvalue Var 0x328df60 'iterator' 'int'
		// `-IntegerLiteral 0x32911a0 <col:28> 'int' 2
		if vv, ok := v.Children()[0].(*ast.DeclRefExpr); ok {
			varName := vv.Name

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			expr = util.NewAnonymousFunction(
				append(preStmts, &goast.ExprStmt{X: expr}),
				postStmts,
				util.NewIdent(varName),
				exprResolveType)
			preStmts = nil
			postStmts = nil
			break
		}
		// CompoundAssignOperator 0x27906c8 <line:450:2, col:6> 'double' '+=' ComputeLHSTy='double' ComputeResultTy='double'
		// |-UnaryOperator 0x2790670 <col:2, col:3> 'double' lvalue prefix '*'
		// | `-ImplicitCastExpr 0x2790658 <col:3> 'double *' <LValueToRValue>
		// |   `-DeclRefExpr 0x2790630 <col:3> 'double *' lvalue Var 0x2790570 'p' 'double *'
		// `-IntegerLiteral 0x32911a0 <col:28> 'int' 2
		if vv, ok := v.Children()[0].(*ast.UnaryOperator); ok && vv.IsPrefix && vv.Operator == "*" {
			if vvv, ok := vv.Children()[0].(*ast.ImplicitCastExpr); ok {
				if vvvv, ok := vvv.Children()[0].(*ast.DeclRefExpr); ok {
					if types.IsPointer(vvvv.Type, p) {
						varName := vvvv.Name

						var exprResolveType string
						exprResolveType, err = types.ResolveType(p, v.Type)
						if err != nil {
							return
						}

						expr = util.NewAnonymousFunction(
							append(preStmts, &goast.ExprStmt{X: expr}),
							postStmts,
							util.NewUnaryExpr(
								util.NewIdent(varName),
								token.AND,
							),
							exprResolveType)
						preStmts = nil
						postStmts = nil
						break
					}
				}
			}
		}

		// CompoundAssignOperator 0x32911c0 <col:18, col:28> 'int' '-=' ComputeLHSTy='int' ComputeResultTy='int'
		// |-DeclRefExpr 0x3291178 <col:18> 'int' lvalue Var 0x328df60 'iterator' 'int'
		// `-IntegerLiteral 0x32911a0 <col:28> 'int' 2
		varName := "tempVar"
		expr, exprType, preStmts, postStmts, err = transpileToExpr(v.Children()[0], p, false)
		if err != nil {
			return
		}
		body := append(preStmts, &goast.AssignStmt{
			Lhs: []goast.Expr{util.NewIdent(varName)},
			Tok: token.DEFINE,
			Rhs: []goast.Expr{util.NewUnaryExpr(expr, token.AND)},
		})
		preStmts = nil

		// CompoundAssignOperator 0x27906c8 <line:450:2, col:6> 'double' '+=' ComputeLHSTy='double' ComputeResultTy='double'
		// |-UnaryOperator 0x2790670 <col:2, col:3> 'double' lvalue prefix '*'
		// | `-ImplicitCastExpr 0x2790658 <col:3> 'double *' <LValueToRValue>
		// |   `-DeclRefExpr 0x2790630 <col:3> 'double *' lvalue Var 0x2790570 'p' 'double *'
		// `-ImplicitCastExpr 0x27906b0 <col:6> 'double' <IntegralToFloating>
		//   `-IntegerLiteral 0x2790690 <col:6> 'int' 1
		var newPre, newPost []goast.Stmt
		expr, exprType, newPre, newPost, err = atomicOperation(v.Children()[1], p)
		if err != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

		var exprResolveType string
		exprResolveType, err = types.ResolveType(p, v.Type)
		if err != nil {
			return
		}

		body = append(preStmts, body...)
		body = append(body, &goast.AssignStmt{
			Lhs: []goast.Expr{&goast.StarExpr{
				X: util.NewIdent(varName),
			}},
			Tok: getTokenForOperator(v.Opcode),
			Rhs: []goast.Expr{expr},
		})

		expr = util.NewAnonymousFunction(body, postStmts,
			&goast.StarExpr{
				X: util.NewIdent(varName),
			},
			exprResolveType)
		preStmts = nil
		postStmts = nil

	case *ast.ParenExpr:
		// ParenExpr 0x3c42468 <col:18, col:40> 'int'
		if len(n.Children()) == 1 {
			return atomicOperation(n.Children()[0], p)
		}
		return

	case *ast.ImplicitCastExpr:
		if _, ok := v.Children()[0].(*ast.MemberExpr); ok {
			return
		}
		if _, ok := v.Children()[0].(*ast.IntegerLiteral); ok {
			return
		}
		if v.Kind == "IntegralToPointer" {
			return
		}
		if v.Kind == "BitCast" {
			return
		}

		v.Type = util.GenerateCorrectType(v.Type)
		v.Type2 = util.GenerateCorrectType(v.Type2)

		// avoid problem :
		//
		// constant -1331 overflows uint32
		//
		// ImplicitCastExpr 'unsigned int' <IntegralCast>
		// `-UnaryOperator 'int' prefix '~'
		if t, ok := ast.GetTypeIfExist(v.Children()[0]); ok && !types.IsSigned(p, v.Type) && types.IsSigned(p, *t) {
			if un, ok := n.Children()[0].(*ast.UnaryOperator); ok && un.Operator == "~" {
				var goType string
				goType, err = types.ResolveType(p, v.Type)
				if err != nil {
					return
				}
				expr = util.ConvertToUnsigned(expr, goType)
				return
			}
		}

		// for case : overflow char
		// ImplicitCastExpr 0x2027358 <col:6, col:7> 'char' <IntegralCast>
		// `-UnaryOperator 0x2027338 <col:6, col:7> 'int' prefix '-'
		//   `-IntegerLiteral 0x2027318 <col:7> 'int' 1
		//
		// another example :
		// ImplicitCastExpr 0x2982630 <col:11, col:14> 'char' <IntegralCast>
		// `-ParenExpr 0x2982610 <col:11, col:14> 'int'
		//   `-UnaryOperator 0x29825f0 <col:12, col:13> 'int' prefix '-'
		//     `-IntegerLiteral 0x29825d0 <col:13> 'int' 1
		if v.Type == "char" {
			if len(v.Children()) == 1 {
				if u, ok := findUnaryWithInteger(n.Children()[0]); ok {
					if u.IsPrefix && u.Type == "int" && u.Operator == "-" {
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

		var isSameBaseType bool
		if impl, ok := v.Children()[0].(*ast.ImplicitCastExpr); ok {
			if types.GetBaseType(v.Type) == types.GetBaseType(impl.Type) {
				isSameBaseType = true
			}
		}

		if v.Kind == "PointerToIntegral" {
			if isSameBaseType {
				expr = &goast.IndexExpr{
					X:      expr,
					Lbrack: 1,
					Index: &goast.BasicLit{
						Kind:  token.INT,
						Value: "0",
					},
				}
				exprType = v.Type
				return
			}
			expr = goast.NewIdent("0")
			expr, _ = types.CastExpr(p, expr, "int", v.Type)
			exprType = v.Type
			return
		}

		expr, exprType, preStmts, postStmts, err = atomicOperation(v.Children()[0], p)
		if err != nil {
			return nil, "", nil, nil, err
		}
		if exprType == types.NullPointer {
			return
		}

		var cast bool = true
		if util.IsFunction(exprType) {
			cast = false
		}
		if v.Kind == ast.ImplicitCastExprArrayToPointerDecay {
			cast = false
		}

		if cast {
			expr, err = types.CastExpr(p, expr, exprType, v.Type)
			if err != nil {
				return nil, "", nil, nil, err
			}
			exprType = v.Type
		}
		return

	case *ast.BinaryOperator:
		defer func() {
			if err != nil {
				err = fmt.Errorf("binary operator : `%v`. %v", v.Operator, err)
			}
		}()
		switch v.Operator {
		case ",":
			// BinaryOperator 0x35b95e8 <col:29, col:51> 'int' ','
			// |-UnaryOperator 0x35b94b0 <col:29, col:31> 'int' postfix '++'
			// | `-DeclRefExpr 0x35b9488 <col:29> 'int' lvalue Var 0x35b8dc8 't' 'int'
			// `-CompoundAssignOperator 0x35b95b0 <col:36, col:51> 'int' '+=' ComputeLHSTy='int' ComputeResultTy='int'
			//   |-MemberExpr 0x35b9558 <col:36, col:44> 'int' lvalue .pos 0x35b8730
			//   | `-ArraySubscriptExpr 0x35b9530 <col:36, col:42> 'struct struct_I_A4':'struct struct_I_A4' lvalue
			//   |   |-ImplicitCastExpr 0x35b9518 <col:36> 'struct struct_I_A4 *' <ArrayToPointerDecay>
			//   |   | `-DeclRefExpr 0x35b94d0 <col:36> 'struct struct_I_A4 [2]' lvalue Var 0x35b88d8 'siia' 'struct struct_I_A4 [2]'
			//   |   `-IntegerLiteral 0x35b94f8 <col:41> 'int' 0
			//   `-IntegerLiteral 0x35b9590 <col:51> 'int' 1

			// `-BinaryOperator 0x3c42440 <col:19, col:32> 'int' ','
			//   |-BinaryOperator 0x3c423d8 <col:19, col:30> 'int' '='
			//   | |-DeclRefExpr 0x3c42390 <col:19> 'int' lvalue Var 0x3c3cf60 'iterator' 'int'
			//   | `-IntegerLiteral 0x3c423b8 <col:30> 'int' 0
			//   `-ImplicitCastExpr 0x3c42428 <col:32> 'int' <LValueToRValue>
			//     `-DeclRefExpr 0x3c42400 <col:32> 'int' lvalue Var 0x3c3cf60 'iterator' 'int'
			varName := "tempVar"

			expr, exprType, preStmts, postStmts, err = transpileToExpr(v.Children()[0], p, false)
			if err != nil {
				return
			}

			inBody := combineStmts(&goast.ExprStmt{X: expr}, preStmts, postStmts)

			expr, exprType, preStmts, postStmts, err = atomicOperation(v.Children()[1], p)
			if err != nil {
				return
			}

			if v, ok := expr.(*goast.CallExpr); ok {
				if vv, ok := v.Fun.(*goast.FuncLit); ok {
					vv.Body.List = append(inBody, vv.Body.List...)
					break
				}
			}

			body := append(inBody, preStmts...)
			preStmts = nil

			body = append(body, &goast.AssignStmt{
				Lhs: []goast.Expr{util.NewIdent(varName)},
				Tok: token.DEFINE,
				Rhs: []goast.Expr{util.NewUnaryExpr(expr, token.AND)},
			})

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				err = fmt.Errorf("exprResolveType error for type `%v`: %v", v.Type, err)
				return
			}

			expr = util.NewAnonymousFunction(body, postStmts,
				util.NewUnaryExpr(util.NewIdent(varName), token.MUL),
				exprResolveType)
			preStmts = nil
			postStmts = nil
			exprType = v.Type
			return

		case "=":
			// Find ast.DeclRefExpr in Children[0]
			// Or
			// Find ast.ArraySubscriptExpr in Children[0]
			decl, ok := getDeclRefExprOrArraySub(v.Children()[0])
			if !ok {
				return
			}
			// BinaryOperator 0x2a230c0 <col:8, col:13> 'int' '='
			// |-UnaryOperator 0x2a23080 <col:8, col:9> 'int' lvalue prefix '*'
			// | `-ImplicitCastExpr 0x2a23068 <col:9> 'int *' <LValueToRValue>
			// |   `-DeclRefExpr 0x2a23040 <col:9> 'int *' lvalue Var 0x2a22f20 'a' 'int *'
			// `-IntegerLiteral 0x2a230a0 <col:13> 'int' 42

			// VarDecl 0x328dc50 <col:3, col:29> col:13 used d 'int' cinit
			// `-BinaryOperator 0x328dd98 <col:17, col:29> 'int' '='
			//   |-DeclRefExpr 0x328dcb0 <col:17> 'int' lvalue Var 0x328dae8 'a' 'int'
			//   `-BinaryOperator 0x328dd70 <col:21, col:29> 'int' '='
			//     |-DeclRefExpr 0x328dcd8 <col:21> 'int' lvalue Var 0x328db60 'b' 'int'
			//     `-BinaryOperator 0x328dd48 <col:25, col:29> 'int' '='
			//       |-DeclRefExpr 0x328dd00 <col:25> 'int' lvalue Var 0x328dbd8 'c' 'int'
			//       `-IntegerLiteral 0x328dd28 <col:29> 'int' 42

			// BinaryOperator 0x364a878 <line:139:7, col:23> 'int' '=='
			// |-ParenExpr 0x364a838 <col:7, col:18> 'int'
			// | `-BinaryOperator 0x364a810 <col:8, col:17> 'int' '='
			// |   |-ArraySubscriptExpr 0x364a740 <col:8, col:11> 'int' lvalue
			// |   | |-ImplicitCastExpr 0x364a728 <col:8> 'int *' <ArrayToPointerDecay>
			// |   | | `-DeclRefExpr 0x364a6e0 <col:8> 'int [5]' lvalue Var 0x3648ea0 'l' 'int [5]'
			// |   | `-IntegerLiteral 0x364a708 <col:10> 'int' 0
			// |   `-BinaryOperator 0x364a7e8 <col:15, col:17> 'int' '-'
			// |     |-ImplicitCastExpr 0x364a7b8 <col:15> 'int' <LValueToRValue>
			// |     | `-DeclRefExpr 0x364a768 <col:15> 'int' lvalue Var 0x3647c00 'y' 'int'
			// |     `-ImplicitCastExpr 0x364a7d0 <col:17> 'int' <LValueToRValue>
			// |       `-DeclRefExpr 0x364a790 <col:17> 'int' lvalue Var 0x364a648 's' 'int'
			// `-IntegerLiteral 0x364a858 <col:23> 'int' 3

			var exprResolveType string
			exprResolveType, err = types.ResolveType(p, v.Type)
			if err != nil {
				return
			}

			e, _, newPre, newPost, _ := transpileToExpr(v, p, false)
			body := combineStmts(&goast.ExprStmt{X: e}, newPre, newPost)

			preStmts = nil
			postStmts = nil

			var returnValue goast.Expr
			if bin, ok := e.(*goast.BinaryExpr); ok {
				returnValue = bin.X
			} else {
				returnValue, _, _, _, _ = transpileToExpr(decl, p, false)
				if d, ok := decl.(*ast.DeclRefExpr); ok &&
					types.IsPointer(d.Type, p) && !types.IsPointer(v.Type, p) {
					returnValue = &goast.IndexExpr{
						X: returnValue,
						Index: &goast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
					}
				}
			}

			expr = util.NewAnonymousFunction(body,
				nil,
				returnValue,
				exprResolveType)
			expr = &goast.ParenExpr{
				X:      expr,
				Lparen: 1,
			}
		}

	}

	return
}

// getDeclRefExprOrArraySub - find ast DeclRefExpr
// Examples of input ast trees:
// UnaryOperator 0x2a23080 <col:8, col:9> 'int' lvalue prefix '*'
// `-ImplicitCastExpr 0x2a23068 <col:9> 'int *' <LValueToRValue>
//   `-DeclRefExpr 0x2a23040 <col:9> 'int *' lvalue Var 0x2a22f20 'a' 'int *'
//
// DeclRefExpr 0x328dd00 <col:25> 'int' lvalue Var 0x328dbd8 'c' 'int'
func getDeclRefExprOrArraySub(n ast.Node) (ast.Node, bool) {
	switch v := n.(type) {
	case *ast.DeclRefExpr:
		return v, true
	case *ast.ParenExpr:
		return getDeclRefExprOrArraySub(n.Children()[0])
	case *ast.ImplicitCastExpr:
		return getDeclRefExprOrArraySub(n.Children()[0])
	case *ast.UnaryOperator:
		return getDeclRefExprOrArraySub(n.Children()[0])
	case *ast.ArraySubscriptExpr:
		return v, true
	case *ast.BinaryOperator:
		for i := range v.Children() {
			if v, ok := getDeclRefExprOrArraySub(v.Children()[i]); ok {
				return v, true
			}
		}
	case *ast.MemberExpr:
		return v, true
	}
	return nil, false
}
