// This file contains functions for transpiling a "switch" statement.

package transpiler

import (
	goast "go/ast"
	"go/token"

	"fmt"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
)

// From :
// CaseStmt
// |-UnaryOperator 'int' prefix '-'
// | `-IntegerLiteral 'int' 1
// |-<<<NULL>>>
// `-CaseStmt
//   |-BinaryOperator 'int' '-'
//   | |- ...
//   |-<<<NULL>>>
//   `-CaseStmt
//     |- ...
//     |-<<<NULL>>>
//     `-DefaultStmt
//       `- ...
// To:
// CaseStmt
// |-UnaryOperator 'int' prefix '-'
// | `-IntegerLiteral 'int' 1
// `-<<<NULL>>>
// <<<NULL>>>
// CaseStmt
// |-BinaryOperator 'int' '-'
// | |- ...
// `-<<<NULL>>>
// <<<NULL>>>
// CaseStmt
// |- ...
// |-<<<NULL>>>
// `-DefaultStmt
//   `- ...
// <<<NULL>>>
//
// From:
// |-CaseStmt
// | `- ...
// |-NullStmt
// |-BreakStmt 0
// |-CaseStmt
// | `-...
// To:
// CaseStmt
// `- ...
// <<<NULL>>>
// NullStmt
// <<<NULL>>>
// BreakStmt 0
// <<<NULL>>>
// CaseStmt
// `-...
// <<<NULL>>>
//
// From:
// |-CaseStmt
// | `- ...
// |-CompoundAssignOperator  'int' '+=' ComputeLHSTy='int' ComputeResultTy='int'
// | `-...
// |-BreakStmt
// |-CaseStmt
// | |- ...
// To:
// `-CaseStmt
//   `- ...
// <<<NULL>>>
// CompoundAssignOperator  'int' '+=' ComputeLHSTy='int' ComputeResultTy='int'
// `-...
// <<<NULL>>>
// BreakStmt
// <<<NULL>>>
// CaseStmt
// `- ...
// <<<NULL>>>
//
func caseSplitter(nodes ...ast.Node) (cs []ast.Node) {
	if len(nodes) == 0 {
		return
	}
	if len(nodes) > 1 {
		for i := range nodes {
			cs = append(cs, nodes[i])
		}
		return
	}

	node := nodes[0]
	if node == nil {
		return
	}

	for i := range node.Children() {
		if node.Children()[i] == nil {
			continue
		}
		switch node.Children()[i].(type) {
		case *ast.CompoundAssignOperator,
			*ast.ParenExpr,
			*ast.UnaryOperator,
			*ast.CallExpr,
			*ast.ArrayFiller,
			*ast.SwitchStmt,
			*ast.ForStmt,
			*ast.IntegerLiteral,
			*ast.BinaryOperator:
			cs = append(cs, node.Children()[i])
			continue
		default:
			fmt.Printf("caseSplitter: %T\n", node.Children()[i])
			cs = append(cs, caseSplitter(node.Children()[i])...)
		}
	}

	return
}

func transpileSwitchStmt(n *ast.SwitchStmt, p *program.Program) (
	_ *goast.SwitchStmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileSwitchStmt : err = %v", err)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error - panic")
		}
	}()

	// The first two children are nil. I don't know what they are supposed to be
	// for. It looks like the number of children is also not reliable, but we
	// know that we need the last two which represent the condition and body
	// respectively.

	if len(n.Children()) < 2 {
		// I don't know what causes this condition. Need to investigate.
		panic(fmt.Sprintf("Less than two children for switch: %#v", n))
	}

	// The condition is the expression to be evaulated against each of the
	// cases.
	condition, conditionType, newPre, newPost, err := atomicOperation(
		n.Children()[len(n.Children())-2], p)
	if err != nil {
		return nil, nil, nil, err
	}
	if conditionType == "bool" {
		condition, err = types.CastExpr(p, condition, conditionType, "int")
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// separation body of switch on cases
	body, ok := n.Children()[len(n.Children())-1].(*ast.CompoundStmt)
	if !ok {
		err = fmt.Errorf("body is not ast.CompoundStmt : %T", n.Children()[len(n.Children())-1])
	}

	// CompoundStmt
	// `-CaseStmt
	//   |-UnaryOperator 'int' prefix '-'
	//   | `-IntegerLiteral 'int' 1
	//   |-<<<NULL>>>
	//   `-CaseStmt
	//     |-BinaryOperator 'int' '-'
	//     | |- ...
	//     |-<<<NULL>>>
	//     `-CaseStmt
	//       |- ...
	//       |-<<<NULL>>>
	//       `-DefaultStmt
	//         `- ...
	fmt.Println("================================================")
	fmt.Println(ast.TypesTree(body))
	fmt.Println("------------------------------------------------")
	parts := caseSplitter(body.Children()...)
	for i := range parts {
		fmt.Println(ast.TypesTree(parts[i]))
	}
	fmt.Println("================================================")

	// The body will always be a CompoundStmt because a switch statement is not
	// valid without curly brackets.
	cases, newPre, newPost, err := normalizeSwitchCases(body, p)
	if err != nil {
		return nil, nil, nil, err
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// TODO:
	//
	// For simplification switch case:
	// from:
	// case 3:
	// 	{
	// 		var c int
	// 		return
	// 	}
	// 	fallthrough
	// to:
	// case 3:
	// 	var c int
	// 	return
	//

	// Convert the normalized cases back into statements so they can be children
	// of goast.SwitchStmt.
	stmts := []goast.Stmt{}
	for _, singleCase := range cases {
		if singleCase == nil {
			panic("nil single case")
		}

		stmts = append(stmts, singleCase)
	}

	return &goast.SwitchStmt{
		Tag: condition,
		Body: &goast.BlockStmt{
			List: stmts,
		},
	}, preStmts, postStmts, nil
}

func normalizeSwitchCases(body *ast.CompoundStmt, p *program.Program) (
	_ []*goast.CaseClause, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	// The body of a switch has a non uniform structure. For example:
	//
	//     switch a {
	//     case 1:
	//         foo();
	//         bar();
	//         break;
	//     default:
	//         baz();
	//         qux();
	//     }
	//
	// Looks like:
	//
	//     *ast.CompountStmt
	//         *ast.CaseStmt     // case 1:
	//             *ast.CallExpr //     foo()
	//         *ast.CallExpr     //     bar()
	//         *ast.BreakStmt    //     break
	//         *ast.DefaultStmt  // default:
	//             *ast.CallExpr //     baz()
	//         *ast.CallExpr     //     qux()
	//
	// Each of the cases contains one child that is the first statement, but all
	// the rest are children of the parent CompountStmt. This makes it
	// especially tricky when we want to remove the 'break' or add a
	// 'fallthrough'.
	//
	// To make it easier we normalise the cases. This means that we iterate
	// through all of the statements of the CompountStmt and merge any children
	// that are not 'case' or 'break' with the previous node to give us a
	// structure like:
	//
	//     []*goast.CaseClause
	//         *goast.CaseClause      // case 1:
	//             *goast.CallExpr    //     foo()
	//             *goast.CallExpr    //     bar()
	//             // *ast.BreakStmt  //     break (was removed)
	//         *goast.CaseClause      // default:
	//             *goast.CallExpr    //     baz()
	//             *goast.CallExpr    //     qux()
	//
	// During this translation we also remove 'break' or append a 'fallthrough'.

	cases := []*goast.CaseClause{}

	for _, x := range body.Children() {
		switch c := x.(type) {
		case *ast.CaseStmt, *ast.DefaultStmt:
			var newPre, newPost []goast.Stmt
			cases, newPre, newPost, err = appendCaseOrDefaultToNormalizedCases(cases, c, p)
			if err != nil {
				return []*goast.CaseClause{}, nil, nil, err
			}

			preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
		case *ast.BreakStmt:

		default:
			var stmt goast.Stmt
			var newPre, newPost []goast.Stmt
			stmt, newPre, newPost, err = transpileToStmt(x, p)
			if err != nil {
				return []*goast.CaseClause{}, nil, nil, err
			}
			preStmts = append(preStmts, newPre...)
			preStmts = append(preStmts, stmt)
			preStmts = append(preStmts, newPost...)
		}
	}

	return cases, preStmts, postStmts, nil
}

func appendCaseOrDefaultToNormalizedCases(cases []*goast.CaseClause,
	stmt ast.Node, p *program.Program) (
	[]*goast.CaseClause, []goast.Stmt, []goast.Stmt, error) {
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	if len(cases) > 0 {
		cases[len(cases)-1].Body = append(cases[len(cases)-1].Body, &goast.BranchStmt{
			Tok: token.FALLTHROUGH,
		})
	}

	var singleCase *goast.CaseClause
	var err error
	var newPre []goast.Stmt
	var newPost []goast.Stmt

	switch c := stmt.(type) {
	case *ast.CaseStmt:
		singleCase, newPre, newPost, err = transpileCaseStmt(c, p)

	case *ast.DefaultStmt:
		singleCase, err = transpileDefaultStmt(c, p)
	}

	if singleCase != nil {
		cases = append(cases, singleCase)
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	if err != nil {
		return []*goast.CaseClause{}, nil, nil, err
	}

	return cases, preStmts, postStmts, nil
}

func transpileCaseStmt(n *ast.CaseStmt, p *program.Program) (
	_ *goast.CaseClause, _ []goast.Stmt, _ []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileCaseStmt: %v", err)
		}
	}()
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	c, cType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, nil, nil, err
	}
	if cType == "bool" {
		c, err = types.CastExpr(p, c, cType, "int")
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	stmts, err := transpileStmts(n.Children()[1:], p)
	if err != nil {
		return nil, nil, nil, err
	}

	return &goast.CaseClause{
		List: []goast.Expr{c},
		Body: stmts,
	}, preStmts, postStmts, nil
}

func transpileDefaultStmt(n *ast.DefaultStmt, p *program.Program) (
	*goast.CaseClause, error) {

	stmts, err := transpileStmts(n.Children()[0:], p)
	if err != nil {
		return nil, err
	}

	return &goast.CaseClause{
		List: nil,
		Body: stmts,
	}, nil
}
