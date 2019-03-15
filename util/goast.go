// This file contains utility and helper methods for making it easier to
// generate parts of the Go AST.

package util

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// NewExprStmt returns a new ExprStmt from an expression. It is used when
// converting a single expression into a statement for another receiver.
//
// It is recommended you use this method of instantiating the ExprStmt yourself
// because NewExprStmt will check that the expr is not nil (or panic). This is
// much more helpful when trying to debug why the Go source build crashes
// because of a nil pointer - which eventually leads back to a nil expr.
func NewExprStmt(expr goast.Expr) *goast.ExprStmt {
	PanicIfNil(expr, "expr is nil")

	return &goast.ExprStmt{
		X: expr,
	}
}

// IsAValidFunctionName performs a check to see if a string would make a
// valid function name in Go. Go allows unicode characters, but C doesn't.
func IsAValidFunctionName(s string) bool {
	return GetRegex(`^[a-zA-Z_]`).
		Match([]byte(s))
}

// Convert a type as a string into a Go AST expression.
func typeToExpr(t string) goast.Expr {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("bad type: '%v'", t))
		}
	}()

	return internalTypeToExpr(t)
}

func internalTypeToExpr(goType string) goast.Expr {

	// I'm not sure if this is an error or not. It is caused by processing the
	// resolved type of "void" which is "". It is used on functions to denote
	// that it does not have a return type.
	if goType == "" {
		return nil
	}

	separator := make([]bool, len(goType)+1)
	for i := range goType {
		switch goType[i] {
		case '.', '*', '(', ')', '-', '+', '&', '{', '}', ' ', '[', ']':
			separator[i] = true
			separator[i+1] = true
		}
	}

	// Specific case for 'interface{}'
	// remove all separator inside that word
	specials := [][]byte{[]byte("func("), []byte("interface{}")}

	for _, special := range specials {
		input := []byte(goType)
	again:
		index := bytes.Index(input, special)
		if index >= 0 {
			for i := index + 1; i < index+len(special); i++ {
				separator[i] = false
			}
			input = input[index+len(special)-1:]
			goto again
		}
	}

	separator[0] = true
	separator[len(separator)-1] = true

	// Separation string 'goType' to slice of bytes
	var indexes []int
	for i := range separator {
		if separator[i] {
			indexes = append(indexes, i)
		}
	}
	var lines [][]byte
	for i := 0; i < len(indexes)-1; i++ {
		lines = append(lines, []byte(goType[indexes[i]:indexes[i+1]]))
	}

	// Checking
	for i := range lines {
		if IsGoKeyword(string(lines[i])) {
			lines[i] = []byte(string(lines[i]) + "_")
		}
	}

	goType = string(bytes.Join(lines, []byte("")))

	return goast.NewIdent(goType)
}

// NewCallExpr creates a new *"go/ast".CallExpr with each of the arguments
// (after the function name) being each of the expressions that represent the
// individual arguments.
//
// The function name is checked with IsAValidFunctionName and will panic if the
// function name is deemed to be not valid.
func NewCallExpr(functionName string, args ...goast.Expr) *goast.CallExpr {
	for i := range args {
		// Argument of function is cannot be nil
		if args[i] == nil {
			args[i] = goast.NewIdent("C4GO_NOT_TRANSPILED_EXPRESSION")
		}
	}
	index := strings.Index(functionName, ".")
	if index > 0 {
		if IsGoPackage(functionName[:index]) {
			return &goast.CallExpr{
				Fun:  goast.NewIdent(functionName),
				Args: args,
			}
		}
	}
	return &goast.CallExpr{
		Fun:  typeToExpr(functionName),
		Args: args,
	}
}

// NewFuncClosure creates a new *"go/ast".CallExpr that calls a function
// literal closure. The first argument is the Go return type of the
// closure, and the remainder of the arguments are the statements of the
// closure body.
func NewFuncClosure(returnType string, stmts ...goast.Stmt) *goast.CallExpr {
	return &goast.CallExpr{
		Fun: &goast.FuncLit{
			Type: NewFuncType(&goast.FieldList{}, returnType, false),
			Body: &goast.BlockStmt{
				List: stmts,
			},
		},
		Args: []goast.Expr{},
	}
}

// NewBinaryExpr create a new Go AST binary expression with a left, operator and
// right operand.
//
// You should use this instead of BinaryExpr directly so that nil left and right
// operands can be caught (and panic) before Go tried to render the source -
// which would result in a very hard to debug error.
//
// Assignment operators in C can be nested inside other expressions, like:
//
//     a + (b += 3)
//
// In Go this is not allowed. Since the operators mutate variables it is not
// possible in some cases to move the statements before or after. The only safe
// (and generic) way around this is to create an immediately executing closure,
// like:
//
//     a + (func () int { b += 3; return b }())
//
// In a lot of cases this may be unnecessary and obfuscate the Go output but
// these will have to be optimised over time and be strict about the
// situation they are simplifying.
//
// If stmt is true then the binary expression is the whole statement. This means
// that the closure above does not need to applied. This makes the output code
// much neater.
func NewBinaryExpr(left goast.Expr, operator token.Token, right goast.Expr,
	returnType string, stmt bool) goast.Expr {
	PanicIfNil(left, "left is nil")
	PanicIfNil(right, "right is nil")

	var b goast.Expr = &goast.BinaryExpr{
		X:  left,
		Op: operator,
		Y:  right,
	}
	return b
}

// NewIdent - create a new Go ast Ident
func NewIdent(name string) *goast.Ident {
	// TODO: The name of a variable or field cannot be a reserved word
	// https://github.com/Konstantin8105/c4go/issues/83
	// Search for this issue in other areas of the codebase.
	if IsGoKeyword(name) {
		name += "_"
	}

	// Remove const prefix as it has no equivalent in Go.
	name = strings.TrimPrefix(name, "const ")

	if !IsAValidFunctionName(name) {
		// Normally we do not panic because we want the transpiler to recover as
		// much as possible so that we always get Go output - even if it's
		// wrong. However, in this case we must panic because we know that this
		// identity will cause the AST renderer in Go to panic with a very
		// unhelpful error message.
		//
		// Panic now so that we can see where the bad identifier is coming from.
		panic(fmt.Sprintf("invalid identity: '%s'", name))
	}

	return goast.NewIdent(name)
}

// NewTypeIdent created a new Go identity that is to be used for a Go type. This
// is different from NewIdent in how the input string is validated.
func NewTypeIdent(name string) goast.Expr {
	return typeToExpr(name)
}

// NewStringLit returns a new Go basic literal with a string value.
func NewStringLit(value string) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.STRING,
		Value: value,
	}
}

// NewIntLit - create a Go ast BasicLit for `INT` value
func NewIntLit(value int) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(value),
	}
}

// NewFloatLit creates a new Float Literal.
func NewFloatLit(value float64) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.FLOAT,
		Value: strconv.FormatFloat(value, 'g', -1, 64),
	}
}

// NewNil returns a Go AST identity that can be used to represent "nil".
func NewNil() *goast.Ident {
	return NewIdent("nil")
}

// NewUnaryExpr creates a new Go unary expression. You should use this function
// instead of instantiating the UnaryExpr directly because this function has
// extra error checking.
func NewUnaryExpr(operator token.Token, right goast.Expr) *goast.UnaryExpr {
	if right == nil {
		panic("right is nil")
	}

	return &goast.UnaryExpr{
		Op: operator,
		X:  right,
	}
}

// IsGoKeyword will return true if a word is one of the reserved words in Go.
// This means that it cannot be used as an identifier, function name, etc.
//
// The list of reserved words has been taken from the spec at
// https://golang.org/ref/spec#Keywords
func IsGoKeyword(w string) bool {
	switch w {
	case "break", "default", "func", "interface", "select", "case", "defer",
		"go", "map", "struct", "chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type", "continue", "for",
		"import", "return", "var", "_", "init":
		return true
	}
	if IsGoPackage(w) {
		return true
	}

	return false
}

// IsGoPackage Go packages
func IsGoPackage(w string) bool {
	switch w {
	case "fmt", "os", "math", "testing", "unsafe", "ioutil":
		return true
	}
	return false
}

// ConvertFunctionNameFromCtoGo - convert function name fromC to Go
func ConvertFunctionNameFromCtoGo(name string) string {
	if name == "_" {
		return "__"
	}
	return name
}

// GetUintptr - return uintptr. Example : uintptr(unsafe.Pointer(&expr))
func GetUintptr(expr goast.Expr) goast.Expr {

	if _, ok := expr.(*goast.SelectorExpr); ok {
		expr = &goast.IndexExpr{
			X:     expr,
			Index: goast.NewIdent("0"),
		}
	}

	if sl, ok := expr.(*goast.SliceExpr); ok {
		// from :
		//
		// 88  0: *ast.SliceExpr {
		// 89  .  X: *ast.Ident {
		// 91  .  .  Name: "b"
		// 93  .  }
		// 95  .  Low: *ast.BasicLit { ... }
		// 99  .  }
		// 102  }
		//
		// to:
		//
		// 0  *ast.IndexExpr {
		// 1  .  X: *ast.Ident {
		// 3  .  .  Name: "b"
		// 4  .  }
		// 6  .  Index: *ast.BasicLit { ... }
		// 12  }
		expr = &goast.IndexExpr{
			X:     sl.X,
			Index: sl.Low,
		}
	}

	if sl, ok := expr.(*goast.SliceExpr); ok {
		if c, ok := sl.X.(*goast.CallExpr); ok {
			if fin, ok := c.Fun.(*goast.Ident); ok && strings.Contains(fin.Name, "100000") {
				if len(c.Args) == 1 {
					if cc, ok := c.Args[0].(*goast.CallExpr); ok {
						if fin, ok := cc.Fun.(*goast.Ident); ok && strings.Contains(fin.Name, "unsafe.Pointer") {
							if len(cc.Args) == 1 {
								if un, ok := cc.Args[0].(*goast.UnaryExpr); ok && un.Op == token.AND {
									expr = un.X
								}
							}
						}
					}
				}
			}
		}
	}

	return NewCallExpr("uintptr",
		NewCallExpr("unsafe.Pointer",
			&goast.UnaryExpr{
				Op: token.AND,
				X:  expr,
			},
		),
	)
}

// GetUintptrForSlice - return uintptr for slice
// Example : uint64(uintptr(unsafe.Pointer((*(**int)(unsafe.Pointer(& ...slice... )))))))
func GetUintptrForSlice(expr goast.Expr, sizeof int) (goast.Expr, string) {
	returnType := "long long"

	return &goast.BinaryExpr{
		X: NewCallExpr("int64", NewCallExpr("uintptr", NewCallExpr("unsafe.Pointer",
			&goast.StarExpr{
				Star: 1,
				X: &goast.CallExpr{
					Fun: &goast.ParenExpr{
						Lparen: 1,
						X: &goast.StarExpr{
							Star: 1,
							X: &goast.StarExpr{
								Star: 1,
								X:    goast.NewIdent("byte"),
							},
						},
					},
					Lparen: 1,
					Args: []goast.Expr{&goast.CallExpr{
						Fun:    goast.NewIdent("unsafe.Pointer"),
						Lparen: 1,
						Args: []goast.Expr{&goast.UnaryExpr{
							Op: token.AND,
							X:  expr,
						}},
					}},
				},
			},
		))),
		Op: token.QUO,
		Y:  NewCallExpr("int64", goast.NewIdent(fmt.Sprintf("%d", sizeof))),
	}, returnType
}

// CreateSliceFromReference - create a slice, like :
// (*[1]int)(unsafe.Pointer(&a))[:]
func CreateSliceFromReference(goType string, expr goast.Expr) goast.Expr {
	// If the Go type is blank it means that the C type is 'void'.
	if goType == "" {
		goType = "interface{}"
	}

	// This is a hack to convert a reference to a variable into a slice that
	// points to the same location. It will look similar to:
	//
	//     (*[1]int)(unsafe.Pointer(&a))[:]
	//
	// You must always call this Go before using CreateSliceFromReference:
	//
	//     p.AddImport("unsafe")
	//
	return &goast.SliceExpr{
		X: NewCallExpr(fmt.Sprintf("(*[100000000]%s)", goType),
			NewCallExpr("unsafe.Pointer",
				&goast.UnaryExpr{
					X:  expr,
					Op: token.AND,
				}),
		),
	}
}

// NewFuncType - create a new function type, example:
// func ...(fieldList)(returnType)
func NewFuncType(fieldList *goast.FieldList, returnType string, addDefaultReturn bool) *goast.FuncType {
	returnTypes := []*goast.Field{}
	if returnType != "" {
		field := goast.Field{
			Type: NewTypeIdent(returnType),
		}
		if addDefaultReturn {
			field.Names = []*goast.Ident{NewIdent("c4goDefaultReturn")}
		}
		returnTypes = append(returnTypes, &field)
	}

	return &goast.FuncType{
		Params: fieldList,
		Results: &goast.FieldList{
			List: returnTypes,
		},
	}
}

// NewGoExpr is used to facilitate the creation of go AST
//
// It should not be used to transpile user code.
func NewGoExpr(expr string) goast.Expr {
	e, err := parser.ParseExpr(expr)
	if err != nil {
		panic("programming error: " + expr)
	}

	return e
}

// NewAnonymousFunction - create a new anonymous function.
// Example:
// func() returnType{
//		defer func(){
//			deferBody
//		}()
// 		body
//		return returnValue
// }
func NewAnonymousFunction(body, deferBody []goast.Stmt,
	returnValue goast.Expr,
	returnType string) *goast.CallExpr {

	if deferBody != nil {
		body = append(body, []goast.Stmt{&goast.DeferStmt{
			Defer: 1,
			Call: &goast.CallExpr{
				Fun: &goast.FuncLit{
					Type: &goast.FuncType{},
					Body: &goast.BlockStmt{List: deferBody},
				},
				Lparen: 1,
			},
		}}...)
	}

	return &goast.CallExpr{Fun: &goast.FuncLit{
		Type: &goast.FuncType{
			Results: &goast.FieldList{List: []*goast.Field{
				{Type: goast.NewIdent(returnType)},
			}},
		},
		Body: &goast.BlockStmt{
			List: append(body, &goast.ReturnStmt{
				Results: []goast.Expr{returnValue},
			}),
		},
	}}
}

// ConvertToUnsigned - return convertion from signed to unsigned type
// s := func() uint {
//		var x int64
//		x = -1
//		if x < 0 {
//			x += 5
//		}
//		return uint(x)
//	}()
func ConvertToUnsigned(expr goast.Expr, returnType string) goast.Expr {

	varName := "c4go_temp_name"
	// maxValue := "123123123"

	if u, ok := expr.(*goast.CallExpr); ok {
		if i, ok := u.Fun.(*goast.Ident); ok {
			switch i.Name {
			case "uint32", "uint", "uint64":
				expr = u.Args[0]
			}
		}
	}

	return &goast.CallExpr{Fun: &goast.FuncLit{
		Type: &goast.FuncType{
			Results: &goast.FieldList{List: []*goast.Field{
				{Type: goast.NewIdent(returnType)},
			}},
		},
		Body: &goast.BlockStmt{
			List: []goast.Stmt{
				&goast.DeclStmt{
					Decl: &goast.GenDecl{
						Tok: token.VAR,
						Specs: []goast.Spec{
							&goast.ValueSpec{
								Names: []*goast.Ident{goast.NewIdent(varName)},
								Type:  goast.NewIdent("int64"),
							},
						},
					},
				},
				&goast.AssignStmt{
					Lhs: []goast.Expr{goast.NewIdent(varName)},
					Tok: token.ASSIGN,
					Rhs: []goast.Expr{expr},
				},
				// &goast.IfStmt{
				// 	Cond: &goast.BinaryExpr{
				// 		X:  goast.NewIdent(varName),
				// 		Op: token.LSS, // <
				// 		Y:  goast.NewIdent("0"),
				// 	},
				// 	Body: &goast.BlockStmt{
				// 		List: []goast.Stmt{
				// 			&goast.AssignStmt{
				// 				Lhs: []goast.Expr{goast.NewIdent(varName)},
				// 				Tok: token.ADD_ASSIGN, // +=
				// 				Rhs: []goast.Expr{goast.NewIdent(maxValue)},
				// 			},
				// 		},
				// 	},
				// },
				&goast.ReturnStmt{
					Results: []goast.Expr{
						&goast.CallExpr{
							Fun:  goast.NewIdent(returnType),
							Args: []goast.Expr{goast.NewIdent(varName)},
						},
					},
				},
			},
		},
	}}
}
