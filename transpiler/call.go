// This file contains functions for transpiling function calls (invocations).

package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

func getMemberName(firstChild ast.Node) (name string, ok bool) {
	switch fc := firstChild.(type) {
	case *ast.MemberExpr:
		return fc.Name, true

	case *ast.ParenExpr:
		return getMemberName(fc.Children()[0])

	case *ast.ImplicitCastExpr:
		return getMemberName(fc.Children()[0])

	case *ast.CStyleCastExpr:
		return getMemberName(fc.Children()[0])

	}
	return "", false
}

func getName(p *program.Program, firstChild ast.Node) (name string, err error) {
	switch fc := firstChild.(type) {
	case *ast.DeclRefExpr:
		return fc.Name, nil

	case *ast.MemberExpr:
		if isUnionMemberExpr(p, fc) {
			var expr goast.Expr
			expr, _, _, _, err = transpileToExpr(fc, p, false)
			if err != nil {
				return
			}
			var buf bytes.Buffer
			err = printer.Fprint(&buf, token.NewFileSet(), expr)
			if err != nil {
				return
			}
			return buf.String(), nil
		}
		if len(fc.Children()) == 0 {
			return fc.Name, nil
		}
		var n string
		n, err = getName(p, fc.Children()[0])
		if err != nil {
			return
		}
		return n + "." + fc.Name, nil

	case *ast.ParenExpr:
		return getName(p, fc.Children()[0])

	case *ast.UnaryOperator:
		return getName(p, fc.Children()[0])

	case *ast.ImplicitCastExpr:
		return getName(p, fc.Children()[0])

	case *ast.CStyleCastExpr:
		return getName(p, fc.Children()[0])

	case *ast.ArraySubscriptExpr:
		var expr goast.Expr
		expr, _, _, _, err = transpileArraySubscriptExpr(fc, p)
		if err != nil {
			return
		}
		var buf bytes.Buffer
		err = printer.Fprint(&buf, token.NewFileSet(), expr)
		if err != nil {
			return
		}
		return buf.String(), nil
	}

	return "", fmt.Errorf("cannot getName for: %#v", firstChild)
}

func getNameOfFunctionFromCallExpr(p *program.Program, n *ast.CallExpr) (string, error) {
	// The first child will always contain the name of the function being
	// called.
	firstChild, ok := n.Children()[0].(*ast.ImplicitCastExpr)
	if !ok {
		err := fmt.Errorf("unable to use CallExpr: %#v", n.Children()[0])
		return "", err
	}

	return getName(p, firstChild.Children()[0])
}

// simplificationCallExprPrintf - minimaze Go code
// transpile C code : printf("Hello")
// to Go code       : fmt_Printf("Hello")
// AST example :
// CallExpr <> 'int'
// |-ImplicitCastExpr <> 'int (*)(const char *, ...)' <FunctionToPointerDecay>
// | `-DeclRefExpr <> 'int (const char *, ...)' Function 0x2fec178 'printf' 'int (const char *, ...)'
// `-ImplicitCastExpr <> 'const char *' <BitCast>
//   `-ImplicitCastExpr <> 'char *' <ArrayToPointerDecay>
//     `-StringLiteral <> 'char [6]' lvalue "Hello"
func simplificationCallExprPrintf(call *ast.CallExpr, p *program.Program) (
	expr *goast.CallExpr, ok bool) {

	var isPrintfCode bool
	var printfText string
	if call.Type == "int" && len(call.ChildNodes) == 2 {
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

	if !isPrintfCode {
		return
	}

	// 0: *ast.ExprStmt {
	// .  X: *ast.CallExpr {
	// .  .  Fun: *ast.SelectorExpr {
	// .  .  .  X: *ast.Ident {
	// .  .  .  .  NamePos: 8:2
	// .  .  .  .  Name: "fmt"
	// .  .  .  }
	// .  .  .  Sel: *ast.Ident {
	// .  .  .  .  NamePos: 8:6
	// .  .  .  .  Name: "Printf"
	// .  .  .  }
	// .  .  }
	// .  .  Lparen: 8:12
	// .  .  Args: []ast.Expr (len = 1) {
	// .  .  .  0: *ast.BasicLit {
	// .  .  .  .  ValuePos: 8:13
	// .  .  .  .  Kind: STRING
	// .  .  .  .  Value: "\"Hello, Golang\\n\""
	// .  .  .  }
	// .  .  }
	// .  .  Ellipsis: -
	// .  .  Rparen: 8:30
	// .  }
	// }
	p.AddImport("fmt")
	printfText = strconv.Quote(printfText)
	return util.NewCallExpr("fmt"+"."+"Printf",
		&goast.BasicLit{
			Kind:  token.STRING,
			Value: printfText,
		}), true
}

// transpileCallExpr transpiles expressions that calls a function, for example:
//
//     foo("bar")
//
// It returns three arguments; the Go AST expression, the C type (that is
// returned by the function) and any error. If there is an error returned you
// can assume the first two arguments will not contain any useful information.
func transpileCallExpr(n *ast.CallExpr, p *program.Program) (
	expr *goast.CallExpr, resultType string,
	preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Error in transpileCallExpr : %v", err)
		}
	}()

	functionName, err := getNameOfFunctionFromCallExpr(p, n)
	if err != nil {
		return nil, "", nil, nil, err
	}
	functionName = util.ConvertFunctionNameFromCtoGo(functionName)

	defer func() {
		if err != nil {
			err = fmt.Errorf("name of call function is %v. %v",
				functionName, err)
		}
	}()

	if functionName == "__builtin_va_start" ||
		functionName == "__builtin_va_end" {
		// ignore function __builtin_va_start, __builtin_va_end
		// see "Variadic functions"
		return nil, "", nil, nil, nil
	}

	// function "malloc" from stdlib.h
	//
	// Change from "malloc" to "calloc"
	//
	// CallExpr <> 'void *'
	// |-ImplicitCastExpr <> 'void *(*)(unsigned long)' <FunctionToPointerDecay>
	// | `-DeclRefExpr <> 'void *(unsigned long)' Function 'malloc' 'void *(unsigned long)'
	// `-ImplicitCastExpr <> 'unsigned long' <IntegralCast>
	//   `- ...
	//
	// CallExpr <> 'void *'
	// |-ImplicitCastExpr <> 'void *(*)(unsigned long, unsigned long)' <FunctionToPointerDecay>
	// | `-DeclRefExpr <> 'void *(unsigned long, unsigned long)' Function 'calloc' 'void *(unsigned long, unsigned long)'
	// |-ImplicitCastExpr <> 'unsigned long' <IntegralCast>
	// | `- ...
	// `-UnaryExprOrTypeTraitExpr <> 'unsigned long' sizeof 'char'
	if p.IncludeHeaderIsExists("stdlib.h") {
		if functionName == "malloc" && len(n.Children()) == 2 {
			// Change from "malloc" to "calloc"
			n.ChildNodes[0] = &ast.ImplicitCastExpr{
				Type: "void *(*)(unsigned long, unsigned long)",
			}
			n.ChildNodes[0].(*ast.ImplicitCastExpr).AddChild(&ast.DeclRefExpr{
				Type: "void *(unsigned long, unsigned long)",
				Name: "calloc",
			})
			n.AddChild(&ast.UnaryExprOrTypeTraitExpr{
				Type1:    "unsigned long",
				Function: "sizeof",
				Type2:    "char",
			})
			return transpileCallExprCalloc(n, p)
		}
	}

	// function "realloc" from stdlib.h
	//
	// Change from "realloc" to "calloc"
	//
	// CallExpr <> 'void *'
	// |-ImplicitCastExpr <> 'void *(*)(void *, unsigned long)' <FunctionToPointerDecay>
	// | `-DeclRefExpr <> 'void *(void *, unsigned long)' Function 0x2c7e3e0 'realloc' 'void *(void *, unsigned long)'
	// |-ImplicitCastExpr <> 'void *' <BitCast>
	// | `-CStyleCastExpr <> 'char *' <BitCast>
	// |   `-ParenExpr <> 'void *'
	// |     `-ParenExpr <> 'void *'
	// |       `-CStyleCastExpr <> 'void *' <NullToPointer>
	// |         `-IntegerLiteral <> 'int' 0
	// `-ImplicitCastExpr <> 'unsigned long' <IntegralCast>
	//   `-BinaryOperator <> 'int' '*'
	//     |-ImplicitCastExpr <> 'int' <LValueToRValue>
	//     | `-DeclRefExpr <> 'int' lvalue Var 0x2ca14e8 'size' 'int'
	//     `-ImplicitCastExpr <> 'int' <LValueToRValue>
	//       `-DeclRefExpr <> 'int' lvalue Var 0x2ca14e8 'size' 'int'
	//
	// CallExpr <> 'void *'
	// |-ImplicitCastExpr <> 'void *(*)(unsigned long, unsigned long)' <FunctionToPointerDecay>
	// | `-DeclRefExpr <> 'void *(unsigned long, unsigned long)' Function 'calloc' 'void *(unsigned long, unsigned long)'
	// |-ImplicitCastExpr <> 'unsigned long' <IntegralCast>
	// | `- ...
	// `-UnaryExprOrTypeTraitExpr <> 'unsigned long' sizeof 'char'
	//
	// function "realloc" from stdlib.h
	if p.IncludeHeaderIsExists("stdlib.h") {
		if functionName == "realloc" && len(n.Children()) == 3 {
			n.ChildNodes[0] = &ast.ImplicitCastExpr{
				Type: "void *(*)(unsigned long, unsigned long)",
			}
			n.ChildNodes[0].(*ast.ImplicitCastExpr).AddChild(&ast.DeclRefExpr{
				Type: "void *(unsigned long, unsigned long)",
				Name: "calloc",
			})
			sizeofType := "char"
			if impl, ok := n.Children()[1].(*ast.ImplicitCastExpr); ok {
				switch v := impl.Children()[0].(type) {
				case *ast.ImplicitCastExpr:
					sizeofType = v.Type
				case *ast.CStyleCastExpr:
					sizeofType = v.Type
				}
			}
			sizeofType = types.GetBaseType(sizeofType)
			n.AddChild(&ast.UnaryExprOrTypeTraitExpr{
				Type1:    "unsigned long",
				Function: "sizeof",
				Type2:    sizeofType,
			})
			n.ChildNodes = []ast.Node{n.ChildNodes[0], n.ChildNodes[2], n.ChildNodes[3]}
			return transpileCallExprCalloc(n, p)
		}
	}

	// function "calloc" from stdlib.h
	if p.IncludeHeaderIsExists("stdlib.h") {
		if functionName == "calloc" && len(n.Children()) == 3 {
			return transpileCallExprCalloc(n, p)
		}
	}

	// function "qsort" from stdlib.h
	if p.IncludeHeaderIsExists("stdlib.h") {
		if functionName == "qsort" && len(n.Children()) == 5 {
			return transpileCallExprQsort(n, p)
		}
	}

	// function "printf" from stdio.h simplification
	if p.IncludeHeaderIsExists("stdio.h") {
		if functionName == "printf" && len(n.Children()) == 2 {
			if e, ok := simplificationCallExprPrintf(n, p); ok {
				return e, "int", nil, nil, nil
			}
		}
	}

	// Get the function definition from it's name. The case where it is not
	// defined is handled below (we haven't seen the prototype yet).
	functionDef := p.GetFunctionDefinition(functionName)

	if functionDef == nil {
		// We do not have a prototype for the function, but we should not exit
		// here. Instead we will create a mock definition for it so that this
		// transpile function will always return something and continue.
		//
		// The mock function definition is never actually saved to the program
		// definitions, so each time we see the CallExpr it will run this every
		// time. This is so if we come across the real prototype later it will
		// be handled correctly. Or at least "more" correctly.
		functionDef = &program.DefinitionFunction{
			Name: functionName,
		}
		if len(n.Children()) > 0 {
			if v, ok := n.Children()[0].(*ast.ImplicitCastExpr); ok &&
				(util.IsFunction(v.Type) || types.IsTypedefFunction(p, v.Type)) {
				t := v.Type
				if v, ok := p.TypedefType[t]; ok {
					t = v
				} else {
					if types.IsTypedefFunction(p, t) {
						t = t[0 : len(t)-len(" *")]
						t = p.TypedefType[t]
					}
				}
				prefix, _, fields, returns, err := util.ParseFunction(t)
				if err != nil {
					p.AddMessage(p.GenerateWarningMessage(fmt.Errorf(
						"Cannot resolve function : %v", err), n))
					return nil, "", nil, nil, err
				}
				if len(prefix) != 0 {
					p.AddMessage(p.GenerateWarningMessage(fmt.Errorf(
						"prefix `%v` is not used in type : %v",
						prefix, t), n))
				}
				functionDef.ReturnType = returns[0]
				functionDef.ArgumentTypes = fields
			}
		}
	} else {
		// type correction for definition function in
		// package program
		var ok bool
		for pos, arg := range n.Children() {
			if pos == 0 {
				continue
			}
			if pos >= len(functionDef.ArgumentTypes) {
				continue
			}
			if arg, ok = arg.(*ast.ImplicitCastExpr); ok {
				arg.(*ast.ImplicitCastExpr).Type = functionDef.ArgumentTypes[pos-1]
			}
		}
	}

	if functionDef.Substitution != "" {
		parts := strings.Split(functionDef.Substitution, ".")
		importName := strings.Join(parts[:len(parts)-1], ".")
		p.AddImport(importName)

		parts2 := strings.Split(functionDef.Substitution, "/")
		functionName = parts2[len(parts2)-1]
	}

	args := []goast.Expr{}
	argTypes := []string{}
	i := 0
	for _, arg := range n.Children()[1:] {
		if bin, ok := arg.(*ast.BinaryOperator); ok && bin.Operator == "=" {
			// example :
			// from :
			// call(val = 43);
			// to:
			// call(val = 43,val);
			var b ast.BinaryOperator
			b.Type = bin.Type
			b.Operator = ","
			b.AddChild(arg)
			b.AddChild(bin.Children()[0])
			arg = &b
		}
		if cmp, ok := arg.(*ast.CompoundAssignOperator); ok {
			// example :
			// from :
			// call(val += 43);
			// to:
			// call(val += 43,val);
			var b ast.BinaryOperator
			b.Type = cmp.Type
			b.Operator = ","
			b.AddChild(arg)
			b.AddChild(cmp.Children()[0])
			arg = &b
		}
		e, eType, newPre, newPost, err := atomicOperation(arg, p)
		if err != nil {
			err = fmt.Errorf("argument position is %d. %v", i, err)
			p.AddMessage(p.GenerateWarningMessage(err, arg))
			return nil, "unknown2", nil, nil, err
		}
		argTypes = append(argTypes, eType)
		preStmts, postStmts = combinePreAndPostStmts(
			preStmts, postStmts, newPre, newPost)

		args = append(args, e)
		i++
	}

	// These are the arguments once any transformations have taken place.
	realArgs := []goast.Expr{}

	// Apply transformation if needed. A transformation rearranges the return
	// value(s) and parameters. It is also used to indicate when a variable must
	// be passed by reference.
	if functionDef.ReturnParameters != nil || functionDef.Parameters != nil {
		for i, a := range functionDef.Parameters {
			byReference := false

			// Negative position means that it must be passed by reference.
			if a < 0 {
				byReference = true
				a = -a
			}

			// Rearrange the arguments. The -1 is because 0 would be the return
			// value.
			realArg := args[a-1]

			if byReference {
				// We have to create a temporary variable to pass by reference.
				// Then we can assign the real variable from it.
				realArg = &goast.UnaryExpr{
					Op: token.AND,
					X:  args[i],
				}
			} else {
				realArg, err = types.CastExpr(p, realArg, argTypes[i],
					functionDef.ArgumentTypes[i])
				p.AddMessage(p.GenerateWarningMessage(err, n))

				if realArg == nil {
					realArg = util.NewNil()
				}
			}

			if realArg == nil {
				return nil, "", preStmts, postStmts,
					fmt.Errorf("Real argument is nil in function : %s", functionName)
			}

			realArgs = append(realArgs, realArg)
		}
	} else {
		// Keep all the arguments the same. But make sure we cast to the correct
		// types.
		// Example of functionDef.ArgumentTypes :
		// [void *, int]
		// [char *, char * , ...]
		// Example of args:
		// [void *, int]
		// [char *, char *, char *, int, double]
		//
		for i, a := range args {
			realType := "unknownType"
			if i < len(functionDef.ArgumentTypes) {
				if len(functionDef.ArgumentTypes) > 1 &&
					i >= len(functionDef.ArgumentTypes)-1 &&
					functionDef.ArgumentTypes[len(functionDef.ArgumentTypes)-1] == "..." {
					realType = functionDef.ArgumentTypes[len(functionDef.ArgumentTypes)-2]
				} else {
					if len(functionDef.ArgumentTypes) > 0 {
						if len(functionDef.ArgumentTypes[i]) != 0 {
							realType = functionDef.ArgumentTypes[i]
							if strings.TrimSpace(realType) != "void" {
								a, err = types.CastExpr(p, a, argTypes[i], realType)

								if p.AddMessage(p.GenerateWarningMessage(err, n)) {
									a = util.NewNil()
								}
							}
						}
					}
				}
			}

			if strings.Contains(realType, "...") {
				p.AddMessage(p.GenerateWarningMessage(
					fmt.Errorf("not acceptable type '...'"), n))
			}

			if a == nil {
				return nil, "", preStmts, postStmts,
					fmt.Errorf("Argument is nil in function : %s", functionName)
			}

			if len(functionDef.ArgumentTypes) > i {
				if !util.IsPointer(functionDef.ArgumentTypes[i]) {
					if strings.HasPrefix(functionDef.ArgumentTypes[i], "union ") {
						a = &goast.CallExpr{
							Fun: &goast.SelectorExpr{
								X:   a,
								Sel: goast.NewIdent("copy"),
							},
							Lparen: 1,
						}
					}
				}
			}

			realArgs = append(realArgs, a)
		}
	}

	// Added for support removing function `free` of <stdlib.h>
	// Example of C code:
	// free(i+=4,buffer)
	// Example of result Go code:
	// i += 4
	// _ = buffer
	if functionDef.Substitution == "_" {
		devNull := &goast.AssignStmt{
			Lhs: []goast.Expr{goast.NewIdent("_")},
			Tok: token.ASSIGN,
			Rhs: []goast.Expr{realArgs[0]},
		}
		preStmts = append(preStmts, devNull)
		return nil, "", preStmts, postStmts, nil
	}

	return util.NewCallExpr(functionName, realArgs...),
		functionDef.ReturnType, preStmts, postStmts, nil
}

func transpileCallExprCalloc(n *ast.CallExpr, p *program.Program) (
	expr *goast.CallExpr, resultType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Function: calloc. err = %v", err)
		}
	}()
	size, sizeType, preStmts, postStmts, err := transpileToExpr(n.Children()[1], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}
	size, err = types.CastExpr(p, size, sizeType, "unsigned long")
	if err != nil {
		return nil, "", nil, nil, err
	}

	//
	// -UnaryExprOrTypeTraitExpr 0x24e0498 <col:27, col:39> 'unsigned long' sizeof 'float'
	//
	// OR
	//
	// t.w = (float*) calloc(t.nw, sizeof(*t.w));
	// -UnaryExprOrTypeTraitExpr <> 'unsigned long' sizeof
	//  `-ParenExpr <> 'float' lvalue
	//    `-UnaryOperator <> 'float' lvalue prefix '*'
	//      `-ImplicitCastExpr <> 'float *' <LValueToRValue>
	//        `-MemberExpr <> 'float *' lvalue .w 0x2548b38
	//          `-DeclRefExpr <> 'struct cws':'struct cws' lvalue Var 0x2548c40 't' 'struct cws': 'struct cws'
	//
	// OR
	//
	// CallExpr <> 'void *'
	// |-ImplicitCastExpr <> 'void *(*)(unsigned long, unsigned long)' <FunctionToPointerDecay>
	// | `-DeclRefExpr <> 'void *(unsigned long, unsigned long)' Function 0x3b9b048 'calloc' 'void *(unsigned long, unsigned long)'
	// |-ImplicitCastExpr <> 'unsigned long' <IntegralCast>
	// | `-ImplicitCastExpr <> 'int' <LValueToRValue>
	// |   `-DeclRefExpr <> 'int' lvalue ParmVar 0x3bc40d8 'n' 'int'
	// `-ImplicitCastExpr <> 'unsigned long' <IntegralCast>
	//   `-ImplicitCastExpr <> 'int' <LValueToRValue>
	//     `-DeclRefExpr <> 'int' lvalue Var 0x3bc4268 'sizeT' 'int'
	var goType goast.Expr
	switch v := n.Children()[2].(type) {
	case *ast.UnaryExprOrTypeTraitExpr:
		var t string
		if v.Type2 == "" {
			_, t, _, _, _ = transpileToExpr(v.Children()[0], p, false)
		} else {
			t = v.Type2
		}
		t, err := types.ResolveType(p, t)
		if err != nil {
			return nil, "", nil, nil, err
		}
		goType = &goast.ArrayType{Elt: goast.NewIdent(t)}
	default:
		var goTypeT string
		goType, goTypeT, _, _, _ = transpileToExpr(v.Children()[0], p, false)
		size, err = types.CastExpr(p, size, sizeType, goTypeT)
		if err != nil {
			return nil, "", nil, nil, err
		}
		size = &goast.BinaryExpr{
			X:  &goast.ParenExpr{X: goType},
			Op: token.MUL,
			Y:  &goast.ParenExpr{X: size},
		}
		goType = &goast.ArrayType{
			Elt: goast.NewIdent("byte"),
		}
	}

	return util.NewCallExpr("make", goType, size),
		n.Type, preStmts, postStmts, nil
}

func transpileCallExprQsort(n *ast.CallExpr, p *program.Program) (
	expr *goast.CallExpr, resultType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Function: qsort. err = %v", err)
		}
	}()
	// CallExpr 0x2c6b1b0 <line:182:2, col:40> 'void'
	// |-ImplicitCastExpr 0x2c6b198 <col:2> 'void (*)(void *, size_t, size_t, __compar_fn_t)' <FunctionToPointerDecay>
	// | `-DeclRefExpr 0x2c6b070 <col:2> 'void (void *, size_t, size_t, __compar_fn_t)' Function 0x2bec110 'qsort' 'void (void *, size_t, size_t, __compar_fn_t)'
	// |-ImplicitCastExpr 0x2c6b210 <col:9> 'void *' <BitCast>
	// | `-ImplicitCastExpr 0x2c6b1f8 <col:9> 'int *' <ArrayToPointerDecay>
	// |   `-DeclRefExpr 0x2c6b098 <col:9> 'int [6]' lvalue Var 0x2c6a6c0 'values' 'int [6]'
	// |-ImplicitCastExpr 0x2c6b228 <col:17> 'size_t':'unsigned long' <IntegralCast>
	// | `-IntegerLiteral 0x2c6b0c0 <col:17> 'int' 6
	// |-UnaryExprOrTypeTraitExpr 0x2c6b0f8 <col:20, col:30> 'unsigned long' sizeof 'int'
	// `-ImplicitCastExpr 0x2c6b240 <col:33> 'int (*)(const void *, const void *)' <FunctionToPointerDecay>
	//   `-DeclRefExpr 0x2c6b118 <col:33> 'int (const void *, const void *)' Function 0x2c6aa70 'compare' 'int (const void *, const void *)'
	var element [4]goast.Expr
	for i := 1; i < 5; i++ {
		el, _, newPre, newPost, err := transpileToExpr(n.Children()[i], p, false)
		if err != nil {
			return nil, "", nil, nil, err
		}
		element[i-1] = el
		preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)
	}
	// found the C type
	t := n.Children()[3].(*ast.UnaryExprOrTypeTraitExpr).Type2
	t, err = types.ResolveType(p, t)
	if err != nil {
		return nil, "", nil, nil, err
	}

	var compareFunc string
	if v, ok := element[3].(*goast.Ident); ok {
		compareFunc = v.Name
	} else {
		return nil, "", nil, nil,
			fmt.Errorf("golang ast for compare function have type %T, expect ast.Ident", element[3])
	}

	var varName string
	if v, ok := element[0].(*goast.Ident); ok {
		varName = v.Name
	} else {
		return nil, "", nil, nil,
			fmt.Errorf("golang ast for variable name have type %T, expect ast.Ident", element[3])
	}

	p.AddImport("sort")
	src := fmt.Sprintf(`package main
		var %s func(a,b interface{})int
		var temp = func(i, j int) bool {
			c4goTempVarA := ([]%s{%s[i]})
			c4goTempVarB := ([]%s{%s[j]})
			return %s(c4goTempVarA, c4goTempVarB) <= 0
		}`, compareFunc, t, varName, t, varName, compareFunc)

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return nil, "", nil, nil, err
	}

	// AST tree part of code after "var temp = ..."
	convertExpr := f.Decls[1].(*goast.GenDecl).Specs[0].(*goast.ValueSpec).Values[0]

	return &goast.CallExpr{
		Fun: &goast.SelectorExpr{
			X:   goast.NewIdent("sort"),
			Sel: goast.NewIdent("SliceStable"),
		},
		Args: []goast.Expr{
			element[0],
			convertExpr,
		},
	}, "", preStmts, postStmts, nil
}
