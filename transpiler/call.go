// This file contains functions for transpiling function calls (invocations).

package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
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

const undefineFunctionName string = "C4GO_UNDEFINE_NAME"

func getName(p *program.Program, firstChild ast.Node) (name string, err error) {
	switch fc := firstChild.(type) {
	case *ast.DeclRefExpr:
		return fc.Name, nil

	case *ast.MemberExpr:
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

	case *ast.CallExpr:
		if len(fc.Children()) == 0 {
			return undefineFunctionName, nil
		}
		return getName(p, fc.Children()[0])

	case *ast.ParenExpr:
		if len(fc.Children()) == 0 {
			return undefineFunctionName, nil
		}
		return getName(p, fc.Children()[0])

	case *ast.UnaryOperator:
		if len(fc.Children()) == 0 {
			return undefineFunctionName, nil
		}
		return getName(p, fc.Children()[0])

	case *ast.ImplicitCastExpr:
		if len(fc.Children()) == 0 {
			return undefineFunctionName, nil
		}
		return getName(p, fc.Children()[0])

	case *ast.CStyleCastExpr:
		if len(fc.Children()) == 0 {
			return undefineFunctionName, nil
		}
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
		if resultType == "" {
			resultType = n.Type
		}
	}()

	functionName, err := getName(p, n)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil
		functionName = undefineFunctionName
	}
	functionName = util.ConvertFunctionNameFromCtoGo(functionName)

	defer func() {
		if err != nil {
			err = fmt.Errorf("name of call function is %v. %v",
				functionName, err)
		}
	}()

	// specific for va_list
	changeVaListFuncs(&functionName)

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
			unary, expression, back, err := findAndReplaceUnaryExprOrTypeTraitExpr(&n.Children()[1])
			if err != nil {
				back()
				return transpileCallExprCalloc(n.Children()[1],
					&ast.UnaryExprOrTypeTraitExpr{
						Function: "sizeof",
						Type1:    "unsigned long",
						Type2:    "char",
					}, p)
			}
			return transpileCallExprCalloc(expression, unary.(*ast.UnaryExprOrTypeTraitExpr), p)
		}
	}

	if p.IncludeHeaderIsExists("stdlib.h") && functionName == "realloc" {
		p.IsHaveRealloc = true
		p.IsHaveMemcpy = true
		p.SetCalled(functionName)
	}

	if p.IncludeHeaderIsExists("string.h") && functionName == "memcpy" {
		p.IsHaveMemcpy = true
		p.SetCalled(functionName)
	}

	// function "calloc" from stdlib.h
	if p.IncludeHeaderIsExists("stdlib.h") {
		if functionName == "calloc" && len(n.Children()) == 3 {
			if unary, ok := n.Children()[2].(*ast.UnaryExprOrTypeTraitExpr); ok {
				return transpileCallExprCalloc(n.Children()[1], unary, p)
			}

			call := &ast.CallExpr{}

			call.AddChild(&ast.ImplicitCastExpr{
				Type: "void *(*)(unsigned long)",
			})
			call.ChildNodes[0].(*ast.ImplicitCastExpr).AddChild(&ast.DeclRefExpr{
				Type: "void *(unsigned long)",
				Name: "malloc",
			})

			bin := &ast.BinaryOperator{
				Operator: "*",
				Type:     "unsigned long",
			}
			bin.AddChild(n.ChildNodes[1])
			bin.AddChild(n.ChildNodes[2])
			call.AddChild(bin)

			return transpileCallExpr(call, p)
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

	if functionDef != nil {
		p.SetCalled(functionName)
	}

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

			checker := func(t string) bool {
				return util.IsFunction(t) || types.IsTypedefFunction(p, t)
			}

			var finder func(n ast.Node) string
			finder = func(n ast.Node) (t string) {
				switch v := n.(type) {
				case *ast.ImplicitCastExpr:
					t = v.Type
				case *ast.ParenExpr:
					t = v.Type
				case *ast.CStyleCastExpr:
					t = v.Type
				default:
					panic(fmt.Errorf("add type %T", n))
				}
				if checker(t) {
					return t
				}
				if len(n.Children()) == 0 {
					return ""
				}
				return finder(n.Children()[0])
			}

			if t := finder(n.Children()[0]); checker(t) {
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
				if !types.IsPointer(functionDef.ArgumentTypes[i], p) {
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
		return nil, n.Type, preStmts, postStmts, nil
	}

	return util.NewCallExpr(functionName, realArgs...),
		functionDef.ReturnType, preStmts, postStmts, nil
}

func findAndReplaceUnaryExprOrTypeTraitExpr(node *ast.Node) (
	unary ast.Node, tree ast.Node, back func(), err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot findAndReplaceUnaryExprOrTypeTraitExpr: err = %v", err)
			if (*node) != nil {
				err = fmt.Errorf("Code line: %d. %v", (*node).Position().Line, err)
			}
		}
	}()

	var counter int
	var lastNode *ast.Node
	one := &ast.IntegerLiteral{Type: "int", Value: "1"}

	var searcher func(*ast.Node) bool
	replacer := func(node *ast.Node) {
		unary = *node
		lastNode = node
		*node = one
		counter++
	}
	searcher = func(node *ast.Node) (modify bool) {
		// find
		if u, ok := (*node).(*ast.UnaryExprOrTypeTraitExpr); ok &&
			u.Function == "sizeof" {
			return true
		}
		for i := range (*node).Children() {
			if searcher(&((*node).Children()[i])) {
				replacer(&((*node).Children()[i]))
			}
		}
		return false
	}
	if searcher(node) {
		unary = *node
		lastNode = node
		*node = one
		counter++
	}

	back = func() {
		// return back node
		if unary != nil {
			*lastNode = unary
		}
	}

	if counter != 1 {
		err = fmt.Errorf("counter is not 1: %d", counter)
		return
	}
	if unary == nil {
		err = fmt.Errorf("pointer is nil")
		return
	}

	tree = *node

	return
}

//
// calloc nodes:
// [0] - function identification
// [1] - expression
// [2] - type UnaryExprOrTypeTraitExpr always
func transpileCallExprCalloc(expression ast.Node, unary *ast.UnaryExprOrTypeTraitExpr, p *program.Program) (
	expr *goast.CallExpr, resultType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Function: calloc. err = %v", err)
		}
	}()

	size, _, newPre, newPost, err := atomicOperation(expression, p)
	if err != nil {
		return
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	var t string = unary.Type2
	if t == "" {
		_, t, _, _, _ = transpileToExpr(unary, p, false)
	}
	resultType = t + "*"
	t, err = types.ResolveType(p, t)
	if err != nil {
		return nil, "", nil, nil, err
	}
	goType := &goast.ArrayType{Elt: goast.NewIdent(t)}

	return util.NewCallExpr("make", goType, size),
		resultType, preStmts, postStmts, nil
}

func transpileCallExprQsort(n *ast.CallExpr, p *program.Program) (
	expr *goast.CallExpr, resultType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Function: qsort. err = %v", err)
		}
		if resultType == "" {
			resultType = n.Type
		}
	}()
	// CallExpr 0x2c6b1b0 'void'
	// |-ImplicitCastExpr 'void (*)(void *, size_t, size_t, __compar_fn_t)' <FunctionToPointerDecay>
	// | `-DeclRefExpr 'void (void *, size_t, size_t, __compar_fn_t)' Function 0x2bec110 'qsort' 'void (void *, size_t, size_t, __compar_fn_t)'
	// |-ImplicitCastExpr 'void *' <BitCast>
	// | `-ImplicitCastExpr 'int *' <ArrayToPointerDecay>
	// |   `-DeclRefExpr 'int [6]' lvalue Var 0x2c6a6c0 'values' 'int [6]'
	// |-ImplicitCastExpr 'size_t':'unsigned long' <IntegralCast>
	// | `-IntegerLiteral 'int' 6
	// |-UnaryExprOrTypeTraitExpr 'unsigned long' sizeof 'int'
	// `-ImplicitCastExpr 'int (*)(const void *, const void *)' <FunctionToPointerDecay>
	//   `-DeclRefExpr 'int (const void *, const void *)' Function 0x2c6aa70 'compare' 'int (const void *, const void *)'
	//
	// CallExpr  'void'
	// |-ImplicitCastExpr 'void (*)(void *, size_t, size_t, __compar_fn_t)' <FunctionToPointerDecay>
	// | `-DeclRefExpr 'void (void *, size_t, size_t, __compar_fn_t)' Function 0x361b6d0 'qsort' 'void (void *, size_t, size_t, __compar_fn_t)'
	// |-ImplicitCastExpr 'void *' <BitCast>
	// | `-ImplicitCastExpr 'int *' <LValueToRValue>
	// |   `-DeclRefExpr 'int *' lvalue ParmVar 0x3668088 'id' 'int *'
	// |-ImplicitCastExpr 'size_t':'unsigned long' <IntegralCast>
	// | `-ImplicitCastExpr 'int' <LValueToRValue>
	// |   `-DeclRefExpr 'int' lvalue Var 0x36684e0 'nid' 'int'
	// |-UnaryExprOrTypeTraitExpr 'unsigned long' sizeof
	// | `-ParenExpr 'int' lvalue
	// |   `-ArraySubscriptExpr 'int' lvalue
	// |     |-ImplicitCastExpr 'int *' <LValueToRValue>
	// |     | `-DeclRefExpr  'int *' lvalue ParmVar 0x3668088 'id' 'int *'
	// |     `-IntegerLiteral 'int' 0
	// `-ImplicitCastExpr '__compar_fn_t':'int (*)(const void *, const void *)' <BitCast>
	//   `-CStyleCastExpr 'void *' <BitCast>
	//     `-ImplicitCastExpr 'int (*)(void *, void *)' <FunctionToPointerDecay>
	//       `-DeclRefExpr 'int (void *, void *)' Function 0x3665148 'intcmp' 'int (void *, void *)'
	//

	arr, _, newPre, newPost, err := atomicOperation(n.Children()[1], p)
	if err != nil {
		err = fmt.Errorf("cannot transpile array node: %v", err)
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	t, ok := ast.GetTypeIfExist(n.Children()[1].Children()[0])
	if !ok {
		err = fmt.Errorf("cannot take type array node")
		return nil, "", nil, nil, err
	}

	*t = strings.Replace(*t, "*", "", 1)

	arrType, err := types.ResolveType(p, *t)
	if err != nil {
		err = fmt.Errorf("cannot resolve array type: %v", err)
		return nil, "", nil, nil, err
	}

	size, sizeType, newPre, newPost, err := atomicOperation(n.Children()[2], p)
	if err != nil {
		err = fmt.Errorf("cannot transpile size node: %v", err)
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	size, err = types.CastExpr(p, size, sizeType, "int")
	if err != nil {
		err = fmt.Errorf("cannot cast size node to int : %v", err)
		return nil, "", nil, nil, err
	}

	f, _, newPre, newPost, err := atomicOperation(n.Children()[4], p)
	if err != nil {
		err = fmt.Errorf("cannot transpile function node: %v", err)
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// cast from `int (void *, void *)` to `bool (int, int)`

	valA := CreateSliceFromReference(arrType, &goast.IndexExpr{
		X:      arr,
		Lbrack: 1,
		Index:  goast.NewIdent("a"),
	})
	valB := CreateSliceFromReference(arrType, &goast.IndexExpr{
		X:      arr,
		Lbrack: 1,
		Index:  goast.NewIdent("b"),
	})
	f = &goast.FuncLit{
		Type: &goast.FuncType{
			Params: &goast.FieldList{
				List: []*goast.Field{
					{
						Names: []*goast.Ident{goast.NewIdent("a"), goast.NewIdent("b")},
						Type:  goast.NewIdent("int"),
					},
				},
			},
			Results: &goast.FieldList{
				List: []*goast.Field{
					{Type: goast.NewIdent("bool")},
				},
			},
		},
		Body: &goast.BlockStmt{
			List: []goast.Stmt{
				&goast.ReturnStmt{
					Results: []goast.Expr{
						&goast.BinaryExpr{
							X: &goast.CallExpr{
								Fun: f,
								Args: []goast.Expr{
									valA,
									valB,
								},
							},
							Op: token.LEQ, // <=
							Y:  goast.NewIdent("0"),
						},
					},
				},
			},
		},
	}

	p.AddImport("sort")

	return util.NewCallExpr("sort.SliceStable",
		&goast.SliceExpr{
			X:    arr,
			High: size,
		},
		f,
	), "", preStmts, postStmts, nil

}
