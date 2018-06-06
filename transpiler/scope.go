// This file contains functions for transpiling scopes. A scope is zero or more
// statements between a set of curly brackets.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strconv"

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

		// minimaze Go code
		// transpile C code : printf("Hello")
		// to Go code       : fmt.Printf("Hello")
		// AST example :
		// CallExpr <> 'int'
		// |-ImplicitCastExpr <> 'int (*)(const char *, ...)' <FunctionToPointerDecay>
		// | `-DeclRefExpr <> 'int (const char *, ...)' Function 0x2fec178 'printf' 'int (const char *, ...)'
		// `-ImplicitCastExpr <> 'const char *' <BitCast>
		//   `-ImplicitCastExpr <> 'char *' <ArrayToPointerDecay>
		//     `-StringLiteral <> 'char [6]' lvalue "Hello"
		var isPrintfCode bool
		var printfText string
		if call, ok := x.(*ast.CallExpr); ok && call.Type == "int" && len(call.ChildNodes) == 2 {
			var step1 bool
			if impl, ok := call.ChildNodes[0].(*ast.ImplicitCastExpr); ok && len(impl.ChildNodes) == 1 {
				if decl, ok := impl.ChildNodes[0].(*ast.DeclRefExpr); ok && decl.Name == "printf" {
					if impl.Type == "int (*)(const char *, ...)" {
						step1 = true
					}
				}
			}
			var step2 bool
			if impl, ok := call.ChildNodes[1].(*ast.ImplicitCastExpr); ok {
				if impl.Type == "const char *" && len(impl.ChildNodes) == 1 {
					if impl2, ok := impl.ChildNodes[0].(*ast.ImplicitCastExpr); ok {
						if impl2.Type == "char *" && len(impl2.ChildNodes) == 1 {
							if str, ok := impl2.ChildNodes[0].(*ast.StringLiteral); ok {
								step2 = true
								printfText = str.Value
							}
						}
					}
				}
			}
			if step1 && step2 {
				isPrintfCode = true
			}
		}

		var result []goast.Stmt
		switch {
		case isPrintfCode:
			/*
			   0: *ast.ExprStmt {
			   .  X: *ast.CallExpr {
			   .  .  Fun: *ast.SelectorExpr {
			   .  .  .  X: *ast.Ident {
			   .  .  .  .  NamePos: 8:2
			   .  .  .  .  Name: "fmt"
			   .  .  .  }
			   .  .  .  Sel: *ast.Ident {
			   .  .  .  .  NamePos: 8:6
			   .  .  .  .  Name: "Printf"
			   .  .  .  }
			   .  .  }
			   .  .  Lparen: 8:12
			   .  .  Args: []ast.Expr (len = 1) {
			   .  .  .  0: *ast.BasicLit {
			   .  .  .  .  ValuePos: 8:13
			   .  .  .  .  Kind: STRING
			   .  .  .  .  Value: "\"Hello, Golang\\n\""
			   .  .  .  }
			   .  .  }
			   .  .  Ellipsis: -
			   .  .  Rparen: 8:30
			   .  }
			   }
			*/
			p.AddImport("fmt")
			printfText = strconv.Quote(printfText)
			result = []goast.Stmt{
				&goast.ExprStmt{
					X: &goast.CallExpr{
						Fun:    goast.NewIdent("fmt.Printf"),
						Lparen: 1,
						Args: []goast.Expr{
							&goast.BasicLit{
								Kind:  token.STRING,
								Value: printfText,
							},
						},
					},
				},
			}
		case isVaList:
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
		default:
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
