// This file contains transpiling functions for literals and constants. Literals
// are single values like 123 or "hello".

package transpiler

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"

	goast "go/ast"

	"strconv"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
	"github.com/Konstantin8105/c4go/util"
)

func transpileFloatingLiteral(n *ast.FloatingLiteral) *goast.BasicLit {
	return util.NewFloatLit(n.Value)
}

// ConvertToGoFlagFormat convert format flags from C to Go
func ConvertToGoFlagFormat(str string) string {
	// %u to %d
	if strings.Contains(str, "%u") {
		str = strings.Replace(str, "%u", "%d", -1)
	}
	// from %lf to %f
	if strings.Contains(str, "%lf") {
		str = strings.Replace(str, "%lf", "%f", -1)
	}
	return str
}

func transpileStringLiteral(p *program.Program, n *ast.StringLiteral, arrayToArray bool) (
	expr goast.Expr, exprType string, err error) {

	// Convert format flags
	n.Value = ConvertToGoFlagFormat(n.Value)

	// Example:
	// StringLiteral 0x280b918 <col:29> 'char [30]' lvalue "%0"
	baseType := types.GetBaseType(n.Type)
	if baseType != "char" {
		err = fmt.Errorf("Type is not `char` : %v", n.Type)
		p.AddMessage(p.GenerateWarningMessage(err, n))
		return
	}
	var s int
	s, err = types.GetAmountArraySize(n.Type)
	if !arrayToArray {
		if err != nil {
			expr = util.NewCallExpr("[]byte",
				util.NewStringLit(strconv.Quote(n.Value+"\x00")))
			exprType = "const char *"
			return
		}
		buf := bytes.NewBufferString(n.Value + "\x00")
		if buf.Len() < s {
			buf.Write(make([]byte, s-buf.Len()))
		}
		expr = util.NewCallExpr("[]byte",
			util.NewStringLit(strconv.Quote(buf.String())))
		exprType = "const char *"
		return
	}
	// Example:
	//
	// var sba SBA = SBA{10, func() (b [100]byte) {
	// 	copy(b[:], "qwe")
	// 	return b
	// }()}
	expr = goast.NewIdent(fmt.Sprintf(
		"func() (b [%v]byte) { copy(b[:],\"%s\" );return }()",
		s, n.Value))
	exprType = n.Type
	return
}

func transpileIntegerLiteral(n *ast.IntegerLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.INT,
		Value: n.Value,
	}
}

func transpileCharacterLiteral(n *ast.CharacterLiteral) *goast.BasicLit {
	return &goast.BasicLit{
		Kind:  token.CHAR,
		Value: fmt.Sprintf("%q", n.Value),
	}
}

func transpilePredefinedExpr(n *ast.PredefinedExpr, p *program.Program) (goast.Expr, string, error) {
	// A predefined expression is a literal that is not given a value until
	// compile time.
	//
	// TODO: Predefined expressions are not evaluated
	// https://github.com/Konstantin8105/c4go/issues/81

	switch n.Name {
	case "__PRETTY_FUNCTION__":
		return util.NewCallExpr(
			"[]byte",
			util.NewStringLit(`"void print_number(int *)"`),
		), "const char*", nil

	case "__func__":
		return util.NewCallExpr(
			"[]byte",
			util.NewStringLit(strconv.Quote(p.Function.Name)),
		), "const char*", nil

	default:
		// There are many more.
		panic(fmt.Sprintf("unknown PredefinedExpr: %s", n.Name))
	}
}

func transpileCompoundLiteralExpr(n *ast.CompoundLiteralExpr, p *program.Program) (goast.Expr, string, error) {
	expr, t, _, _, err := transpileToExpr(n.Children()[0], p, false)
	return expr, t, err
}
