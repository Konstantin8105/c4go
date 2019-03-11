// This file contains functions for transpiling declarations of variables and
// types. The usage of variables is handled in variables.go.

package transpiler

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

// Example of AST for union without name inside struct:
// -RecordDecl 0x40d41b0 <...> line:453:8 struct EmptyName definition
//  |-RecordDecl 0x40d4260 <...> line:454:2 union definition
//  | |-FieldDecl 0x40d4328 <...> col:8 referenced l1 'long'
//  | `-FieldDecl 0x40d4388 <...> col:8 referenced l2 'long'
//  |-FieldDecl 0x40d4420 <...> col:2 implicit referenced 'union EmptyName::(anonymous at struct.c:454:2)'
//  |-IndirectFieldDecl 0x40d4478 <...> col:8 implicit l1 'long'
//  | |-Field 0x40d4420 '' 'union EmptyName::(anonymous at /struct.c:454:2)'
//  | `-Field 0x40d4328 'l1' 'long'
//  `-IndirectFieldDecl 0x40d44c8 <...> col:8 implicit l2 'long'
//    |-Field 0x40d4420 '' 'union EmptyName::(anonymous at /struct.c:454:2)'
//    `-Field 0x40d4388 'l2' 'long'

func newFunctionField(p *program.Program, name, cType string) (
	_ *goast.Field, err error) {
	if name == "" {
		err = fmt.Errorf("Name of function field cannot be empty")
		return
	}
	if !util.IsFunction(cType) {
		err = fmt.Errorf("Cannot create function field for type : %s", cType)
		return
	}

	// TODO : add err handling
	fieldType, _ := types.ResolveType(p, cType)

	return &goast.Field{
		Names: []*goast.Ident{util.NewIdent(name)},
		Type:  goast.NewIdent(fieldType),
	}, nil
}

func generateNameFieldDecl(t string) string {
	return "implicit_" + strings.Replace(t, " ", "S", -1)
}

func transpileFieldDecl(p *program.Program, n *ast.FieldDecl) (
	field *goast.Field, err error) {
	defer func() {
		if field != nil {
			if field.Type == nil {
				err = fmt.Errorf("Found nil transpileFieldDecl in field Type %v , %v : %v",
					n, field, err)
				field.Type = util.NewIdent(n.Type)
			}
		}
	}()
	if util.IsFunction(n.Type) {
		field, err = newFunctionField(p, n.Name, n.Type)
		if err == nil {
			return
		}
	}

	if n.Name == "" {
		//&ast.FieldDecl{Addr:0x3157420, Pos:ast.Position{...}, Position2:"col:2", Name:"", Type:"union EmptyNameDD__at__home_lepricon_go_src_github_com_Konstantin8105_c4go_tests_struct_c_454_2_", Type2:"", Implicit:true, Referenced:true, ChildNodes:[]ast.Node{}}
		n.Name = generateNameFieldDecl(n.Type)
	}

	name := n.Name

	fieldType, err := types.ResolveType(p, n.Type)
	p.AddMessage(p.GenerateWarningMessage(err, n))

	// TODO: The name of a variable or field cannot be a reserved word
	// https://github.com/Konstantin8105/c4go/issues/83
	// Search for this issue in other areas of the codebase.
	if util.IsGoKeyword(name) {
		name += "_"
	}

	arrayType, arraySize := types.GetArrayTypeAndSize(n.Type)
	if arraySize != -1 {
		fieldType, err = types.ResolveType(p, arrayType)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		fieldType = fmt.Sprintf("[%d]%s", arraySize, fieldType)
	}

	return &goast.Field{
		Names: []*goast.Ident{util.NewIdent(name)},
		Type:  util.NewTypeIdent(fieldType),
	}, nil
}

func transpileRecordDecl(p *program.Program, n *ast.RecordDecl) (
	decls []goast.Decl, err error) {

	var addPackageUnsafe bool

	n.Name = util.GenerateCorrectType(n.Name)
	name := n.Name
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileRecordDecl `%v`. %v",
				n.Name, err)
		} else {
			if addPackageUnsafe {
				p.AddImports("unsafe")
			}
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error - panic : %#v", r)
		}
	}()

	// ignore if haven`t definition
	if !n.IsDefinition {
		return
	}

	if name == "" || p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	name = util.GenerateCorrectType(name)
	p.DefineType(name)
	defer func() {
		if err != nil {
			p.UndefineType(name)
		}
	}()

	var fields []*goast.Field

	// repair name for anonymous RecordDecl
	for pos := range n.Children() {
		if rec, ok := n.Children()[pos].(*ast.RecordDecl); ok && rec.Name == "" {
			if pos < len(n.Children()) {
				switch v := n.Children()[pos+1].(type) {
				case *ast.FieldDecl:
					rec.Name = types.GetBaseType(util.GenerateCorrectType(v.Type))
				default:
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Cannot find name for anon RecordDecl: %T",
							v), n))
					rec.Name = "UndefinedNameC2GO"
				}
			}
		}
	}

	for pos := range n.Children() {
		switch field := n.Children()[pos].(type) {
		case *ast.FieldDecl:
			field.Type = util.GenerateCorrectType(field.Type)
			field.Type2 = util.GenerateCorrectType(field.Type2)
			var f *goast.Field
			f, err = transpileFieldDecl(p, field)
			if err != nil {
				err = fmt.Errorf("cannot transpile field. %v", err)
				p.AddMessage(p.GenerateWarningMessage(err, field))
				// TODO ignore error
				// return
				err = nil
			} else {
				// ignore fields without name
				if len(f.Names) != 1 {
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Ignore FieldDecl with more then 1 names"+
							" in RecordDecl : `%v`", n.Name), n))
					continue
				}
				if f.Names[0].Name == "" {
					p.AddMessage(p.GenerateWarningMessage(
						fmt.Errorf("Ignore FieldDecl without name "+
							" in RecordDecl : `%v`", n.Name), n))
					continue
				}
				// remove dublicates of fields
				var isDublicate bool
				for i := range fields {
					if fields[i].Names[0].Name == f.Names[0].Name {
						isDublicate = true
					}
				}
				if isDublicate {
					f.Names[0].Name += strconv.Itoa(pos)
				}
				fields = append(fields, f)
			}

		case *ast.IndirectFieldDecl:
			// ignore

		case *ast.TransparentUnionAttr:
			// Don't do anythink
			// Example of AST:
			// |-RecordDecl 0x3632d78 <...> line:67:9 union definition
			// | |-TransparentUnionAttr 0x3633050 <...>
			// | |-FieldDecl 0x3632ed0 <...> col:17 __uptr 'union wait *'
			// | `-FieldDecl 0x3632f60 <...> col:10 __iptr 'int *'
			// |-TypedefDecl 0x3633000 <...> col:5 __WAIT_STATUS 'union __WAIT_STATUS':'__WAIT_STATUS'
			// | `-ElaboratedType 0x3632fb0 'union __WAIT_STATUS' sugar
			// |   `-RecordType 0x3632e00 '__WAIT_STATUS'
			// |     `-Record 0x3632d78 ''

		default:
			// For case anonymous enum:

			// |-EnumDecl 0x26c3970 <...> line:77:5
			// | `-EnumConstantDecl 0x26c3a50 <...> col:9 referenced SWE_ENUM_THREE 'int'
			// |   `-IntegerLiteral 0x26c3a30 <...> 'int' 3
			// |-FieldDecl 0x26c3af0 <...> col:7 EnumThree 'enum (anonymous enum at ...
			if eDecl, ok := field.(*ast.EnumDecl); ok && eDecl.Name == "" {
				if pos+1 <= len(n.Children())-1 {
					if f, ok := n.Children()[pos+1].(*ast.FieldDecl); ok {
						n.Children()[pos].(*ast.EnumDecl).Name = f.Type
					}
				}
			}

			// default
			var declsIn []goast.Decl
			declsIn, err = transpileToNode(field, p)
			if err != nil {
				err = fmt.Errorf("Cannot transpile %T", field)
				// p.AddMessage(p.GenerateWarningMessage(err, field))
				return
			}
			decls = append(decls, declsIn...)
		}
	}

	s, err := program.NewStruct(p, n)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return
	}
	switch s.Type {
	case program.UnionType:
		if strings.HasPrefix(s.Name, "union ") {
			p.Structs[s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, s.Name)
					p.UndefineType(s.Name)
				}
			}()
		} else {
			p.Unions["union "+s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, "union "+s.Name)
					p.UndefineType("union " + s.Name)
				}
			}()
		}

	case program.StructType:
		if strings.HasPrefix(s.Name, "struct ") {
			p.Structs[s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, s.Name)
					p.UndefineType(s.Name)
				}
			}()
		} else {
			p.Structs["struct "+s.Name] = s
			defer func() {
				if err != nil {
					delete(p.Structs, "struct "+s.Name)
					p.UndefineType("struct " + s.Name)
				}
			}()
		}

	default:
		err = fmt.Errorf("Undefine type of struct : %v", s.Type)
		return
	}

	name = strings.TrimPrefix(name, "struct ")
	name = strings.TrimPrefix(name, "union ")

	var d []goast.Decl
	switch s.Type {
	case program.UnionType:
		// Union size
		var size int
		size, err = types.SizeOf(p, "union "+name)

		// In normal case no error is returned,
		if err != nil {
			// but if we catch one, send it as a warning
			err = fmt.Errorf("could not determine the size of type `union %s`"+
				" for that reason: %s", name, err)
			return
		}
		// So, we got size, then
		// Add imports needed
		addPackageUnsafe = true

		// Declaration for implementing union type
		d, err = transpileUnion(name, size, fields)
		if err != nil {
			return nil, err
		}

	case program.StructType:
		d = append(d, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: util.NewIdent(name),
					Type: &goast.StructType{
						Fields: &goast.FieldList{
							List: fields,
						},
					},
				},
			},
		})

	default:
		err = fmt.Errorf("Undefine type of struct : %v", s.Type)
		return
	}

	decls = append(decls, d...)

	return
}

func transpileCXXRecordDecl(p *program.Program, n *ast.RecordDecl) (
	decls []goast.Decl, err error) {

	n.Name = util.GenerateCorrectType(n.Name)
	name := n.Name

	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpileCXXRecordDecl : `%v`. %v",
				n.Name, err)
			p.AddMessage(p.GenerateWarningMessage(err, n))
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error - panic : %#v", r)
		}
	}()

	// ignore if haven`t definition
	if !n.IsDefinition {
		return
	}

	if name == "" || p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	p.DefineType(n.Kind + " " + name)
	defer func() {
		if err != nil {
			p.UndefineType(n.Kind + " " + name)
		}
	}()

	var fields []*goast.Field
	for _, v := range n.Children() {
		switch v := v.(type) {
		case *ast.CXXRecordDecl:
			// ignore

		case *ast.FieldDecl:
			var f *goast.Field
			f, err = transpileFieldDecl(p, v)
			if err != nil {
				return
			}
			fields = append(fields, f)

		default:
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("Cannot transpilation field in CXXRecordDecl : %T", v), n))
		}
	}

	return []goast.Decl{&goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: util.NewIdent(name),
				Type: &goast.StructType{
					Fields: &goast.FieldList{
						List: fields,
					},
				},
			},
		},
	}}, nil
}

func transpileTypedefDecl(p *program.Program, n *ast.TypedefDecl) (
	decls []goast.Decl, err error) {

	// implicit code from clang at the head of each clang AST tree
	if n.IsImplicit && n.Pos.File == ast.PositionBuiltIn {
		return
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpile Typedef Decl : err = %v", err)
		} else {
			if !p.IncludeHeaderIsExists(n.Pos.File) {
				// no need add struct from C STD
				decls = nil
				return
			}
		}
	}()
	n.Name = util.CleanCType(util.GenerateCorrectType(n.Name))
	n.Type = util.CleanCType(util.GenerateCorrectType(n.Type))
	n.Type2 = util.CleanCType(util.GenerateCorrectType(n.Type2))
	name := n.Name

	if "struct "+n.Name == n.Type || "union "+n.Name == n.Type {
		p.TypedefType[n.Name] = n.Type
		return
	}

	if util.IsFunction(n.Type) {
		var field *goast.Field
		field, err = newFunctionField(p, n.Name, n.Type)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
		} else {
			// registration type
			p.TypedefType[n.Name] = n.Type

			decls = append(decls, &goast.GenDecl{
				Tok: token.TYPE,
				Specs: []goast.Spec{
					&goast.TypeSpec{
						Name: util.NewIdent(name),
						Type: field.Type,
					},
				},
			})
			err = nil
			return
		}
	}

	// added for support "typedef enum {...} dd" with empty name of struct
	// Result in Go: "type dd int"
	if strings.Contains(n.Type, "enum") {
		// Registration new type in program.Program
		if !p.IsTypeAlreadyDefined(n.Name) {
			p.DefineType(n.Name)
			p.EnumTypedefName[n.Name] = true
		}
		decls = append(decls, &goast.GenDecl{
			Tok: token.TYPE,
			Specs: []goast.Spec{
				&goast.TypeSpec{
					Name: util.NewIdent(name),
					Type: util.NewTypeIdent("int32"),
				},
			},
		})
		err = nil
		return
	}

	if p.IsTypeAlreadyDefined(name) {
		err = nil
		return
	}

	p.DefineType(name)

	resolvedType, err := types.ResolveType(p, n.Type)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
	}

	// There is a case where the name of the type is also the definition,
	// like:
	//
	//     type _RuneEntry _RuneEntry
	//
	// This of course is impossible and will cause the Go not to compile.
	// It itself is caused by lack of understanding (at this time) about
	// certain scenarios that types are defined as. The above example comes
	// from:
	//
	//     typedef struct {
	//        // ... some fields
	//     } _RuneEntry;
	//
	// Until which time that we actually need this to work I am going to
	// suppress these.
	if name == resolvedType {
		err = nil
		return
	}

	err = nil
	if resolvedType == "" {
		resolvedType = "interface{}"
	}

	p.TypedefType[n.Name] = n.Type

	// 0: *ast.GenDecl {
	// .  Tok: type
	// .  Specs: []ast.Spec (len = 1) {
	// .  .  0: *ast.TypeSpec {
	// .  .  .  Name: *ast.Ident {
	// .  .  .  .  Name: "R"
	// .  .  .  }
	// .  .  .  Assign: 3:8        // <- This is important
	// .  .  .  Type: *ast.Ident {
	// .  .  .  .  Name: "int"
	// .  .  .  }
	// .  .  }
	// .  }
	// }
	decls = append(decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name:   util.NewIdent(name),
				Assign: 1,
				Type:   util.NewTypeIdent(resolvedType),
			},
		},
	})

	return
}

func transpileVarDecl(p *program.Program, n *ast.VarDecl) (
	decls []goast.Decl, theType string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileVarDecl : err = %v", err)
		}
	}()

	n.Name = util.GenerateCorrectType(n.Name)
	n.Type = util.GenerateCorrectType(n.Type)
	n.Type2 = util.GenerateCorrectType(n.Type2)

	// Ignore extern as there is no analogy for Go right now.
	if n.IsExtern && len(n.ChildNodes) == 0 {
		return
	}

	if strings.Contains(n.Type, "va_list") &&
		strings.Contains(n.Type2, "va_list_tag") {
		// variable for va_list. see "variadic function"
		// header : <stdarg.h>
		// Example :
		// DeclStmt 0x2fd87e0 <line:442:2, col:14>
		// `-VarDecl 0x2fd8780 <col:2, col:10> col:10 used args 'va_list':'struct __va_list_tag [1]'
		// Result:
		// ... - convert to - c4goArgs ...interface{}
		// var args = c4goArgs
		return []goast.Decl{&goast.GenDecl{
			Tok: token.VAR,
			Specs: []goast.Spec{
				&goast.ValueSpec{
					Names:  []*goast.Ident{util.NewIdent(n.Name)},
					Values: []goast.Expr{util.NewIdent("c4goArgs")},
				},
			},
		}}, "", nil
	}

	// Example of DeclStmt for C code:
	// void * a = NULL;
	// void(*t)(void) = a;
	// Example of AST:
	// `-VarDecl 0x365fea8 <col:3, col:20> col:9 used t 'void (*)(void)' cinit
	//   `-ImplicitCastExpr 0x365ff48 <col:20> 'void (*)(void)' <BitCast>
	//     `-ImplicitCastExpr 0x365ff30 <col:20> 'void *' <LValueToRValue>
	//       `-DeclRefExpr 0x365ff08 <col:20> 'void *' lvalue Var 0x365f8c8 'r' 'void *'

	if len(n.Children()) > 0 {
		if v, ok := (n.Children()[0]).(*ast.ImplicitCastExpr); ok {
			if len(v.Type) > 0 {
				// Is it function ?
				if util.IsFunction(v.Type) {
					var prefix string
					var fields, returns []string
					prefix, fields, returns, err = types.SeparateFunction(p, v.Type)
					if err != nil {
						err = fmt.Errorf("Cannot resolve function : %v", err)
						return
					}
					if len(prefix) != 0 {
						p.AddMessage(p.GenerateWarningMessage(
							fmt.Errorf("Prefix is not used : `%v`", prefix), n))
					}
					functionType := GenerateFuncType(fields, returns)
					nameVar1 := n.Name

					if vv, ok := v.Children()[0].(*ast.ImplicitCastExpr); ok {
						if decl, ok := vv.Children()[0].(*ast.DeclRefExpr); ok {
							nameVar2 := decl.Name

							return []goast.Decl{&goast.GenDecl{
								Tok: token.VAR,
								Specs: []goast.Spec{&goast.ValueSpec{
									Names: []*goast.Ident{{Name: nameVar1}},
									Type:  functionType,
									Values: []goast.Expr{&goast.TypeAssertExpr{
										X:    &goast.Ident{Name: nameVar2},
										Type: functionType,
									}},
									Doc: p.GetMessageComments(),
								},
								}}}, "", nil
						}
					}
				}
			}
		}
	}

	theType = n.Type

	p.GlobalVariables[n.Name] = theType

	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	defaultValue, _, newPre, newPost, err := getDefaultValueForVar(p, n)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil // Error is ignored
	}
	// for ignore zero value. example:
	// int i = 0;
	// tranpile to:
	// var i int // but not "var i int = 0"
	if len(defaultValue) == 1 && defaultValue[0] != nil {
		if bl, ok := defaultValue[0].(*goast.BasicLit); ok {
			if bl.Kind == token.INT && bl.Value == "0" {
				defaultValue = nil
			}
			if bl.Kind == token.FLOAT && bl.Value == "0" {
				defaultValue = nil
			}
		} else if call, ok := defaultValue[0].(*goast.CallExpr); ok {
			if len(call.Args) == 1 {
				if bl, ok := call.Args[0].(*goast.BasicLit); ok {
					if bl.Kind == token.INT && bl.Value == "0" {
						defaultValue = nil
					}
					if bl.Kind == token.FLOAT && bl.Value == "0" {
						defaultValue = nil
					}
				}
			}
		} else if ind, ok := defaultValue[0].(*goast.Ident); ok {
			if ind.Name == "nil" {
				defaultValue = nil
			}
		}
	}

	preStmts, postStmts = combinePreAndPostStmts(preStmts, postStmts, newPre, newPost)

	// Allocate slice so that it operates like a fixed size array.
	arrayType, arraySize := types.GetArrayTypeAndSize(n.Type)

	if arraySize != -1 && defaultValue == nil {
		var goArrayType string
		goArrayType, err = types.ResolveType(p, arrayType)
		if err != nil {
			p.AddMessage(p.GenerateWarningMessage(err, n))
			err = nil // Error is ignored
		}

		defaultValue = []goast.Expr{
			util.NewCallExpr(
				"make",
				&goast.ArrayType{
					Elt: util.NewTypeIdent(goArrayType),
				},
				util.NewIntLit(arraySize),
				// If len and capacity is same, then
				// capacity is not need
				// util.NewIntLit(arraySize),
			),
		}
	}

	if len(preStmts) != 0 || len(postStmts) != 0 {
		p.AddMessage(p.GenerateWarningMessage(
			fmt.Errorf("Not acceptable length of Stmt : pre(%d), post(%d)",
				len(preStmts), len(postStmts)), n))
	}

	names := map[string]string{
		"ptrdiff_t": "noarch.PtrdiffT",
	}

	var typeResult goast.Expr

	if n, ok := names[n.Type]; ok {
		typeResult = util.NewTypeIdent(n)
		goto ignoreType
	}

	theType, err = types.ResolveType(p, n.Type)
	if err != nil {
		p.AddMessage(p.GenerateWarningMessage(err, n))
		err = nil // Error is ignored
		theType = "UnknownType"
	}
	typeResult = util.NewTypeIdent(theType)

ignoreType:

	return []goast.Decl{&goast.GenDecl{
		Tok: token.VAR,
		Specs: []goast.Spec{
			&goast.ValueSpec{
				Names:  []*goast.Ident{util.NewIdent(n.Name)},
				Type:   typeResult,
				Values: defaultValue,
				Doc:    p.GetMessageComments(),
			},
		},
	}}, "", nil
}
