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

	// for values
	if v, ok := n.Children()[0].(*ast.DeclRefExpr); ok &&
		!types.IsCArray(v.Type, p) && !types.IsCPointer(v.Type, p) {
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

	// for other case
	op := "+"
	if operator == token.DEC {
		op = "-"
	}

	if len(n.ChildNodes) != 1 {
		err = fmt.Errorf("not enought ChildNodes: %d", len(n.ChildNodes))
		return
	}

	if !types.IsCPointer(n.Type, p) && !types.IsCArray(n.Type, p) {

		binaryOperator := "+="
		if operator == token.DEC {
			binaryOperator = "-="
		}

		return transpileBinaryOperator(&ast.BinaryOperator{
			Type:     n.Type,
			Operator: binaryOperator,
			ChildNodes: []ast.Node{
				n.ChildNodes[0],
				&ast.IntegerLiteral{
					Type:  "int",
					Value: "1",
				},
			},
		}, p, false)
	}

	return transpileBinaryOperator(&ast.BinaryOperator{
		Type:     n.Type,
		Operator: "=",
		ChildNodes: []ast.Node{
			n.ChildNodes[0],
			&ast.BinaryOperator{
				Type:     n.Type,
				Operator: op,
				ChildNodes: append(n.ChildNodes, &ast.IntegerLiteral{
					Type:  "int",
					Value: "1",
				}),
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
			t = "int32"
			e = util.NewCallExpr("int32", e)
		}
	}

	p.AddImport("github.com/Konstantin8105/c4go/noarch")

	functionName := fmt.Sprintf("noarch.Not%s",
		util.GetExportedName(t))
	eType = "int"

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
// UnaryOperator 'float *' prefix '&'
// `-DeclRefExpr 'float' lvalue Var 0x409e2a0 't' 'float'
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
			X:      ind.X,
			Low:    ind.Index,
			Slice3: false,
		}
		eType = n.Type
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

	// We now have a pointer to the original type.
	eType += " *"

	p.AddImport("unsafe")

	// UnaryOperator 'float *' prefix '&'
	// `-DeclRefExpr 'float' lvalue Var 0x409e2a0 't' 'float'
	if len(n.Children()) == 1 {
		if decl, ok := n.Children()[0].(*ast.DeclRefExpr); ok {
			// can simplify
			p.UnsafeConvertValueToPointer = append(p.UnsafeConvertValueToPointer,
				resolvedType)
			expr = util.NewCallExpr(fmt.Sprintf("%s%s", unsafeConvertFunctionName, resolvedType),
				&goast.UnaryExpr{
					Op: token.AND,
					X:  goast.NewIdent(decl.Name),
				})
			return
		}
	}

	expr = util.CreateSliceFromReference(resolvedType, expr)

	return
}

// pointerParts - separate to pointer and value.
// 	* change type for all nodes to `int`
//
// BinaryOperator <col:13, col:57> 'int *' '-'
// |-BinaryOperator <col:13, col:40> 'int *' '+'
// | |-BinaryOperator <col:13, col:32> 'int *' '+'
// | | |-BinaryOperator <col:13, col:21> 'int *' '+'
// | | | |-BinaryOperator <col:13, col:17> 'int' '+'
// | | | | |-IntegerLiteral <col:13> 'int' 1
// | | | | `-IntegerLiteral <col:17> 'int' 0
// | | | `-ImplicitCastExpr <col:21> 'int *' <LValueToRValue>
// | | |   `-DeclRefExpr <col:21> 'int *' lvalue Var 0x29a91a8 'i5' 'int *'
// | | `-BinaryOperator <col:26, col:32> 'long' '*'
// | |   |-ImplicitCastExpr <col:26> 'long' <IntegralCast>
// | |   | `-IntegerLiteral <col:26> 'int' 5
// | |   `-CallExpr <col:28, col:32> 'long'
// | |     `-ImplicitCastExpr <col:28> 'long (*)()' <FunctionToPointerDecay>
// | |       `-DeclRefExpr <col:28> 'long ()' Function 0x29a8470 'get' 'long ()'
// | `-CallExpr <col:36, col:40> 'long'
// |   `-ImplicitCastExpr <col:36> 'long (*)()' <FunctionToPointerDecay>
// |     `-DeclRefExpr <col:36> 'long ()' Function 0x29a8470 'get' 'long ()'
// `-BinaryOperator <col:44, col:57> 'long' '*'
//   |-ImplicitCastExpr <col:44, col:51> 'long' <IntegralCast>
//   | `-ParenExpr <col:44, col:51> 'int'
//   |   `-BinaryOperator <col:45, col:50> 'int' '+'
//   |     |-IntegerLiteral <col:45> 'int' 12
//   |     `-IntegerLiteral <col:50> 'int' 3
//   `-CallExpr <col:53, col:57> 'long'
//     `-ImplicitCastExpr <col:53> 'long (*)()' <FunctionToPointerDecay>
//       `-DeclRefExpr <col:53> 'long ()' Function 0x29a8470 'get' 'long ()'
//
// ParenExpr <col:25, col:31> 'char *'
// `-UnaryOperator <col:26, col:29> 'char *' postfix '++'
//   `-DeclRefExpr <col:26> 'char *' lvalue Var 0x3c05ae8 'pos' 'char *'
//
func pointerParts(node *ast.Node, p *program.Program) (
	pnt ast.Node, value ast.Node, back func(), err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot pointerParts: err = %v", err)
			if (*node) != nil {
				err = fmt.Errorf("Code line: %d. %v", (*node).Position().Line, err)
			}
		}
	}()

	var counter int

	var lastNode *ast.Node

	// replacer zero
	zero := &ast.IntegerLiteral{Type: "int", Value: "0"}

	var baseTypes []*string
	var searcher func(*ast.Node) bool
	replacer := func(node *ast.Node) {
		pnt = *node
		lastNode = node
		*node = zero
		counter++
	}
	searcher = func(node *ast.Node) (modify bool) {
		// save types of all nodes
		t, ok := ast.GetTypeIfExist(*node)
		if !ok {
			panic(fmt.Errorf("Not support parent type %T in pointer seaching", node))
		}
		baseTypes = append(baseTypes, t)

		// find
		if types.IsCPointer(*t, p) || types.IsCArray(*t, p) {
			switch (*node).(type) {
			case *ast.BinaryOperator,
				*ast.ImplicitCastExpr,
				*ast.ParenExpr:
				// go deeper
			default:
				return true
			}
		} else {
			// type is not pointer
			switch (*node).(type) {
			case *ast.CallExpr,
				*ast.ArraySubscriptExpr,
				*ast.MemberExpr,
				*ast.CStyleCastExpr:
				return
			}
		}
		switch (*node).(type) {
		case *ast.UnaryOperator:
			baseTypes = baseTypes[:len(baseTypes)-1]
			return
		}
		for i := range (*node).Children() {
			if searcher(&((*node).Children()[i])) {
				replacer(&((*node).Children()[i]))
			}
		}
		return false
	}
	if searcher(node) {
		pnt = *node
		lastNode = node
		*node = zero
		counter++
	}

	if counter != 1 {
		err = fmt.Errorf("counter is not 1: %d", counter)
		return
	}
	if pnt == nil {
		err = fmt.Errorf("pointer is nil")
		return
	}

	copyTypes := make([]string, len(baseTypes))
	for i := range baseTypes {
		copyTypes[i] = *(baseTypes[i])
	}
	back = func() {
		// return back types
		for i := range baseTypes {
			*(baseTypes[i]) = copyTypes[i]
		}
		// return back node
		*lastNode = pnt
	}

	// replace all types to `int`
	for i := range baseTypes {
		*baseTypes[i] = `int`
	}

	value = *node

	return
}

// transpilePointerArith - transpile pointer aripthmetic
// Example of using:
// *(t + 1) = ...
func transpilePointerArith(n *ast.UnaryOperator, p *program.Program) (
	expr goast.Expr, eType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpilePointerArith: err = %v", err)
		}
	}()

	if n.Operator != "*" {
		err = fmt.Errorf("not valid operator : %s", n.Operator)
		return
	}

	var pnt, value ast.Node
	var back func()
	pnt, value, back, err = pointerParts(&(n.Children()[0]), p)
	if err != nil {
		return
	}

	e, eType, newPre, newPost, err := atomicOperation(value, p)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
	eType = n.Type

	// return all types
	back()

	arr, arrType, newPre, newPost, err := atomicOperation(pnt, p)
	if err != nil {
		return
	}
	_ = arrType
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	switch v := pnt.(type) {
	case *ast.MemberExpr:
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
		return &goast.IndexExpr{
			X: &goast.ParenExpr{
				X:      arr,
				Lparen: 1,
			},
			Index: e,
		}, eType, preStmts, postStmts, err

	case *ast.UnaryOperator:
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
	return nil, "", nil, nil, fmt.Errorf("Cannot found : %#v", pnt)
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

	// Example:
	// UnaryOperator 'int' prefix '-'
	// `-ImplicitCastExpr 'int' <LValueToRValue>
	//   `-DeclRefExpr 'int' lvalue Var 0x3b42898 'c' 'int'

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
		if typ, ok := ast.GetTypeIfExist(n.Children()[0]); ok {
			t = *typ
		} else {
			panic(fmt.Sprintf("cannot find first child from: %#v", n.Children()[0]))
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
