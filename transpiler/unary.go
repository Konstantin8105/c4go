// This file contains functions for transpiling unary operator expressions.

package transpiler

import (
	"fmt"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"

	goast "go/ast"
	"go/token"
)

func transpileUnaryOperatorInc(n *ast.UnaryOperator, p *program.Program, operator token.Token) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileUnaryOperatorInc. err = %v", err)
		}
	}()

	if !(operator == token.INC || operator == token.DEC) {
		err = fmt.Errorf("not acceptable operator '%v'", operator)
		return
	}

	if util.IsPointer(n.Type) {
		switch operator {
		case token.INC: // ++
			operator = token.ADD
		case token.DEC: // --
			operator = token.SUB
		}
		// remove paren - ()
	remove_paren:
		if p, ok := n.Children()[0].(*ast.ParenExpr); ok {
			n.Children()[0] = p.Children()[0]
			goto remove_paren
		}

		var left goast.Expr
		var leftType string
		var newPre, newPost []goast.Stmt
		left, leftType, newPre, newPost, err = transpileToExpr(
			n.Children()[0], p, false)
		if err != nil {
			return
		}

		preStmts, postStmts = combinePreAndPostStmts(
			preStmts, postStmts, newPre, newPost)

		rightType := "int"
		right := &goast.BasicLit{
			Kind:  token.INT,
			Value: "1",
		}

		expr, eType, newPre, newPost, err = pointerArithmetic(
			p, left, leftType, right, rightType, operator)
		if err != nil {
			return
		}
		if expr == nil {
			return nil, "", nil, nil, fmt.Errorf("Expr is nil")
		}

		preStmts, postStmts = combinePreAndPostStmts(
			preStmts, postStmts, newPre, newPost)

		var name string
		name, err = getName(p, n.Children()[0])
		if err != nil {
			return
		}

		found := false
		if sl, ok := expr.(*goast.SliceExpr); ok {
			if ind, ok := sl.X.(*goast.IndexExpr); ok {
				expr = &goast.BinaryExpr{
					X:  ind,
					Op: token.ASSIGN,
					Y:  expr,
				}
				found = true
			}
		}
		if !found {
			expr = &goast.BinaryExpr{
				X:  goast.NewIdent(name),
				Op: token.ASSIGN,
				Y:  expr,
			}
		}
		return
	}

	if v, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
		switch n.Operator {
		case "++":
			return &goast.BinaryExpr{
				X:  util.NewIdent(v.Name),
				Op: token.ADD_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}, n.Type, nil, nil, nil
		case "--":
			return &goast.BinaryExpr{
				X:  util.NewIdent(v.Name),
				Op: token.SUB_ASSIGN,
				Y:  &goast.BasicLit{Kind: token.INT, Value: "1"},
			}, n.Type, nil, nil, nil
		}
	}

	// Unfortunately we cannot use the Go increment operators because we are not
	// providing any position information for tokens. This means that the ++/--
	// would be placed before the expression and would be invalid in Go.
	//
	// Until it can be properly fixed (can we trick Go into to placing it after
	// the expression with a magic position?) we will have to return a
	// BinaryExpr with the same functionality.

	binaryOperator := "+="
	if operator == token.DEC {
		binaryOperator = "-="
	}

	return transpileBinaryOperator(&ast.BinaryOperator{
		Type:     n.Type,
		Operator: binaryOperator,
		ChildNodes: []ast.Node{
			n.Children()[0], &ast.IntegerLiteral{
				Type:       "int",
				Value:      "1",
				ChildNodes: []ast.Node{},
			},
		},
	}, p, false)
}

func transpileUnaryOperatorNot(n *ast.UnaryOperator, p *program.Program) (
	goast.Expr, string, []goast.Stmt, []goast.Stmt, error) {
	e, eType, preStmts, postStmts, err := atomicOperation(n.Children()[0], p)
	if err != nil {
		return nil, "", nil, nil, err
	}

	// UnaryOperator <> 'int' prefix '!'
	// `-ImplicitCastExpr <> 'int (*)(int, double)' <LValueToRValue>
	//   `-DeclRefExpr <> 'int (*)(int, double)' lvalue Var 0x2be4e80 'T' 'int (*)(int, double)'
	if util.IsFunction(eType) {
		return &goast.BinaryExpr{
			X:  e,
			Op: token.EQL, // ==
			Y:  goast.NewIdent("nil"),
		}, "bool", preStmts, postStmts, nil
	}

	// specific case:
	//
	// UnaryOperator 'int' prefix '!'
	// `-ParenExpr 'int'
	//   `-BinaryOperator 'int' '='
	//     |-DeclRefExpr 'int' lvalue Var 0x3329b60 'y' 'int'
	//     `-ImplicitCastExpr 'int' <LValueToRValue>
	//       `-DeclRefExpr 'int' lvalue Var 0x3329ab8 'p' 'int'
	if par, ok := e.(*goast.ParenExpr); ok {
		if bi, ok := par.X.(*goast.BinaryExpr); ok {
			if bi.Op == token.ASSIGN { // =
				preStmts = append(preStmts, &goast.ExprStmt{
					X: bi,
				})
				e = bi.X
			}
		}
	}

	// null in C is zero
	if eType == types.NullPointer {
		e = &goast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		}
		eType = "int"
	}

	if eType == "bool" {
		return &goast.UnaryExpr{
			X:  e,
			Op: token.NOT,
		}, "bool", preStmts, postStmts, nil
	}

	if strings.HasSuffix(eType, "*") {
		// `!pointer` has to be converted to `pointer == nil`
		return &goast.BinaryExpr{
			X:  e,
			Op: token.EQL,
			Y:  util.NewIdent("nil"),
		}, "bool", preStmts, postStmts, nil
	}

	t, err := types.ResolveType(p, eType)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	if t == "[]byte" {
		return util.NewUnaryExpr(
			token.NOT, util.NewCallExpr("noarch.CStringIsNull", e),
		), "bool", preStmts, postStmts, nil
	}

	// only if added "stdbool.h"
	if p.IncludeHeaderIsExists("stdbool.h") {
		if t == "_Bool" {
			t = "int"
			e = util.NewCallExpr("int", e)
		}
	}

	p.AddImport("github.com/Konstantin8105/c4go/noarch")

	functionName := fmt.Sprintf("noarch.Not%s",
		util.GetExportedName(t))

	return util.NewCallExpr(functionName, e),
		eType, preStmts, postStmts, nil
}

// tranpileUnaryOperatorAmpersant - operator ampersant &
// Example of AST:
//
// UnaryOperator 'int (*)[5]' prefix '&'
// `-DeclRefExpr 'int [5]' lvalue Var 0x2d0fb20 'arr' 'int [5]'
//
// UnaryOperator 'char **' prefix '&'
// `-DeclRefExpr 'char *' lvalue Var 0x39b95f0 'line' 'char *'
//
func transpileUnaryOperatorAmpersant(n *ast.UnaryOperator, p *program.Program) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileUnaryOperatorAmpersant : err = %v", err)
		}
	}()

	expr, eType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}
	if expr == nil {
		err = fmt.Errorf("Expr is nil")
		return
	}

	if util.IsFunction(eType) {
		return
	}

	if util.IsLastArray(eType) {
		// In : eType = 'int [5]'
		// Out: eType = 'int *'
		f := strings.Index(eType, "[")
		e := strings.Index(eType, "]")
		if e == len(eType)-1 {
			eType = eType[:f] + "*"
		} else {
			eType = eType[:f] + "*" + eType[e+1:]
		}
		return
	}

	// In : eType = 'int'
	// Out: eType = 'int *'
	// FIXME: This will need to use a real slice to reference the original
	// value.
	resolvedType, err := types.ResolveType(p, eType)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return
	}

	p.AddImport("unsafe")
	expr = util.CreateSliceFromReference(resolvedType, expr)

	// We now have a pointer to the original type.
	eType += " *"
	return
}

// transpilePointerArith - transpile pointer aripthmetic
// Example of using:
// *(t + 1) = ...
func transpilePointerArith(n *ast.UnaryOperator, p *program.Program) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	// pointer - expression with name of array pointer
	var pointer ast.Node

	// locationPointer
	var locPointer ast.Node
	var locPosition int

	// counter - count of amount of changes in AST tree
	var counter int

	var parents []ast.Node
	var found bool

	var f func(ast.Node)
	f = func(n ast.Node) {
		for i := range n.Children() {
			state := func() {
				counter++
				if counter > 1 {
					err = fmt.Errorf("Not acceptable :"+
						" change counter is more then 1. found = %T,%T",
						pointer, n.Children()[i])
					return
				}
				// found pointer
				pointer = n.Children()[i]
				// Replace pointer to zero
				var zero ast.IntegerLiteral
				zero.Type = "int"
				zero.Value = "0"
				locPointer = n
				locPosition = i
				n.Children()[i] = &zero
				found = true
			}

			switch v := n.Children()[i].(type) {
			case *ast.ArraySubscriptExpr,
				*ast.UnaryOperator,
				*ast.VAArgExpr,
				*ast.DeclRefExpr:
				state()
				return

			case *ast.CStyleCastExpr:
				if v.Type == "int" {
					continue
				}
				state()
				return

			case *ast.MemberExpr:
				// check - if member of union
				a := n.Children()[i]
				var isUnion bool
				for {
					if a == nil {
						break
					}
					if len(a.Children()) == 0 {
						break
					}
					switch vv := a.Children()[0].(type) {
					case *ast.MemberExpr, *ast.DeclRefExpr:
						var typeVV string
						switch vt := vv.(type) {
						case *ast.MemberExpr:
							typeVV = vt.Type
						case *ast.DeclRefExpr:
							typeVV = vt.Type
						}
						typeVV = types.GetBaseType(typeVV)

						if _, ok := p.Structs[typeVV]; ok {
							isUnion = true
						}
						if _, ok := p.Structs["struct "+typeVV]; ok {
							isUnion = true
						}
						if strings.HasPrefix(typeVV, "union ") || strings.HasPrefix(typeVV, "struct ") {
							isUnion = true
						}
						if isUnion {
							break
						}
						a = vv
						continue
					case *ast.ImplicitCastExpr, *ast.CStyleCastExpr,
						*ast.ArraySubscriptExpr:
						a = vv
						continue
					}
					break
				}
				if isUnion {
					state()
					return
				}
				// member of struct
				f(v)

			case *ast.CallExpr:
				if v.Type == "int" {
					continue
				}
				state()
				return

			default:
				if found {
					break
				}
				if len(v.Children()) > 0 {
					if found {
						break
					}
					parents = append(parents, v)
					deep := true
					if vv, ok := v.(*ast.ImplicitCastExpr); ok && types.IsCInteger(p, vv.Type) {
						deep = false
					}
					if vv, ok := v.(*ast.CStyleCastExpr); ok && types.IsCInteger(p, vv.Type) {
						deep = false
					}
					if deep {
						f(v)
					}
					if !found {
						parents = parents[:len(parents)-1]
					}
				}
			}
		}
	}
	f(n)

	if err != nil {
		return
	}

	if pointer == nil {
		err = fmt.Errorf("pointer is nil")
		return
	}

	defer func() {
		if pointer != nil && locPointer != nil {
			locPointer.Children()[locPosition] = pointer.(ast.Node)
		}
	}()

	var typesParentBefore []string
	for i := range parents {
		if t, ok := ast.GetTypeIfExist(parents[i]); ok {
			typesParentBefore = append(typesParentBefore, *t)
			*t = "int"
		} else {
			panic(fmt.Errorf("Not support parent type %T in pointer seaching", parents[i]))
		}
	}
	defer func() {
		for i := range parents {
			if t, ok := ast.GetTypeIfExist(parents[i]); ok {
				*t = typesParentBefore[i]
			} else {
				panic(fmt.Errorf("Not support parent type %T in pointer seaching", parents[i]))
			}
		}
	}()

	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
	eType = n.Type

	switch v := pointer.(type) {
	case *ast.MemberExpr:
		arr, _, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		return &goast.IndexExpr{
			X:     arr,
			Index: e,
		}, eType, preStmts, postStmts, err

	case *ast.DeclRefExpr:
		return &goast.IndexExpr{
			X:     util.NewIdent(v.Name),
			Index: e,
		}, eType, preStmts, postStmts, err

	case *ast.CStyleCastExpr, *ast.VAArgExpr, *ast.CallExpr, *ast.ArraySubscriptExpr:
		arr, _, newPre, newPost, err2 := transpileToExpr(v, p, false)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				Lparen: 1,
				X:      arr,
			},
			Index: e,
		}, eType, preStmts, postStmts, err

	case *ast.UnaryOperator:
		arr, _, newPre, newPost, err2 := atomicOperation(v, p)
		if err2 != nil {
			return
		}
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		if memberName, ok := getMemberName(n.Children()[0]); ok {
			return &goast.IndexExpr{
				X: &goast.SelectorExpr{
					X:   arr,
					Sel: util.NewIdent(memberName),
				},
				Index: &goast.BasicLit{
					Kind:  token.INT,
					Value: "0",
				},
			}, eType, preStmts, postStmts, err
		}
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				Lparen: 1,
				X:      arr,
			},
			Index: e,
		}, eType, preStmts, postStmts, err

	}
	return nil, "", nil, nil, fmt.Errorf("Cannot found : %#v", pointer)
}

func transpileUnaryOperator(n *ast.UnaryOperator, p *program.Program) (
	_ goast.Expr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile UnaryOperator: err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	operator := getTokenForOperator(n.Operator)

	switch operator {
	case token.MUL: // *
		// Prefix "*" is not a multiplication.
		// Prefix "*" used for pointer ariphmetic
		// Example of using:
		// *(t + 1) = ...
		return transpilePointerArith(n, p)
	case token.INC, token.DEC: // ++, --
		return transpileUnaryOperatorInc(n, p, operator)
	case token.NOT: // !
		return transpileUnaryOperatorNot(n, p)
	case token.AND: // &
		return transpileUnaryOperatorAmpersant(n, p)
	}

	// Otherwise handle like a unary operator.
	e, eType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	return &goast.UnaryExpr{
		Op: operator,
		X:  e,
	}, eType, preStmts, postStmts, nil

}

func transpileUnaryExprOrTypeTraitExpr(n *ast.UnaryExprOrTypeTraitExpr, p *program.Program) (
	*goast.BasicLit, string, []goast.Stmt, []goast.Stmt, error) {
	t := n.Type2

	// It will have children if the sizeof() is referencing a variable.
	// Fortunately clang already has the type in the AST for us.
	if len(n.Children()) > 0 {
		var realFirstChild interface{}
		t = ""

		switch c := n.Children()[0].(type) {
		case *ast.ParenExpr:
			realFirstChild = c.Children()[0]
		case *ast.DeclRefExpr:
			t = c.Type
		case *ast.UnaryOperator:
			t = c.Type
		case *ast.MemberExpr:
			t = c.Type
		case *ast.ArraySubscriptExpr:
			t = c.Type
		default:
			panic(fmt.Sprintf("cannot find first child from: %#v", n.Children()[0]))
		}

		if t == "" {
			if node, ok := realFirstChild.(ast.Node); ok {
				if ty, ok := ast.GetTypeIfExist(node); ok {
					t = *ty
				} else {
					panic(fmt.Sprintf("cannot do unary on: %#v", realFirstChild))
				}
			}
		}
	}

	sizeInBytes, err := types.SizeOf(p, t)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	return util.NewIntLit(sizeInBytes), n.Type1, nil, nil, nil
}

func transpileStmtExpr(n *ast.StmtExpr, p *program.Program) (
	*goast.CallExpr, string, []goast.Stmt, []goast.Stmt, error) {
	returnType, err := types.ResolveType(p, n.Type)
	if err != nil {
		return nil, "", nil, nil, err
	}

	body, pre, post, err := transpileCompoundStmt(n.Children()[0].(*ast.CompoundStmt), p)
	if err != nil {
		return nil, "", pre, post, err
	}

	// The body of the StmtExpr is always a CompoundStmt. However, the last
	// statement needs to be transformed into an explicit return statement.
	lastStmt := body.List[len(body.List)-1]
	body.List[len(body.List)-1] = &goast.ReturnStmt{
		Results: []goast.Expr{lastStmt.(*goast.ExprStmt).X},
	}

	return util.NewFuncClosure(returnType, body.List...), n.Type, pre, post, nil
}
