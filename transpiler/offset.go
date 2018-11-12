package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/util"
)

func transpileOffsetOfExpr(n *ast.OffsetOfExpr, p *program.Program) (
	expr goast.Expr, exprType string, err error) {
	// clang ast haven`t enought information about OffsetOfExpr
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpile OffsetOfExpr. %v", err)
		}
	}()

	var buffer []byte
	pos := n.Position()
	buffer, err = p.PreprocessorFile.GetSnippet(pos.File,
		pos.Line, pos.LineEnd,
		pos.Column, pos.ColumnEnd)
	if err != nil {
		err = fmt.Errorf("cannot found snippet position is %v. %v",
			n.Position(), err)
		return
	}

	if len(buffer) == 0 {
		err = fmt.Errorf("Buffer is empty")
		return
	}

	if !bytes.HasPrefix(buffer, []byte("__builtin_offsetof(")) {
		err = fmt.Errorf("Haven`t prefix `__builtin_offsetof(` in buffer `%v`",
			string(buffer))
		return
	}

	buffer = buffer[len("__builtin_offsetof(") : len(buffer)-1]

	// separate by `,`
	arguments := bytes.Split(buffer, []byte(","))
	if len(arguments) != 2 {
		err = fmt.Errorf("Not correct amount of arguments in `%v` found %v",
			string(buffer), len(arguments))
		return
	}

	for i := range arguments {
		arguments[i] = bytes.TrimSpace(arguments[i])
	}

	// preparing name of struct
	if bytes.HasPrefix(arguments[0], []byte("struct ")) {
		arguments[0] = arguments[0][len("struct "):]
	}

	p.AddImport("unsafe")
	expr = util.NewCallExpr("unsafe.Offsetof",
		&goast.SelectorExpr{
			X: &goast.CompositeLit{
				Type:   goast.NewIdent(string(arguments[0])),
				Lbrace: 1,
			},
			Sel: goast.NewIdent(string(arguments[1])),
		})

	exprType = n.Type
	return
}
