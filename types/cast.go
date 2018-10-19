package types

import (
	"fmt"
	"go/token"
	"strings"

	goast "go/ast"

	"strconv"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/util"
)

// GetArrayTypeAndSize returns the size and type of a fixed array. If the type
// is not an array with a fixed size then the the size will be -1 and the
// returned type should be ignored.
func GetArrayTypeAndSize(s string) (string, int) {
	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	match := util.GetRegex(`([\w\* ]*)\[(\d+)\]((\[\d+\])*)`).FindStringSubmatch(s)
	if len(match) > 0 {
		var t = fmt.Sprintf("%s%s", match[1], match[3])
		return strings.Trim(t, " "), util.Atoi(match[2])
	}

	return s, -1
}

// CastExpr returns an expression that casts one type to another. For
// reliability and flexability the existing type (fromType) must be structly
// provided.
//
// There are lots of rules about how an expression is cast, but here are some
// main points:
//
// 1. If fromType == toType (casting to the same type) OR toType == "void *",
//    the original expression is returned unmodified.
//
// 2. There is a special type called "null" which is not defined in C, but
//    rather an estimate of the NULL macro which evaluates to: (0). We cannot
//    guarantee that original C used the NULL macro but it is a safe assumption
//    for now.
//
//    The reason why NULL is special (or at least seamingly) is that it is often
//    used in different value contexts. As a number, testing pointers and
//    strings. Being able to better understand the original purpose of the code
//    helps to generate cleaner and more Go-like output.
//
// 3. There is a set of known primitive number types like "int", "float", etc.
//    These we know can be safely cast between each other by using the data type
//    as a function. For example, 3 (int) to a float would produce:
//    "float32(3)".
//
//    There are also some platform specific types and types that are shared in
//    Go packages that are common aliases kept in this list.
//
// 4. If all else fails the fallback is to cast using a function. For example,
//    Foo -> Bar, would return an expression similar to "noarch.FooToBar(expr)".
//    This code would certainly fail with custom types, but that would likely be
//    a bug. It is most useful to do this when dealing with compound types like
//    FILE where those function probably exist (or should exist) in the noarch
//    package.
func CastExpr(p *program.Program, expr goast.Expr, cFromType, cToType string) (
	_ goast.Expr, err2 error) {

	defer func() {
		if err2 != nil {
			err2 = fmt.Errorf(
				"Cannot casting {%s -> %s}. err = %v", cFromType, cToType, err2)
		}
	}()

	if cFromType == "" {
		err2 = fmt.Errorf("cFromType is empty")
		return
	}

	if cToType == "" {
		err2 = fmt.Errorf("cToType is empty")
		return
	}

	// only for debugging
	// fmt.Printf("%#v -> %#v\n", cFromType, cToType)

	// Uncomment only for debugging
	if strings.Contains(cFromType, ":") {
		err2 = fmt.Errorf("Found mistake `cFromType` `%v` C type : %#v",
			cFromType, expr)
		return
	}
	if strings.Contains(cToType, ":") {
		err2 = fmt.Errorf("Found mistake `cToType` `%v` C type: %#v",
			cToType, expr)
		return
	}

	cFromType = CleanCType(cFromType)
	cToType = CleanCType(cToType)

	// Only for "stddef.h"
	if p.IncludeHeaderIsExists("stddef.h") {
		if cFromType == "long" && cToType == "ptrdiff_t" {
			expr = &goast.BinaryExpr{
				X:  expr,
				Op: token.QUO,
				Y: &goast.BasicLit{
					Kind:  token.INT,
					Value: "8",
				},
			}
		}
	}
	fromType := cFromType
	toType := cToType

	if cFromType == cToType {
		return expr, nil
	}

	if expr == nil {
		return nil, fmt.Errorf("Expr is nil")
	}

	if IsFunction(cFromType) && toType == "bool" {
		return &goast.BinaryExpr{
			X:  expr,
			Op: token.NEQ,
			Y:  goast.NewIdent("nil"),
		}, nil
	}

	// Pointer address
	// Example :
	// From : "int"
	// To   : "int *"
	{
		s := strings.TrimSpace(strings.TrimLeft(cToType, cFromType))
		if s == "*" {
			return &goast.UnaryExpr{
				Op: token.AND,
				X:  expr,
			}, nil
		}
	}

	// Function casting
	// Example :
	// cFromType  : double (int, float, double)
	// cToType    : double (*)(int, float, double)
	if IsFunction(cFromType) {
		return expr, nil
	}

	// Exceptions for stdout, stdin, stderr
	if fromType == "FILE *" && toType == "struct _IO_FILE *" {
		return expr, nil
	}
	if fromType == "struct _IO_FILE *" && toType == "FILE *" {
		return expr, nil
	}

	// casting
	if fromType == "void *" && toType[len(toType)-1] == '*' &&
		!strings.Contains(toType, "FILE") && toType[len(toType)-2] != '*' {
		toType = strings.Replace(toType, "*", " ", -1)
		t, err := ResolveType(p, toType)
		if err != nil {
			return nil, err
		}
		return &goast.TypeAssertExpr{
			X:      expr,
			Lparen: 1,
			Type: &goast.ArrayType{
				Lbrack: 1,
				Elt:    util.NewIdent(t),
			}}, nil
	}

	// Checking amount recursive typedef element
	list := []string{fromType, toType}
	for _, l := range list {
		if v, ok := p.TypedefType[l]; ok {
			var typedefs []string
			for {
				if vv, ok := p.TypedefType[v]; ok {
					for i := range typedefs {
						if vv == typedefs[i] {
							return expr,
								fmt.Errorf("recursive typedef %s in : %v", vv, typedefs)
						}
					}
					v = vv
					typedefs = append(typedefs, vv)
				} else {
					break
				}
			}
		}
	}

	// Checking registated typedef types in program
	if v, ok := p.TypedefType[toType]; ok {
		if fromType == v {
			toType, err := ResolveType(p, toType)
			if err != nil {
				return expr, err
			}

			return &goast.CallExpr{
				Fun: &goast.Ident{
					Name: toType,
				},
				Lparen: 1,
				Args: []goast.Expr{
					&goast.ParenExpr{
						Lparen: 1,
						X:      expr,
						Rparen: 2,
					},
				},
				Rparen: 2,
			}, nil
		}
		e, err := CastExpr(p, expr, fromType, v)
		if err != nil {
			return nil, err
		}
		return CastExpr(p, e, v, toType)
	}
	if v, ok := p.TypedefType[fromType]; ok {
		t, err := ResolveType(p, v)
		if err != nil {
			return expr, err
		}
		expr = &goast.CallExpr{
			Fun: &goast.Ident{
				Name: t,
			},
			Lparen: 1,
			Args: []goast.Expr{
				&goast.ParenExpr{
					Lparen: 1,
					X:      expr,
					Rparen: 2,
				},
			},
			Rparen: 2,
		}
		if toType == v {
			return expr, nil
		}
		return CastExpr(p, expr, v, toType)
	}

	// C null pointer can cast to any pointer
	if cFromType == NullPointer {
		return goast.NewIdent("nil"), nil
	}

	// For : union (anonymous union at sqlite3.c:619241696:3)
	if strings.Contains(fromType, "anonymous union") {
		// I don't understood - How to change correctly
		// Try change to : `union` , but it is FAIL with that
		fromType = ""
	}

	// convert enum to int and recursive
	if strings.Contains(fromType, "enum") && !strings.Contains(toType, "enum") {
		in := goast.CallExpr{
			Fun: &goast.Ident{
				Name: "int",
			},
			Lparen: 1,
			Args: []goast.Expr{
				&goast.ParenExpr{
					Lparen: 1,
					X:      expr,
					Rparen: 2,
				},
			},
			Rparen: 2,
		}
		return CastExpr(p, &in, "int", toType)
	}
	// convert int to enum and recursive
	if !strings.Contains(fromType, "enum") && strings.Contains(toType, "enum") {
		in := goast.CallExpr{
			Fun: &goast.Ident{
				Name: strings.TrimSpace(strings.Replace(toType, "enum", "", -1)),
			},
			Lparen: 1,
			Args: []goast.Expr{
				&goast.ParenExpr{
					Lparen: 1,
					X:      expr,
					Rparen: 2,
				},
			},
			Rparen: 2,
		}
		return CastExpr(p, &in, toType, toType)
	}

	// Let's assume that anything can be converted to a void pointer.
	if toType == "void *" {
		return expr, nil
	}

	fromType, err := ResolveType(p, fromType)
	if err != nil {
		return expr, err
	}

	toType, err = ResolveType(p, toType)
	if err != nil {
		return expr, err
	}

	if fromType == "null" && toType == "[][]byte" {
		return util.NewNil(), nil
	}

	if fromType == "null" && toType == "float64" {
		return util.NewFloatLit(0.0), nil
	}

	if fromType == "null" && toType == "bool" {
		return util.NewIdent("false"), nil
	}

	// FIXME: This is a hack to avoid casting in some situations.
	if fromType == "" || toType == "" {
		return expr, nil
	}

	if fromType == "null" && toType == "[]byte" {
		return util.NewNil(), nil
	}

	// This if for linux.
	if fromType == "*_IO_FILE" && toType == "*noarch.File" {
		return expr, nil
	}

	if fromType == toType {
		return expr, nil
	}

	// Compatible integer types
	types := []string{
		// Integer types
		"byte",
		"int", "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64",

		// Floating-point types.
		"float32", "float64",

		// Known aliases
		"__uint16_t", "size_t",
	}
	for _, v := range types {
		if fromType == v && toType == "bool" {
			e := util.NewBinaryExpr(
				expr,
				token.NEQ,
				util.NewIntLit(0),
				toType,
				false,
			)

			return e, nil
		}
		if fromType == "bool" && toType == v {
			e := util.NewGoExpr(`map[bool]int{false: 0, true: 1}[replaceme]`)
			// Swap replaceme with the current expression
			e.(*goast.IndexExpr).Index = expr
			return CastExpr(p, e, "int", cToType)
		}
	}

	// cast size_t to int
	{
		_, fok := program.DefinitionType[cFromType]
		t, tok := program.DefinitionType[cToType]
		if fok && tok {
			return &goast.CallExpr{
				Fun:  goast.NewIdent(p.ImportType(t)),
				Args: []goast.Expr{expr},
			}, nil
		}
	}

	// In the forms of:
	// - `string` -> `[]byte`
	// - `string` -> `char *[13]`
	match1 := util.GetRegex(`\[\]byte`).FindStringSubmatch(toType)
	match2 := util.GetRegex(`char \*\[(\d+)\]`).FindStringSubmatch(toType)
	if fromType == "string" && (len(match1) > 0 || len(match2) > 0) {
		// Construct a byte array from "first":
		//
		//     var str []byte = []byte{'f','i','r','s','t'}

		value := &goast.CompositeLit{
			Type: &goast.ArrayType{
				Elt: util.NewTypeIdent("byte"),
			},
			Elts: []goast.Expr{},
		}

		strValue, err := strconv.Unquote(expr.(*goast.BasicLit).Value)
		if err != nil {
			panic(fmt.Sprintf("Failed to Unquote %s\n", expr.(*goast.BasicLit).Value))
		}

		for _, c := range []byte(strValue) {
			value.Elts = append(value.Elts, &goast.BasicLit{
				Kind:  token.CHAR,
				Value: fmt.Sprintf("%q", c),
			})
		}

		value.Elts = append(value.Elts, util.NewIntLit(0))

		return value, nil
	}

	// In the forms of:
	// - `[7]byte` -> `string`
	// - `char *[12]` -> `string`
	match1 = util.GetRegex(`\[(\d+)\]byte`).FindStringSubmatch(fromType)
	match2 = util.GetRegex(`char \*\[(\d+)\]`).FindStringSubmatch(fromType)
	if (len(match1) > 0 || len(match2) > 0) && toType == "string" {
		var size int
		if len(match1) > 0 {
			size = util.Atoi(match1[1])
		} else {
			size = util.Atoi(match2[1])
		}

		// The following code builds this:
		//
		//     string(expr[:size - 1])
		//
		return util.NewCallExpr(
			"string",
			&goast.SliceExpr{
				X:    expr,
				High: util.NewIntLit(size - 1),
			},
		), nil
	}

	// Anything that is a pointer can be compared to nil
	if fromType[0] == '*' && toType == "bool" {
		e := util.NewBinaryExpr(expr, token.NEQ, util.NewNil(), toType, false)

		return e, nil
	}

	if fromType == "[]byte" && toType == "bool" {
		return util.NewUnaryExpr(
			token.NOT, util.NewCallExpr("noarch.CStringIsNull", expr),
		), nil
	}

	if fromType == "int" && toType == "*int" {
		return util.NewNil(), nil
	}
	if fromType == "int" && toType == "*byte" {
		return util.NewStringLit(`""`), nil
	}

	if fromType == "_Bool" && toType == "int" {
		return expr, nil
	}

	if util.InStrings(fromType, types) && util.InStrings(toType, types) {
		return util.NewCallExpr(toType, expr), nil
	}

	p.AddImport("github.com/Konstantin8105/c4go/noarch")

	leftName := fromType
	rightName := toType

	if strings.Contains(leftName, ".") {
		parts := strings.Split(leftName, ".")
		leftName = parts[len(parts)-1]
	}
	if strings.Contains(rightName, ".") {
		parts := strings.Split(rightName, ".")
		rightName = parts[len(parts)-1]
	}

	if cFromType == "void *" && cToType == "char *" {
		return expr, nil
	}

	if IsPointer(cFromType) && cToType == "bool" {
		expr = &goast.BinaryExpr{
			Op: token.NEQ,
			X:  expr,
			Y:  goast.NewIdent("nil"),
		}
		return expr, nil
	}

	if IsCInteger(p, cFromType) && IsPointer(cToType) {
		expr = goast.NewIdent("nil")
		return expr, nil
	}

	p.GenerateWarningMessage(fmt.Errorf("Cannot cast types `%s`->`%s`", cFromType, cToType), nil)

	return expr, nil
}

// IsNullExpr tries to determine if the expression is the result of the NULL
// macro. In C, NULL is actually a macro that produces an expression like "(0)".
//
// There are no guarantees if the original C code used the NULL macro, but it is
// usually a pretty good guess when we see this specific exression signature.
//
// Either way the return value from IsNullExpr should not change the
// functionality of the code but can lead to hints that allow the Go produced to
// be cleaner and more Go-like.
func IsNullExpr(n goast.Expr) bool {
	if p1, ok := n.(*goast.ParenExpr); ok {
		if p2, ok := p1.X.(*goast.BasicLit); ok && p2.Value == "0" {
			return true
		}
	}

	return false
}
