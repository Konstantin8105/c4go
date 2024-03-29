// This file contains functions for transpiling binary operator expressions.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

// Comma problem. Example:
// for (int i=0,j=0;i+=1,j<5;i++,j++){...}
// For solving - we have to separate the
// binary operator "," to 2 parts:
// part 1(pre ): left part  - typically one or more some expressions
// part 2(stmt): right part - always only one expression, with or without
//
//	logical operators like "==", "!=", ...
func transpileBinaryOperatorComma(n *ast.BinaryOperator, p *program.Program) (
	stmt goast.Stmt, preStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpile operator comma : err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	left, err := transpileToStmts(n.Children()[0], p)
	if err != nil {
		return nil, nil, err
	}

	right, err := transpileToStmts(n.Children()[1], p)
	if err != nil {
		return nil, nil, err
	}

	if left == nil || right == nil {
		return nil, nil,
			fmt.Errorf("cannot transpile binary operator comma: right = %v , left = %v",
				right, left)
	}

	preStmts = append(preStmts, left...)
	preStmts = append(preStmts, right...)

	if len(preStmts) >= 2 {
		return preStmts[len(preStmts)-1], preStmts[:len(preStmts)-1], nil
	}

	if len(preStmts) == 1 {
		return preStmts[0], nil, nil
	}
	return nil, nil, nil
}

func transpileBinaryOperator(n *ast.BinaryOperator, p *program.Program, exprIsStmt bool) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf(
				"cannot transpile BinaryOperator with type '%s' :"+
					" result type = {%s}. Error: %v", n.Type, eType, err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	operator := getTokenForOperatorNoError(n.Operator)
	n.Type = util.GenerateCorrectType(n.Type)
	n.Type2 = util.GenerateCorrectType(n.Type2)

	defer func() {
		if err != nil {
			err = fmt.Errorf("operator is `%v`. %v", operator, err)
		}
	}()

	// Char overflow
	// BinaryOperator  'int' '!='
	// |-ImplicitCastExpr 'int' <IntegralCast>
	// | `-ImplicitCastExpr 'char' <LValueToRValue>
	// |   `-...
	// `-ParenExpr 'int'
	//   `-UnaryOperator 'int' prefix '-'
	//     `-IntegerLiteral 'int' 1
	if n.Operator == "!=" {
		var leftOk bool
		if l0, ok := n.ChildNodes[0].(*ast.ImplicitCastExpr); ok && l0.Type == "int" {
			if len(l0.ChildNodes) > 0 {
				if l1, ok := l0.ChildNodes[0].(*ast.ImplicitCastExpr); ok && l1.Type == "char" {
					leftOk = true
				}
			}
		}
		if leftOk {
			if r0, ok := n.ChildNodes[1].(*ast.ParenExpr); ok && r0.Type == "int" {
				if len(r0.ChildNodes) > 0 {
					if r1, ok := r0.ChildNodes[0].(*ast.UnaryOperator); ok && r1.IsPrefix && r1.Operator == "-" {
						if r2, ok := r1.ChildNodes[0].(*ast.IntegerLiteral); ok && r2.Type == "int" {
							r0.ChildNodes[0] = &ast.BinaryOperator{
								Type:     "int",
								Type2:    "int",
								Operator: "+",
								ChildNodes: []ast.Node{
									r1,
									&ast.IntegerLiteral{
										Type:  "int",
										Value: "256",
									},
								},
							}
						}
					}
				}
			}
		}
	}

	// Example of C code
	// a = b = 1
	// // Operation equal transpile from right to left
	// Solving:
	// b = 1, a = b
	// // Operation comma transpile from left to right
	// If we have for example:
	// a = b = c = 1
	// then solution is:
	// c = 1, b = c, a = b
	// |-----------|
	// this part, created in according to
	// recursive working
	// Example of AST tree for problem:
	// |-BinaryOperator 0x2f17870 <line:13:2, col:10> 'int' '='
	// | |-DeclRefExpr 0x2f177d8 <col:2> 'int' lvalue Var 0x2f176d8 'x' 'int'
	// | `-BinaryOperator 0x2f17848 <col:6, col:10> 'int' '='
	// |   |-DeclRefExpr 0x2f17800 <col:6> 'int' lvalue Var 0x2f17748 'y' 'int'
	// |   `-IntegerLiteral 0x2f17828 <col:10> 'int' 1
	//
	// Example of AST tree for solution:
	// |-BinaryOperator 0x368e8d8 <line:13:2, col:13> 'int' ','
	// | |-BinaryOperator 0x368e820 <col:2, col:6> 'int' '='
	// | | |-DeclRefExpr 0x368e7d8 <col:2> 'int' lvalue Var 0x368e748 'y' 'int'
	// | | `-IntegerLiteral 0x368e800 <col:6> 'int' 1
	// | `-BinaryOperator 0x368e8b0 <col:9, col:13> 'int' '='
	// |   |-DeclRefExpr 0x368e848 <col:9> 'int' lvalue Var 0x368e6d8 'x' 'int'
	// |   `-ImplicitCastExpr 0x368e898 <col:13> 'int' <LValueToRValue>
	// |     `-DeclRefExpr 0x368e870 <col:13> 'int' lvalue Var 0x368e748 'y' 'int'
	//
	// Example
	// BinaryOperator 'const char *' '='
	// |-...
	// `-ImplicitCastExpr 'const char *' <BitCast>
	//   `-BinaryOperator 'char *' '='
	//     |-...
	//     `-...
	if getTokenForOperatorNoError(n.Operator) == token.ASSIGN {
		child := n.Children()[1]
		if impl, ok := child.(*ast.ImplicitCastExpr); ok {
			child = impl.Children()[0]
		}
		switch c := child.(type) {
		case *ast.BinaryOperator:
			if getTokenForOperatorNoError(c.Operator) == token.ASSIGN {
				bSecond := ast.BinaryOperator{
					Type:     c.Type,
					Operator: "=",
				}
				bSecond.AddChild(n.Children()[0])

				var impl ast.ImplicitCastExpr
				impl.Type = c.Type
				impl.Kind = "LValueToRValue"
				impl.AddChild(c.Children()[0])
				bSecond.AddChild(&impl)

				var bComma ast.BinaryOperator
				bComma.Operator = ","
				bComma.Type = c.Type
				bComma.AddChild(c)
				bComma.AddChild(&bSecond)

				// exprIsStmt now changes to false to stop any AST children from
				// not being safely wrapped in a closure.
				return transpileBinaryOperator(&bComma, p, false)
			}
		}
	}

	// Example of C code
	// a = 1, b = a
	// Solving
	// a = 1; // preStmts
	// b = a; // n
	// Example of AST tree for problem:
	// |-BinaryOperator 0x368e8d8 <line:13:2, col:13> 'int' ','
	// | |-BinaryOperator 0x368e820 <col:2, col:6> 'int' '='
	// | | |-DeclRefExpr 0x368e7d8 <col:2> 'int' lvalue Var 0x368e748 'y' 'int'
	// | | `-IntegerLiteral 0x368e800 <col:6> 'int' 1
	// | `-BinaryOperator 0x368e8b0 <col:9, col:13> 'int' '='
	// |   |-DeclRefExpr 0x368e848 <col:9> 'int' lvalue Var 0x368e6d8 'x' 'int'
	// |   `-ImplicitCastExpr 0x368e898 <col:13> 'int' <LValueToRValue>
	// |     `-DeclRefExpr 0x368e870 <col:13> 'int' lvalue Var 0x368e748 'y' 'int'
	//
	// Example of AST tree for solution:
	// |-BinaryOperator 0x21a7820 <line:13:2, col:6> 'int' '='
	// | |-DeclRefExpr 0x21a77d8 <col:2> 'int' lvalue Var 0x21a7748 'y' 'int'
	// | `-IntegerLiteral 0x21a7800 <col:6> 'int' 1
	// |-BinaryOperator 0x21a78b0 <line:14:2, col:6> 'int' '='
	// | |-DeclRefExpr 0x21a7848 <col:2> 'int' lvalue Var 0x21a76d8 'x' 'int'
	// | `-ImplicitCastExpr 0x21a7898 <col:6> 'int' <LValueToRValue>
	// |   `-DeclRefExpr 0x21a7870 <col:6> 'int' lvalue Var 0x21a7748 'y' 'int'
	if getTokenForOperatorNoError(n.Operator) == token.COMMA {
		stmts, _, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
		if err != nil {
			err = fmt.Errorf("cannot transpile expr `token.COMMA` child 0. %v", err)
			return nil, "unknown50", nil, nil, err
		}
		preStmts = append(preStmts, newPre...)
		preStmts = append(preStmts, util.NewExprStmt(stmts))
		preStmts = append(preStmts, newPost...)

		var st string
		stmts, st, newPre, newPost, err = transpileToExpr(n.Children()[1], p, false)
		if err != nil {
			err = fmt.Errorf("cannot transpile expr `token.COMMA` child 1. %v", err)
			return nil, "unknown51", nil, nil, err
		}
		// Theoretically , we don't have any preStmts or postStmts
		// from n.Children()[1]
		if len(newPre) > 0 || len(newPost) > 0 {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("not support length pre or post stmts: {%d,%d}",
					len(newPre), len(newPost)), n))
		}
		return stmts, st, preStmts, postStmts, nil
	}

	// pointer arithmetic
	if types.IsPointer(n.Type, p) {
		if operator == token.ADD || // +
			false {

			// not acceptable binaryOperator with operator `-`
			haveSub := false
			{
				var check func(ast.Node)
				check = func(node ast.Node) {
					if node == nil {
						return
					}
					if bin, ok := node.(*ast.BinaryOperator); ok && bin.Operator == "-" {
						haveSub = true
					}
					for i := range node.Children() {
						check(node.Children()[i])
					}
				}
				check(n)
			}

			if !haveSub {

				fakeUnary := &ast.UnaryOperator{
					Type:     n.Type,
					Operator: "*",
					ChildNodes: []ast.Node{
						n,
					},
				}

				var newPre, newPost []goast.Stmt
				expr, eType, newPre, newPost, err =
					transpilePointerArith(fakeUnary, p)
				eType = n.Type

				if err != nil {
					return
				}
				if expr == nil {
					return nil, "", nil, nil, fmt.Errorf("expr is nil")
				}
				preStmts, postStmts =
					combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

				if ind, ok := expr.(*goast.IndexExpr); ok {
					expr = &goast.SliceExpr{
						X:      ind.X,
						Low:    ind.Index,
						Slice3: false,
					}
				}

				return
			}
		}
	}

	left, leftType, newPre, newPost, err := atomicOperation(n.Children()[0], p)
	if err != nil {
		err = fmt.Errorf("cannot atomic for left part. %v", err)
		return nil, "unknown52", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	right, rightType, newPre, newPost, err := atomicOperation(n.Children()[1], p)
	if err != nil {
		err = fmt.Errorf("cannot atomic for right part. %v", err)
		return nil, "unknown53", nil, nil, err
	}

	if types.IsPointer(leftType, p) && types.IsPointer(rightType, p) {
		switch operator {
		case token.SUB: // -
			p.AddImport("unsafe")
			var sizeof int
			sizeof, err = types.SizeOf(p, types.GetBaseType(leftType))
			if err != nil {
				return nil, "PointerOperation_unknown01", nil, nil, err
			}
			var e goast.Expr
			var newPost []goast.Stmt
			e, newPost, err = SubTwoPnts(p, left, leftType, right, rightType, sizeof)
			if err != nil {
				return nil, "PointerOperation_unknown02", nil, nil, err
			}
			postStmts = append(postStmts, newPost...)

			expr, err = types.CastExpr(p, e, "long long", n.Type)
			if err != nil {
				return nil, "PointerOperation_unknown03", nil, nil, err
			}
			eType = n.Type
			return

		case token.GTR, token.GEQ, // >  >=
			token.LOR,            // ||
			token.LAND,           // &&
			token.LSS, token.LEQ, // <  <=
			token.EQL, token.NEQ: // == !=

			// IfStmt 0x2369a68 <line:210:3, line:211:11>
			// |-BinaryOperator 0x2369a00 <line:210:7, col:21> 'int' '=='
			// | |-ImplicitCastExpr 0x23699d0 <col:7, col:16> 'struct font *' <LValueToRValue>
			// | | `-ArraySubscriptExpr 0x2369990 <col:7, col:16> 'struct font *' lvalue
			// | |   |-ImplicitCastExpr 0x2369960 <col:7> 'struct font **' <ArrayToPointerDecay>
			// | |   | `-DeclRefExpr 0x2369920 <col:7> 'struct font *[32]' lvalue Var 0x235fbe8 'fn_font' 'struct font *[32]'
			// | |   `-ImplicitCastExpr 0x2369978 <col:15> 'int' <LValueToRValue>
			// | |     `-DeclRefExpr 0x2369940 <col:15> 'int' lvalue Var 0x2369790 'i' 'int'
			// | `-ImplicitCastExpr 0x23699e8 <col:21> 'struct font *' <LValueToRValue>
			// |   `-DeclRefExpr 0x23699b0 <col:21> 'struct font *' lvalue ParmVar 0x2369638 'fn' 'struct font *'
			// `-ReturnStmt 0x2369a58 <line:211:4, col:11>
			//   `-ImplicitCastExpr 0x2369a40 <col:11> 'int' <LValueToRValue>
			//     `-DeclRefExpr 0x2369a20 <col:11> 'int' lvalue Var 0x2369790 'i' 'int'

			var sizeof int
			baseType := types.GetBaseType(leftType)
			sizeof, err = types.SizeOf(p, baseType)
			if err != nil {
				err = fmt.Errorf("{'%s' %v '%s'}. sizeof = %d for baseType = '%s'. %v",
					leftType, operator, rightType, sizeof, baseType, err)
				return nil, "PointerOperation_unknown04", nil, nil, err
			}
			var e goast.Expr
			var newPost []goast.Stmt
			e, newPost, err = PntCmpPnt(
				p,
				left, leftType,
				right, rightType,
				sizeof, operator,
			)
			if err != nil {
				err = fmt.Errorf("{'%s' %v '%s'}. for base type: `%s`. %v",
					leftType, operator, rightType, baseType, err)
				return nil, "PointerOperation_unknown05", nil, nil, err
			}
			postStmts = append(postStmts, newPost...)
			expr = e
			eType = "bool"

			return

		case token.ASSIGN: // =
			// ignore

		default:
			err = fmt.Errorf("not implemented pointer operation: %v", operator)
			return
		}
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	returnType := types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType)

	if operator == token.LAND || operator == token.LOR { // && ||
		left, err = types.CastExpr(p, left, leftType, "bool")
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
			// ignore error
			left = util.NewNil()
			err = nil
		}

		right, err = types.CastExpr(p, right, rightType, "bool")
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
			// ignore error
			right = util.NewNil()
			err = nil
		}

		resolvedLeftType, err := types.ResolveType(p, leftType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}

		expr := util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt)

		return expr, "bool", preStmts, postStmts, nil
	}

	// The right hand argument of the shift left or shift right operators
	// in Go must be unsigned integers. In C, shifting with a negative shift
	// count is undefined behaviour (so we should be able to ignore that case).
	// To handle this, cast the shift count to a uint64.
	if operator == token.SHL || // <<
		operator == token.SHR || // >>
		operator == token.SHL_ASSIGN || // <<=
		operator == token.SHR_ASSIGN || // <<=
		false {
		right, err = types.CastExpr(p, right, rightType, "unsigned long long")
		p.AddMessage(p.GenerateWarningMessage(err, n))
		if right == nil {
			right = util.NewNil()
		}
	}

	// pointer arithmetic
	if types.IsPointer(n.Type, p) {
		if operator == token.ADD || // +
			operator == token.SUB || // -
			false {

			el, extl, prel, postl, errl := atomicOperation(n.Children()[0], p)
			er, extr, prer, postr, errr := atomicOperation(n.Children()[1], p)
			if errl != nil || errr != nil {
				err = fmt.Errorf("pointer operation is not valid : %v. %v", errl, errr)
				return
			}
			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, prel, postl)
			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, prer, postr)

			if types.IsCPointer(extl, p) {
				expr, eType, prel, postl, errl = pointerArithmetic(p, el, extl, er, extr, operator)
			} else {
				expr, eType, prel, postl, errl = pointerArithmetic(p, er, extr, el, extl, operator)
			}

			if errl != nil {
				err = fmt.Errorf("pointer operation is not valid : %v", errl)
				return
			}

			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, prel, postl)

			return
		}
	}

	if operator == token.NEQ || // !=
		operator == token.EQL || // ==
		operator == token.LSS || // <
		operator == token.GTR || // >
		operator == token.AND || // &
		operator == token.ADD || // +
		operator == token.SUB || // -
		operator == token.MUL || // *
		operator == token.QUO || // /
		operator == token.REM || // %
		operator == token.LEQ || // <=
		operator == token.GEQ || // >=
		operator == token.ADD_ASSIGN || // +=
		operator == token.SUB_ASSIGN || // -=
		operator == token.MUL_ASSIGN || // *=
		operator == token.QUO_ASSIGN || // /=
		operator == token.REM_ASSIGN || // %=

		operator == token.AND_ASSIGN || // &=
		operator == token.OR_ASSIGN || // |=
		operator == token.XOR_ASSIGN || // ^=
		operator == token.SHL_ASSIGN || // <<=
		operator == token.SHR_ASSIGN || // >>=
		operator == token.AND_NOT_ASSIGN || // &^=
		false {

		if rightType == types.NullPointer && leftType == types.NullPointer {
			// example C code :
			// if ( NULL == NULL )
			right = goast.NewIdent("1")
			rightType = "int"
			left = goast.NewIdent("1")
			leftType = "int"

		} else if rightType != types.NullPointer && leftType != types.NullPointer {
			// We may have to cast the right side to the same type as the left
			// side. This is a bit crude because we should make a better
			// decision of which type to cast to instead of only using the type
			// of the left side.

			if operator == token.ADD || // +
				operator == token.SUB || // -
				operator == token.MUL || // *
				operator == token.QUO || // /
				operator == token.REM || // %
				false {

				if rightType == "bool" {
					right, err = types.CastExpr(p, right, rightType, "int")
					rightType = "int"
					p.AddMessage(p.GenerateWarningMessage(err, n))
				}
				if leftType == "bool" {
					left, err = types.CastExpr(p, left, leftType, "int")
					leftType = "int"
					p.AddMessage(p.GenerateWarningMessage(err, n))
				}
			}
			right, err = types.CastExpr(p, right, rightType, leftType)
			rightType = leftType
			p.AddMessage(p.GenerateWarningMessage(err, n))

			// compare pointers
			//
			// BinaryOperator 'int' '<'
			// |-ImplicitCastExpr 'char *' <LValueToRValue>
			// | `-DeclRefExpr 'char *' lvalue Var 0x26ba988 'c' 'char *'
			// `-ImplicitCastExpr 'char *' <LValueToRValue>
			//   `-DeclRefExpr 'char *' lvalue Var 0x26ba8a8 'b' 'char *'
			if types.IsPointer(leftType, p) || types.IsPointer(rightType, p) {
				err = fmt.Errorf("need add pointer operator : %s %v %s",
					leftType, n.Operator, rightType)
				return
			}
		}
	}

	if operator == token.ASSIGN { // =

		// BinaryOperator 'double *' '='
		// |-DeclRefExpr 'double *' lvalue Var 0x2a7fa48 'd' 'double *'
		// `-ImplicitCastExpr 'double *' <BitCast>
		//   `-CStyleCastExpr 'char *' <BitCast>
		//     `-...
		right, err = types.CastExpr(p, right, rightType, returnType)
		rightType = returnType
		if err != nil {
			return
		}

		if _, ok := right.(*goast.UnaryExpr); ok && types.IsDereferenceType(rightType) {
			deref, err := types.GetDereferenceType(rightType)

			if !p.AddMessage(p.GenerateWarningMessage(err, n)) {
				resolvedDeref, err := types.ResolveType(p, deref)

				// FIXME: I'm not sure how this situation arises.
				if resolvedDeref == "" {
					resolvedDeref = "interface{}"
				}

				if !p.AddMessage(p.GenerateWarningMessage(err, n)) {
					p.AddImport("unsafe")
					right = CreateSliceFromReference(resolvedDeref, right)
				}
			}
		}

		if p.AddMessage(p.GenerateWarningMessage(err, n)) && right == nil {
			right = util.NewNil()
		}
	}

	var resolvedLeftType = n.Type
	if !util.IsFunction(n.Type) && !types.IsTypedefFunction(p, n.Type) {
		if leftType != types.NullPointer {
			resolvedLeftType, err = types.ResolveType(p, leftType)
		} else {
			resolvedLeftType, err = types.ResolveType(p, rightType)
		}
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	// Enum casting
	if operator != token.ASSIGN && strings.Contains(rightType, "enum ") {
		right, err = types.CastExpr(p, right, rightType, "int")
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	if left == nil {
		err = fmt.Errorf("left part of binary operation is nil. left : %#v", n.Children()[0])
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return nil, "", nil, nil, err
	}

	if right == nil {
		err = fmt.Errorf("right part of binary operation is nil. right : %#v", n.Children()[1])
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return nil, "", nil, nil, err
	}

	return util.NewBinaryExpr(left, operator, right, resolvedLeftType, exprIsStmt),
		types.ResolveTypeForBinaryOperator(p, n.Operator, leftType, rightType),
		preStmts, postStmts, nil
}

func foundCallExpr(n ast.Node) *ast.CallExpr {
	switch v := n.(type) {
	case *ast.ImplicitCastExpr, *ast.CStyleCastExpr:
		return foundCallExpr(n.Children()[0])
	case *ast.CallExpr:
		return v
	}
	return nil
}
