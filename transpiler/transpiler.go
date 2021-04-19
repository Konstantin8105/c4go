// Package transpiler handles the conversion between the Clang AST and the Go
// AST.
package transpiler

import (
	"errors"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"runtime/debug"
	"strings"
	"unicode"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

var AddOutsideStruct bool

// TranspileAST iterates through the Clang AST and builds a Go AST
func TranspileAST(fileName, packageName string, withOutsideStructs bool,
	p *program.Program, root ast.Node, clangFlags []string) (
	source string, // result Go source
	err error) {
	// Start by parsing an empty file.
	p.FileSet = token.NewFileSet()
	packageSignature := fmt.Sprintf("package %v", packageName)
	f, err := parser.ParseFile(p.FileSet, fileName, packageSignature, 0)
	p.File = f
	AddOutsideStruct = withOutsideStructs

	if err != nil {
		return
	}

	// replace if type name and variable name
	{
		var replacer func(ast.Node)
		replacer = func(node ast.Node) {
			if node == nil {
				return
			}
			var vName *string
			var vType *string
			switch v := node.(type) {
			case *ast.DeclRefExpr:
				vName = &v.Name
				vType = &v.Type
			case *ast.VarDecl:
				vName = &v.Name
				vType = &v.Type
			case *ast.ParmVarDecl:
				vName = &v.Name
				vType = &v.Type
			}

			// examples:
			//   vName        vType
			//   `wb`         `wb`
			//   `wb`        `wb *`
			//   `wb`      `struct wb`
			//   `wb`      `struct wb *`
			//   `wb`      `struct wb*`
			//   `wb`      `struct wb [10]`
			// not ok:
			//   `wb`      `struct wba`
			postfix := "_c4go_postfix"
			if vType != nil && vName != nil &&
				len(strings.TrimSpace(*vName)) > 0 &&
				strings.Contains(*vType, *vName) {

				for _, pr := range []string{*vName, "struct " + *vName, "union " + *vName} {
					if pr == *vType {
						*vName += postfix
						break
					}
					if len(*vType) > len(pr) && pr == (*vType)[:len(pr)] && len(pr) > 0 {
						letter := (*vType)[len(pr)]
						if unicode.IsLetter(rune(letter)) {
							continue
						}
						if unicode.IsNumber(rune(letter)) {
							continue
						}
						if letter == '*' || letter == '[' || letter == ' ' {
							*vName += postfix
							break
						}
					}
				}
			}
			for i := range node.Children() {
				replacer(node.Children()[i])
			}
		}
		replacer(root)
	}

	// Now begin building the Go AST.
	decls, err := transpileToNode(root, p)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("error of transpiling: err = %v", err), root))
		err = nil // Error is ignored
	}
	p.File.Decls = append(p.File.Decls, decls...)

	// only for "stdbool.h"
	if p.IncludeHeaderIsExists("stdbool.h") {
		p.File.Decls = append(p.File.Decls, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: goast.NewIdent("_Bool"),
					Type: goast.NewIdent("int32"),
				},
			},
		})
	}

	// add functions from CSTD
	std := p.GetCstdFunction()

	// add convertion value to slice
	GetUnsafeConvertDecls(p)

	// checking implementation for all called functions
	bindHeader, bindCode := generateBinding(p, clangFlags)

	// Add the imports after everything else so we can ensure that they are all
	// placed at the top.
	for _, quotedImportPath := range p.Imports() {
		importSpec := &goast.ImportSpec{
			Path: &goast.BasicLit{
				Kind:  token.IMPORT,
				Value: quotedImportPath,
			},
		}
		importDecl := &goast.GenDecl{
			Tok: token.IMPORT,
		}

		importDecl.Specs = append(importDecl.Specs, importSpec)
		p.File.Decls = append([]goast.Decl{importDecl}, p.File.Decls...)
	}

	// generate Go source
	source = p.String()

	// add functions from CSTD
	source += std

	// inject binding code
	if len(bindCode) > 0 {
		index := strings.Index(source, "package")
		index += strings.Index(source[index:], "\n")
		src := source[:index]
		src += "\n"
		src += bindHeader
		src += "\n"
		src += source[index:]
		src += "\n"
		src += bindCode
		source = src
	}

	// only for "stdarg.h"
	if (p.IncludeHeaderIsExists("stdarg.h") && p.IsHaveVaList) || strings.Contains(source, "va_list") {
		source += getVaListStruct()
	}

	// generate pointer arithmetic functions
	source += getPointerArithFunctions(p)

	return
}

func transpileToExpr(node ast.Node, p *program.Program, exprIsStmt bool) (
	expr goast.Expr,
	exprType string,
	preStmts []goast.Stmt,
	postStmts []goast.Stmt,
	err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileToExpr. err = %v", err)
		}
	}()
	if node == nil {
		err = fmt.Errorf("not acceptable nil node")
		return
	}
	defer func() {
		preStmts = nilFilterStmts(preStmts)
		postStmts = nilFilterStmts(postStmts)
	}()

	switch n := node.(type) {
	case *ast.StringLiteral:
		expr, exprType, err = transpileStringLiteral(p, n, false)
		return

	case *ast.FloatingLiteral:
		expr, exprType, err = transpileFloatingLiteral(n), "double", nil

	case *ast.PredefinedExpr:
		expr, exprType, err = transpilePredefinedExpr(n, p)

	case *ast.BinaryConditionalOperator:
		expr, exprType, preStmts, postStmts, err = transpileBinaryConditionalOperator(n, p)

	case *ast.ConditionalOperator:
		expr, exprType, preStmts, postStmts, err = transpileConditionalOperator(n, p)

	case *ast.ArraySubscriptExpr:
		expr, exprType, preStmts, postStmts, err = transpileArraySubscriptExpr(n, p)

	case *ast.BinaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileBinaryOperator(n, p, exprIsStmt)

	case *ast.UnaryOperator:
		expr, exprType, preStmts, postStmts, err = transpileUnaryOperator(n, p)

	case *ast.MemberExpr:
		expr, exprType, preStmts, postStmts, err = transpileMemberExpr(n, p)

	case *ast.ImplicitCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileImplicitCastExpr(n, p, exprIsStmt)

	case *ast.DeclRefExpr:
		expr, exprType, err = transpileDeclRefExpr(n, p)

	case *ast.IntegerLiteral:
		expr, exprType, err = transpileIntegerLiteral(n), "int", nil

	case *ast.ParenExpr:
		expr, exprType, preStmts, postStmts, err = transpileParenExpr(n, p)

	case *ast.CStyleCastExpr:
		expr, exprType, preStmts, postStmts, err = transpileCStyleCastExpr(n, p, exprIsStmt)

	case *ast.CharacterLiteral:
		expr, exprType, err = transpileCharacterLiteral(n), "char", nil

	case *ast.CallExpr:
		expr, exprType, preStmts, postStmts, err = transpileCallExpr(n, p)

	case *ast.CompoundAssignOperator:
		return transpileCompoundAssignOperator(n, p, exprIsStmt)

	case *ast.UnaryExprOrTypeTraitExpr:
		return transpileUnaryExprOrTypeTraitExpr(n, p)

	case *ast.InitListExpr:
		expr, exprType, err = transpileInitListExpr(n, p)

	case *ast.CompoundLiteralExpr:
		expr, exprType, err = transpileCompoundLiteralExpr(n, p)

	case *ast.StmtExpr:
		return transpileStmtExpr(n, p)

	case *ast.ImplicitValueInitExpr:
		return transpileImplicitValueInitExpr(n, p)

	case *ast.OffsetOfExpr:
		expr, exprType, err = transpileOffsetOfExpr(n, p)

	case *ast.VAArgExpr:
		expr, exprType, preStmts, postStmts, err = transpileVAArgExpr(n, p)

	case *ast.ConstantExpr:
		switch len(n.Children()) {
		case 0:
			// ignore
		case 1:
			expr, exprType, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, exprIsStmt)
		default:
			err = fmt.Errorf("ConstantExpr: %v. has many nodes", err)
		}

	case *ast.VisibilityAttr:
		// ignore

	case *ast.WeakAttr:
		// ignore

	default:
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("cannot transpile to expr in transpileToExpr : %T : %#v", node, node), node))
		expr = util.NewNil()
	}

	// Real return is through named arguments.
	return
}

func transpileToStmts(node ast.Node, p *program.Program) (
	stmts []goast.Stmt, err error) {

	if node == nil {
		return nil, nil
	}
	defer func() {
		stmts = nilFilterStmts(stmts)
	}()

	switch n := node.(type) {
	case *ast.DeclStmt:
		stmts, err = transpileDeclStmt(n, p)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("error in DeclStmt: %v", err), n))
			err = nil // Error is ignored
		}
		return
	}

	var (
		stmt      goast.Stmt
		preStmts  []goast.Stmt
		postStmts []goast.Stmt
	)
	stmt, preStmts, postStmts, err = transpileToStmt(node, p)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("error in DeclStmt: %v", err), node))
		err = nil // Error is ignored
	}
	return combineStmts(stmt, preStmts, postStmts), err
}

func transpileToStmt(node ast.Node, p *program.Program) (
	stmt goast.Stmt, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	if node == nil {
		return
	}

	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileToStmt : %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, node))
			err = nil // Error is ignored
		}
	}()
	defer func() {
		preStmts = nilFilterStmts(preStmts)
		postStmts = nilFilterStmts(postStmts)
	}()
	defer func() {
		com := p.GetComments(node.Position())
		for i := range com {
			preStmts = append(preStmts, &goast.ExprStmt{
				X: goast.NewIdent(com[i].Text),
			})
		}
		cg := p.GetMessageComments()
		for i := range cg.List {
			preStmts = append(preStmts, &goast.ExprStmt{
				X: goast.NewIdent(cg.List[i].Text),
			})
		}
	}()

	var expr goast.Expr

	switch n := node.(type) {
	// case *ast.DefaultStmt:
	// 	stmt, err = transpileDefaultStmt(n, p)
	// 	return
	//
	// case *ast.CaseStmt:
	// 	stmt, preStmts, postStmts, err = transpileCaseStmt(n, p)
	// 	return

	case *ast.SwitchStmt:
		stmt, preStmts, postStmts, err = transpileSwitchStmt(n, p)
		return

	case *ast.BreakStmt:
		stmt = &goast.BranchStmt{
			Tok: token.BREAK,
		}
		return

	case *ast.WhileStmt:
		return transpileWhileStmt(n, p)

	case *ast.DoStmt:
		return transpileDoStmt(n, p)

	case *ast.ContinueStmt:
		stmt, err = transpileContinueStmt(n, p)
		return

	case *ast.IfStmt:
		stmt, preStmts, postStmts, err = transpileIfStmt(n, p)
		return

	case *ast.ForStmt:
		return transpileForStmt(n, p)

	case *ast.ReturnStmt:
		return transpileReturnStmt(n, p)

	case *ast.CompoundStmt:
		stmt, preStmts, postStmts, err = transpileCompoundStmt(n, p)
		return

	case *ast.BinaryOperator:
		if n.Operator == "," {
			stmt, preStmts, err = transpileBinaryOperatorComma(n, p)
			return
		}

	case *ast.LabelStmt:
		stmt, preStmts, postStmts, err = transpileLabelStmt(n, p)
		return

	case *ast.GotoStmt:
		stmt, err = transpileGotoStmt(n, p)
		return

	case *ast.GCCAsmStmt:
		// Go does not support inline assembly. See:
		// https://github.com/Konstantin8105/c4go/issues/228
		p.AddMessage(p.GenerateWarningMessage(
			errors.New("cannot transpile asm, will be ignored"), n))

		stmt = &goast.EmptyStmt{}
		return
	case *ast.DeclStmt:
		var stmts []goast.Stmt
		stmts, err = transpileDeclStmt(n, p)
		if err != nil {
			return
		}
		stmt = stmts[len(stmts)-1]
		if len(stmts) > 1 {
			preStmts = stmts[0 : len(stmts)-2]
		}
		return
	}

	// We do not care about the return type.
	var theType string
	expr, theType, preStmts, postStmts, err = transpileToExpr(node, p, true)
	if err != nil {
		return
	}

	// nil is happen, when we remove function `free` of <stdlib.h>
	// see function CallExpr in transpiler
	if expr == (*goast.CallExpr)(nil) {
		return
	}

	// CStyleCastExpr.Kind == ToVoid
	var foundToVoid bool
	if theType == types.ToVoid {
		foundToVoid = true
	}
	if v, ok := node.(*ast.CStyleCastExpr); ok && v.Kind == ast.CStyleCastExprToVoid {
		foundToVoid = true
	}
	if len(node.Children()) > 0 {
		if v, ok := node.Children()[0].(*ast.CStyleCastExpr); ok &&
			v.Kind == ast.CStyleCastExprToVoid {
			foundToVoid = true
		}
	}
	if foundToVoid {
		stmt = &goast.AssignStmt{
			Lhs: []goast.Expr{goast.NewIdent("_")},
			Tok: token.ASSIGN,
			Rhs: []goast.Expr{expr},
		}
		return
	}

	// For all other cases
	if expr == nil {
		err = fmt.Errorf("expr is nil")
		return
	}
	stmt = util.NewExprStmt(expr)

	return
}

func transpileToNode(node ast.Node, p *program.Program) (
	decls []goast.Decl, err error) {
	defer func() {
		if err != nil {
			if _, ok := node.(*ast.RecordDecl); !ok {
				// ignore error for all case except RecordDecl
				p.AddMessage(p.GenerateWarningMessage(err, node))
				err = nil // Error is ignored
			}
		}
	}()

	defer func() {
		decls = nilFilterDecl(decls)
	}()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("transpileToNode: error - panic : %#v. %s", r, string(debug.Stack()))
		}
	}()

	switch n := node.(type) {
	case *ast.TranslationUnitDecl:
		return transpileTranslationUnitDecl(p, n)
	}

	if !AddOutsideStruct {
		if node != nil {
			if (!p.PreprocessorFile.IsUserSource(node.Position().File)) &&
				(!strings.HasSuffix(node.Position().File, "stdint.h")) {
				return
			}
		}
	}

	defer func() {
		if len(decls) > 0 && err == nil {
			for i := range decls {
				if decls[i] == nil {
					continue
				}

				var (
					doc   *goast.CommentGroup
					name  string
					found bool
				)

				if p.Function != nil {
					continue
				}

				switch decls[i].(type) {
				case *goast.GenDecl:
					if decls[i].(*goast.GenDecl).Doc == nil {
						decls[i].(*goast.GenDecl).Doc = &goast.CommentGroup{}
					}
					doc = decls[i].(*goast.GenDecl).Doc
					found = true

					// try to find name
					name = "c4go_name_is_not_found"
					specs := decls[i].(*goast.GenDecl).Specs
					if len(specs) > 0 {
						switch v := specs[0].(type) {
						case *goast.TypeSpec:
							if v.Name != nil {
								name = v.Name.Name
							}

						case *goast.ValueSpec:
							if len(v.Names) > 0 {
								if v.Names[0] != nil {
									name = v.Names[0].Name
								}
							}

						default:
							// ignored
						}
					}

				case *goast.FuncDecl:
					if decls[i].(*goast.FuncDecl).Doc == nil {
						decls[i].(*goast.FuncDecl).Doc = &goast.CommentGroup{}
					}
					if decls[i].(*goast.FuncDecl).Name == nil {
						decls[i].(*goast.FuncDecl).Name = goast.NewIdent("c4go_noname")
					}
					doc = decls[i].(*goast.FuncDecl).Doc
					name = decls[i].(*goast.FuncDecl).Name.Name
					found = true

				default:
					// ignore that goast.Decl
					found = false
					continue
				}

				if !found {
					continue
				}

				com := p.GetComments(node.Position())
				msg := p.GetMessageComments().List
				doc.List = append(doc.List, com...)
				doc.List = append(doc.List, msg...)

				// location of file
				location := node.Position().GetSimpleLocation()
				location = program.PathSimplification(location)
				doc.List = append([]*goast.Comment{{
					Text: fmt.Sprintf("// %s - transpiled function from %s",
						name, location),
				}}, doc.List...)

				break
			}
		}
	}()

	switch n := node.(type) {
	case *ast.FunctionDecl:
		decls, err = transpileFunctionDecl(n, p)

	case *ast.CXXRecordDecl:
		if !strings.Contains(n.RecordDecl.Kind, "class") {
			decls, err = transpileToNode(n.RecordDecl, p)
		} else {
			decls, err = transpileCXXRecordDecl(p, n.RecordDecl)
		}

	case *ast.TypedefDecl:
		decls, err = transpileTypedefDecl(p, n)

	case *ast.RecordDecl:
		decls, err = transpileRecordDecl(p, n)

	case *ast.VarDecl:
		decls, _, err = transpileVarDecl(p, n)

	case *ast.EnumDecl:
		decls, err = transpileEnumDecl(p, n)

	case *ast.LinkageSpecDecl:
		// ignore

	case *ast.EmptyDecl:
		if len(n.Children()) == 0 {
			// ignore if length is zero, for avoid
			// mistake warning
		} else {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("EmptyDecl is not transpiled"), n))
		}
		err = nil
		return

	default:
		err = fmt.Errorf("cannot transpile to node: %#v", node)
	}

	return
}

func transpileStmts(nodes []ast.Node, p *program.Program) (stmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("error in transpileToStmts: %v", err), nodes[0]))
			err = nil // Error is ignored
		}
	}()

	for _, s := range nodes {
		if s != nil {
			var (
				stmt      goast.Stmt
				preStmts  []goast.Stmt
				postStmts []goast.Stmt
			)
			stmt, preStmts, postStmts, err = transpileToStmt(s, p)
			if err != nil {
				return
			}
			stmts = append(stmts, combineStmts(stmt, preStmts, postStmts)...)
		}
	}

	return stmts, nil
}
