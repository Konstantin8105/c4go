package transpiler

import (
	"fmt"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"

	goast "go/ast"
	"go/parser"
	"go/token"
)

func transpileDeclRefExpr(n *ast.DeclRefExpr, p *program.Program) (
	expr *goast.Ident, exprType string, err error) {

	if n.For == "EnumConstant" {
		// clang don`t show enum constant with enum type,
		// so we have to use hack for repair the type
		if v, ok := p.EnumConstantToEnum[n.Name]; ok {
			expr, exprType, err = util.NewIdent(n.Name), v, nil
			return
		}
	}

	if name, ok := program.DefinitionVariable[n.Name]; ok {
		name = p.ImportType(name)
		return util.NewIdent(name), n.Type, nil
	}

	if n.For == "Function" {
		var includeFile string
		includeFile, err = p.GetIncludeFileNameByFunctionSignature(n.Name, n.Type)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		if includeFile != "" && p.IncludeHeaderIsExists(includeFile) {
			name := p.GetFunctionDefinition(n.Name).Substitution
			if strings.Contains(name, ".") {
				p.AddImport(strings.Split(name, ".")[0])
			}
			return goast.NewIdent(name), n.Type, nil
		}
	}

	theType := n.Type

	return util.NewIdent(n.Name), theType, nil
}

func getDefaultValueForVar(p *program.Program, a *ast.VarDecl) (
	expr []goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot getDefaultValueForVar : err = %v", err)
		}
	}()
	if len(a.Children()) == 0 {
		return nil, "", nil, nil, nil
	}

	// Memory allocation is translated into the Go-style.
	if allocSize := getAllocationSizeNode(p, a.Children()[0]); allocSize != nil {
		// type
		var t string
		if v, ok := a.Children()[0].(*ast.ImplicitCastExpr); ok {
			t = v.Type
		}
		if v, ok := a.Children()[0].(*ast.CStyleCastExpr); ok {
			t = v.Type
		}
		if t != "" {
			right, newPre, newPost, err := generateAlloc(p, allocSize, t)
			if err != nil {
				p.AddMessage(p.GenerateWarningMessage(err, a))
				return nil, "", nil, nil, err
			}

			return []goast.Expr{right}, t, newPre, newPost, nil
		}
	}

	if va, ok := a.Children()[0].(*ast.VAArgExpr); ok {
		outType, err := types.ResolveType(p, va.Type)
		if err != nil {
			return nil, "", nil, nil, err
		}
		var argsName string
		if a, ok := va.Children()[0].(*ast.ImplicitCastExpr); ok {
			if a, ok := a.Children()[0].(*ast.DeclRefExpr); ok {
				argsName = a.Name
			} else {
				return nil, "", nil, nil, fmt.Errorf(
					"Expect DeclRefExpr for vaar, but we have %T", a)
			}
		} else {
			return nil, "", nil, nil, fmt.Errorf(
				"Expect ImplicitCastExpr for vaar, but we have %T", a)
		}
		src := fmt.Sprintf(`package main
var temp = func() %s {
	var ret %s
	if v, ok := %s[c4goVaListPosition].(int32); ok{
		// for 'rune' type
		ret = %s(v)
	} else {
		ret = %s[c4goVaListPosition].(%s)
	}
	c4goVaListPosition++
	return ret
}()`, outType,
			outType,
			argsName,
			outType,
			argsName, outType)

		// Create the AST by parsing src.
		fset := token.NewFileSet() // positions are relative to fset
		f, err := parser.ParseFile(fset, "", src, 0)
		if err != nil {
			return nil, "", nil, nil, err
		}

		expr := f.Decls[0].(*goast.GenDecl).Specs[0].(*goast.ValueSpec).Values
		return expr, va.Type, nil, nil, nil
	}

	defaultValue, defaultValueType, newPre, newPost, err := atomicOperation(a.Children()[0], p)
	if err != nil {
		return nil, defaultValueType, newPre, newPost, err
	}

	var values []goast.Expr
	if !types.IsNullExpr(defaultValue) {
		t, err := types.CastExpr(p, defaultValue, defaultValueType, a.Type)
		if !p.AddMessage(p.GenerateWarningMessage(err, a)) {
			values = append(values, t)
			defaultValueType = a.Type
		}
	}

	return values, defaultValueType, newPre, newPost, nil
}

// GenerateFuncType in according to types
// Type: *ast.FuncType {
// .  Func: 13:7
// .  Params: *ast.FieldList {
// .  .  Opening: 13:12
// .  .  List: []*ast.Field (len = 2) {
// .  .  .  0: *ast.Field {
// .  .  .  .  Type: *ast.Ident {
// .  .  .  .  .  NamePos: 13:13
// .  .  .  .  .  Name: "int"
// .  .  .  .  }
// .  .  .  }
// .  .  .  1: *ast.Field {
// .  .  .  .  Type: *ast.Ident {
// .  .  .  .  .  NamePos: 13:17
// .  .  .  .  .  Name: "int"
// .  .  .  .  }
// .  .  .  }
// .  .  }
// .  }
// .  Results: *ast.FieldList {
// .  .  Opening: -
// .  .  List: []*ast.Field (len = 1) {
// .  .  .  0: *ast.Field {
// .  .  .  .  Type: *ast.Ident {
// .  .  .  .  .  NamePos: 13:21
// .  .  .  .  .  Name: "string"
// .  .  .  .  }
// .  .  .  }
// .  .  }
// .  }
// }
func GenerateFuncType(fields, returns []string) *goast.FuncType {
	var ft goast.FuncType
	{
		var fieldList goast.FieldList
		fieldList.Opening = 1
		fieldList.Closing = 2
		for i := range fields {
			fieldList.List = append(fieldList.List, &goast.Field{Type: &goast.Ident{Name: fields[i]}})
		}
		ft.Params = &fieldList
	}
	{
		var fieldList goast.FieldList
		for i := range returns {
			fieldList.List = append(fieldList.List, &goast.Field{Type: &goast.Ident{Name: returns[i]}})
		}
		ft.Results = &fieldList
	}
	return &ft
}

func transpileInitListExpr(e *ast.InitListExpr, p *program.Program) (
	expr goast.Expr, exprType string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileInitListExpr. err = %v", err)
		}
	}()
	resp := []goast.Expr{}
	var hasArrayFiller = false
	e.Type1 = util.GenerateCorrectType(e.Type1)
	e.Type2 = util.GenerateCorrectType(e.Type2)

	for _, node := range e.Children() {
		// Skip ArrayFiller
		if _, ok := node.(*ast.ArrayFiller); ok {
			hasArrayFiller = true
			continue
		}

		var expr goast.Expr
		var err error
		if sl, ok := node.(*ast.StringLiteral); ok {
			expr, _, err = transpileStringLiteral(p, sl, true)
		} else {
			expr, _, _, _, err = transpileToExpr(node, p, true)
		}
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, node))
			return nil, "", err
		}

		resp = append(resp, expr)
	}

	var t goast.Expr
	var cTypeString string

	arrayType, arraySize := types.GetArrayTypeAndSize(e.Type1)
	if arraySize != -1 {
		goArrayType, err := types.ResolveType(p, arrayType)
		p.AddMessage(p.GenerateWarningMessage(err, e))

		cTypeString = fmt.Sprintf("%s[%d]", arrayType, arraySize)

		if hasArrayFiller {
			t = &goast.ArrayType{
				Elt: &goast.Ident{
					Name: goArrayType,
				},
				Len: util.NewIntLit(arraySize),
			}

			// Array fillers do not work with slices.
			// We initialize the array first, then convert to a slice.
			// For example: (&[4]int{1,2})[:]
			return &goast.SliceExpr{
				X: &goast.ParenExpr{
					X: &goast.UnaryExpr{
						Op: token.AND,
						X: &goast.CompositeLit{
							Type: t,
							Elts: resp,
						},
					},
				},
			}, cTypeString, nil
		}
	}

	goType, err := types.ResolveType(p, e.Type1)
	if err != nil {
		return nil, "", err
	}

	return &goast.CompositeLit{
		Type: goast.NewIdent(goType),
		Elts: resp,
	}, e.Type1, nil
}

func transpileDeclStmt(n *ast.DeclStmt, p *program.Program) (
	stmts []goast.Stmt, err error) {

	if len(n.Children()) == 0 {
		return
	}
	var tud ast.TranslationUnitDecl
	tud.ChildNodes = n.Children()
	var decls []goast.Decl
	decls, err = transpileToNode(&tud, p)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil
	}
	stmts = convertDeclToStmt(decls)

	return
}

func transpileArraySubscriptExpr(n *ast.ArraySubscriptExpr, p *program.Program) (
	_ *goast.IndexExpr, theType string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile ArraySubscriptExpr. err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	children := n.Children()

	expression, _, newPre, newPost, err := transpileToExpr(children[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	index, _, newPre, newPost, err := atomicOperation(children[1], p)
	if err != nil {
		return nil, "", nil, nil, err
	}
	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	return &goast.IndexExpr{
		X:     expression,
		Index: index,
	}, n.Type, preStmts, postStmts, nil
}

func transpileMemberExpr(n *ast.MemberExpr, p *program.Program) (
	_ goast.Expr, _ string, preStmts []goast.Stmt, postStmts []goast.Stmt, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile MemberExpr. err = %v", err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	n.Type = util.GenerateCorrectType(n.Type)
	n.Type2 = util.GenerateCorrectType(n.Type2)

	originTypes := []string{n.Type, n.Type2}
	if n.Children()[0] != nil {
		switch v := n.Children()[0].(type) {
		case *ast.ParenExpr:
			originTypes = append(originTypes, v.Type)
			originTypes = append(originTypes, v.Type2)
		}
	}

	lhs, lhsType, newPre, newPost, err := transpileToExpr(n.Children()[0], p, false)
	if err != nil {
		return nil, "", nil, nil, err
	}

	baseType := lhsType
	lhsType = util.GenerateCorrectType(lhsType)

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// lhsType will be something like "struct foo"
	structType := p.GetStruct(lhsType)
	// added for support "struct typedef"
	if structType == nil {
		structType = p.GetStruct("struct " + lhsType)
	}
	// added for support "union typedef"
	if structType == nil {
		structType = p.GetStruct("union " + lhsType)
	}
	// for anonymous structs
	if structType == nil {
		structType = p.GetStruct(baseType)
	}
	// for anonymous structs
	if structType == nil {
		structType = p.GetStruct(util.CleanCType(baseType))
	}
	// typedef types
	if structType == nil {
		structType = p.GetStruct(p.TypedefType[baseType])
	}
	if structType == nil {
		t := types.GetBaseType(baseType)
		structType = p.GetStruct(p.TypedefType[t])
	}
	// other case
	for _, t := range originTypes {
		if structType == nil {
			structType = p.GetStruct(util.CleanCType(t))
		} else {
			break
		}
		if structType == nil {
			structType = p.GetStruct(types.GetBaseType(t))
		} else {
			break
		}
	}

	if n.Name == "" {
		n.Name = generateNameFieldDecl(util.GenerateCorrectType(n.Type))
	}
	rhs := n.Name
	rhsType := "void *"
	if structType == nil {
		// This case should not happen in the future. Any structs should be
		// either parsed correctly from the source or be manually setup when the
		// parser starts if the struct if hidden or shared between libraries.
		//
		// Some other things to keep in mind:
		//   1. Types need to be stripped of their pointer, 'FILE *' -> 'FILE'.
		//   2. Types may refer to one or more other types in a chain that have
		//      to be resolved before the real field type can be determined.
		err = fmt.Errorf("cannot determine type for LHS '%v'"+
			", will use 'void *' for all fields. Is lvalue = %v. n.Name = %v",
			lhsType, n.IsLvalue, n.Name)
		p.AddMessage(p.GenerateWarningMessage(err, n))
	} else {
		if s, ok := structType.Fields[rhs].(string); ok {
			rhsType = s
		} else {
			err = fmt.Errorf("cannot determine type for RHS '%v', will use"+
				" 'void *' for all fields. Is lvalue = %v. n.Name = `%v`",
				rhs, n.IsLvalue, n.Name)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}

	x := lhs
	if n.IsPointer {
		x = &goast.IndexExpr{X: x, Index: util.NewIntLit(0)}
	}

	// Check for member name translation.
	lhsType = strings.TrimSpace(lhsType)
	if lhsType[len(lhsType)-1] == '*' {
		lhsType = lhsType[:len(lhsType)-len(" *")]
	}
	if str := p.GetStruct("c4go_" + lhsType); str != nil {
		if alias, ok := str.Fields[rhs]; ok {
			rhs = alias.(string)
			goto Selector
		}
	}

	// anonymous struct member?
	if rhs == "" {
		rhs = "anon"
	}

	if isUnionMemberExpr(p, n) {
		return &goast.ParenExpr{
			Lparen: 1,
			X: &goast.StarExpr{
				Star: 1,
				X: &goast.CallExpr{
					Fun: &goast.SelectorExpr{
						X:   x,
						Sel: util.NewIdent(rhs),
					},
					Lparen: 1,
				},
			},
		}, n.Type, preStmts, postStmts, nil
	}

Selector:
	_ = rhsType

	return &goast.SelectorExpr{
		X:   x,
		Sel: util.NewIdent(rhs),
	}, n.Type, preStmts, postStmts, nil
}

func transpileImplicitValueInitExpr(n *ast.ImplicitValueInitExpr, p *program.Program) (
	expr goast.Expr, exprType string, _ []goast.Stmt, _ []goast.Stmt, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileImplicitValueInitExpr. err = %v", err)
		}
	}()

	exprType = n.Type1

	var t string
	t = n.Type1
	t, err = types.ResolveType(p, t)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	var isStruct bool
	if _, ok := p.Structs[t]; ok {
		isStruct = true
	}
	if _, ok := p.Structs["struct "+t]; ok {
		isStruct = true
	}
	if isStruct {
		expr = &goast.CompositeLit{
			Type:   util.NewIdent(t),
			Lbrace: 1,
		}
		return
	}

	expr = util.NewCallExpr(t,
		&goast.BasicLit{
			Kind:  token.INT,
			Value: "0",
		})
	return
}
