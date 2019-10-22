// This file contains transpiling functions for literals and constants. Literals
// are single values like 123 or "hello".

package transpiler

import (
	"bytes"
	"fmt"
	"go/token"
	"regexp"
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

var regexpUnsigned *regexp.Regexp
var regexpLongDouble *regexp.Regexp

func init() {
	regexpUnsigned = regexp.MustCompile(`%(\d+)?u`)
	regexpLongDouble = regexp.MustCompile(`%(\d+)?(.\d+)?lf`)
}

// ConvertToGoFlagFormat convert format flags from C to Go
func ConvertToGoFlagFormat(str string) string {
	// %u to %d
	{
		match := regexpUnsigned.FindAllStringSubmatch(str, -1)
		for _, sub := range match {
			str = strings.Replace(str, sub[0], sub[0][:len(sub[0])-1]+"d", -1)
		}
	}
	// from %lf to %f
	{
		match := regexpLongDouble.FindAllStringSubmatch(str, -1)
		for _, sub := range match {
			str = strings.Replace(str, sub[0], sub[0][:len(sub[0])-2]+"f", -1)
		}
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
	if baseType != "char" &&
		!strings.Contains(baseType, "int") &&
		!strings.Contains(baseType, "wchar_t") {
		if t, ok := p.TypedefType[baseType]; ok {
			n.Type = t
		} else {
			err = fmt.Errorf("type is not valid : `%v`", n.Type)
			p.AddMessage(p.GenerateWarningMessage(err, n))
			return
		}
	}
	var s int
	s, err = types.GetAmountArraySize(n.Type, p)
	if !arrayToArray {
		if err != nil {
			expr = util.NewCallExpr("[]byte",
				util.NewStringLit(strconv.Quote(n.Value+"\x00")))
			exprType = "const char *"
			err = nil // ignore error
			return
		}
		buf := bytes.NewBufferString(n.Value + "\x00")
		if buf.Len() < s {
			buf.Write(make([]byte, s-buf.Len()))
		}
		switch {
		case strings.Contains(baseType, "int"), strings.Contains(baseType, "wchar_t"):
			expr = util.NewCallExpr("[]rune",
				util.NewStringLit(strconv.Quote(buf.String())))
			exprType = "const wchar_t *"
		default:
			expr = util.NewCallExpr("[]byte",
				util.NewStringLit(strconv.Quote(buf.String())))
			exprType = "const char *"
		}
		return
	}
	// Example:
	//
	// var sba SBA = SBA{10, func() (b [100]byte) {
	// 	copy(b[:], "qwe")
	// 	return b
	// }()}
	expr = goast.NewIdent(fmt.Sprintf(
		"func() (b [%v]byte) {copy(b[:], %s);return }()",
		s, strconv.Quote(n.Value)))
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

func transpilePredefinedExpr(n *ast.PredefinedExpr, p *program.Program) (
	expr goast.Expr, exprType string, err error) {

	if len(n.Children()) == 1 {
		expr, exprType, _, _, err = transpileToExpr(n.Children()[0], p, false)
		return
	}

	// There are many more.
	err = fmt.Errorf("unknown PredefinedExpr: %s", n.Name)
	return
}

func transpileCompoundLiteralExpr(n *ast.CompoundLiteralExpr, p *program.Program) (goast.Expr, string, error) {
	expr, t, _, _, err := transpileToExpr(n.Children()[0], p, false)
	return expr, t, err
}
