// This file contains transpiling for enums.

package transpiler

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	goast "go/ast"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

// ctypeEnumValue generates a specific expression for values used by some
// constants in ctype.h. This is to get around an issue that the real values
// need to be evaulated by the compiler; which c4go does not yet do.
//
// TODO: Ability to evaluate constant expressions at compile time
// https://github.com/Konstantin8105/c4go/issues/77
func ctypeEnumValue(value int, t token.Token) goast.Expr {
	// Produces an expression like: ((1 << (0)) << 8)
	return &goast.ParenExpr{
		X: util.NewBinaryExpr(
			&goast.ParenExpr{
				X: util.NewBinaryExpr(
					util.NewIntLit(1),
					token.SHL,
					util.NewIntLit(value),
					"int",
					false,
				),
			},
			t,
			util.NewIntLit(8),
			"int",
			false,
		),
	}
}

func transpileEnumConstantDecl(p *program.Program, n *ast.EnumConstantDecl) (
	*goast.ValueSpec, []goast.Stmt, []goast.Stmt) {
	var value goast.Expr = util.NewIdent("iota")
	valueType := "int"
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// Special cases for linux ctype.h. See the description for the
	// ctypeEnumValue() function.
	switch n.Name {
	case "_ISupper":
		value = ctypeEnumValue(0, token.SHL) // "((1 << (0)) << 8)"
		valueType = "uint16"
	case "_ISlower":
		value = ctypeEnumValue(1, token.SHL) // "((1 << (1)) << 8)"
		valueType = "uint16"
	case "_ISalpha":
		value = ctypeEnumValue(2, token.SHL) // "((1 << (2)) << 8)"
		valueType = "uint16"
	case "_ISdigit":
		value = ctypeEnumValue(3, token.SHL) // "((1 << (3)) << 8)"
		valueType = "uint16"
	case "_ISxdigit":
		value = ctypeEnumValue(4, token.SHL) // "((1 << (4)) << 8)"
		valueType = "uint16"
	case "_ISspace":
		value = ctypeEnumValue(5, token.SHL) // "((1 << (5)) << 8)"
		valueType = "uint16"
	case "_ISprint":
		value = ctypeEnumValue(6, token.SHL) // "((1 << (6)) << 8)"
		valueType = "uint16"
	case "_ISgraph":
		value = ctypeEnumValue(7, token.SHL) // "((1 << (7)) << 8)"
		valueType = "uint16"
	case "_ISblank":
		value = ctypeEnumValue(8, token.SHR) // "((1 << (8)) >> 8)"
		valueType = "uint16"
	case "_IScntrl":
		value = ctypeEnumValue(9, token.SHR) // "((1 << (9)) >> 8)"
		valueType = "uint16"
	case "_ISpunct":
		value = ctypeEnumValue(10, token.SHR) // "((1 << (10)) >> 8)"
		valueType = "uint16"
	case "_ISalnum":
		value = ctypeEnumValue(11, token.SHR) // "((1 << (11)) >> 8)"
		valueType = "uint16"
	default:
		if len(n.Children()) > 0 {
			var err error
			value, _, preStmts, postStmts, err = transpileToExpr(n.Children()[0], p, false)
			if err != nil {
				panic(err)
			}
		}
	}

	return &goast.ValueSpec{
		Names:  []*goast.Ident{util.NewIdent(n.Name)},
		Type:   util.NewTypeIdent(valueType),
		Values: []goast.Expr{value},
		Doc:    p.GetMessageComments(),
	}, preStmts, postStmts
}

func transpileEnumDecl(p *program.Program, n *ast.EnumDecl) (
	decls []goast.Decl, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileEnumDecl. %v", err)
		}
	}()

	n.Name = types.GenerateCorrectType(n.Name)
	if strings.HasPrefix(n.Name, "enum ") {
		n.Name = n.Name[len("enum "):]
	}

	// For case `enum` without name
	if n.Name == "" {
		return transpileEnumDeclWithType(p, n, "int")
	}

	// For case `enum` with name

	// Create alias of enum for int
	decls = append(decls, &goast.GenDecl{
		Tok: token.TYPE,
		Specs: []goast.Spec{
			&goast.TypeSpec{
				Name: &goast.Ident{
					Name: n.Name,
					Obj:  goast.NewObj(goast.Typ, n.Name),
				},
				// by defaults enum in C is INT
				Type: util.NewTypeIdent("int"),
			},
		},
	})

	// Registration new type in program.Program
	if !p.IsTypeAlreadyDefined(n.Name) {
		p.DefineType(n.Name)
	}

	var d []goast.Decl
	d, err = transpileEnumDeclWithType(p, n, n.Name)
	decls = append(decls, d...)
	return
}

func transpileEnumDeclWithType(p *program.Program, n *ast.EnumDecl, enumType string) (
	decls []goast.Decl, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Cannot transpileEnumDeclWithoutName. %v", err)
		}
	}()
	preStmts := []goast.Stmt{}
	postStmts := []goast.Stmt{}

	// initialization decls
	d := &goast.GenDecl{
		Tok: token.CONST,
	}

	// create all EnumConstant like just constants
	var counter int
	for i, child := range n.Children() {
		switch child.(type) {
		case *ast.FullComment, *ast.BlockCommandComment,
			*ast.HTMLStartTagComment, *ast.HTMLEndTagComment,
			*ast.InlineCommandComment, *ast.ParagraphComment,
			*ast.ParamCommandComment, *ast.TextComment,
			*ast.VerbatimLineComment, *ast.VerbatimBlockComment,
			*ast.VerbatimBlockLineComment:
			// comments are ignored
			continue
		case *ast.EnumConstantDecl:
			// go to next
		default:
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("Unsupported type `%T` in enum.", child), child))
			return
		}
		child := child.(*ast.EnumConstantDecl)
		var (
			e       *goast.ValueSpec
			newPre  []goast.Stmt
			newPost []goast.Stmt
			val     *goast.ValueSpec
		)
		val, newPre, newPost = transpileEnumConstantDecl(p, child)

		if len(newPre) > 0 || len(newPost) > 0 {
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("Check - added in code : (%d)(%d)",
					len(newPre), len(newPost)), n))
		}

		preStmts, postStmts = combinePreAndPostStmts(
			preStmts, postStmts, newPre, newPost)

	remove_parent_expr:
		if v, ok := val.Values[0].(*goast.ParenExpr); ok {
			val.Values[0] = v.X
			goto remove_parent_expr
		}

		var sign int = 1
		if unary, ok := val.Values[0].(*goast.UnaryExpr); ok {
			if unary.Op == token.SUB {
				sign = -1
			}
			val.Values[0] = unary.X
		}

		switch v := val.Values[0].(type) {
		case *goast.Ident:
			e = &goast.ValueSpec{
				Names: []*goast.Ident{{Name: child.Name}},
				Values: []goast.Expr{&goast.BasicLit{
					Kind:  token.INT,
					Value: strconv.Itoa(counter),
				}},
				Type: val.Type,
				Doc:  p.GetMessageComments(),
			}
			counter++

		case *goast.BasicLit:
			var value int
			value, err = strconv.Atoi(v.Value)
			if err != nil {
				err = fmt.Errorf("Cannot parse '%s' in BasicLit", v.Value)
				return
			}
			if err != nil {
				e = val
				counter++
				p.AddMessage(p.GenerateWarningMessage(
					fmt.Errorf("Cannot parse '%s' in BasicLit", v.Value), n))
				break
			}
			if sign == -1 {
				e = &goast.ValueSpec{
					Names: []*goast.Ident{{Name: child.Name}},
					Values: []goast.Expr{&goast.UnaryExpr{
						X: &goast.BasicLit{
							Kind:  token.INT,
							Value: v.Value,
						},
						Op: token.SUB,
					}},
					Type: val.Type,
					Doc:  p.GetMessageComments(),
				}
			} else {
				e = &goast.ValueSpec{
					Names: []*goast.Ident{{Name: child.Name}},
					Values: []goast.Expr{&goast.BasicLit{
						Kind:  token.INT,
						Value: v.Value,
					}},
					Type: val.Type,
					Doc:  p.GetMessageComments(),
				}
			}
			counter = value * sign
			counter++

		default:
			e = val
			p.AddMessage(p.GenerateWarningMessage(
				fmt.Errorf("Add support of continues counter for type : %T",
					v), n))
		}

		valSpec := &goast.ValueSpec{
			Names:  e.Names,
			Values: e.Values,
		}

		if i == 0 {
			valSpec.Type = goast.NewIdent(enumType)
		}

		d.Specs = append(d.Specs, valSpec)
	}
	d.Lparen = 1
	decls = append(decls, d)
	err = nil
	return
}
