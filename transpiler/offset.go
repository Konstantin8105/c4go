package transpiler

import (
	"bytes"
	"fmt"
	goast "go/ast"
	"io/ioutil"

	"github.com/Konstantin8105/c4go/ast"
	"github.com/Konstantin8105/c4go/program"
)

func transpileOffsetOfExpr(n *ast.OffsetOfExpr, p *program.Program) (
	expr goast.Expr, exprType string, err error) {
	// clang ast haven`t enought information about OffsetOfExpr
	defer func() {
		if err != nil {
			err = fmt.Errorf("cannot transpile OffsetOfExpr. %v", err)
		}
	}()
	// read file
	var dat []byte
	dat, err = ioutil.ReadFile(n.Pos.File)
	if err != nil {
		err = fmt.Errorf("cannot read file. %v", err)
		return
	}

	lines := bytes.Split(dat, []byte("\n"))
	if n.Pos.Line >= len(lines) && n.Pos.LineEnd >= len(lines) {
		err = fmt.Errorf("not correct position of line {%v,%v}. Amount lines %d",
			n.Pos.Line, n.Pos.LineEnd, len(lines))
		return
	}

	var buffer []byte
	if n.Pos.Line != n.Pos.LineEnd {
		buffer = lines[n.Pos.Line-1][n.Pos.Column:n.Pos.ColumnEnd]
	} else {
		// TODO
		fmt.Println("TODO")
	}

	buffer = bytes.TrimSpace(buffer)
	buffer = bytes.Replace(buffer, []byte("\x00"), []byte(""), -1)

	if len(buffer) == 0 {
		err = fmt.Errorf("Buffer is empty")
		return
	}

	// find `(` and `)`
	if buffer[0] != '(' {
		err = fmt.Errorf("Not start from `(` in buffer : `%s`",
			string(buffer))
		return
	}
	var endPosition int
	for i := range buffer {
		if buffer[i] == ')' {
			endPosition = i
			break
		}
	}
	if buffer[endPosition] != ')' {
		err = fmt.Errorf("Not start from `)` in buffer : `%v`",
			string(buffer))
		return
	}
	buffer = buffer[1:endPosition]

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
	expr = &goast.CallExpr{
		Fun: &goast.SelectorExpr{
			X:   goast.NewIdent("unsafe"),
			Sel: goast.NewIdent("Offsetof"),
		},
		Lparen: 1,
		Args: []goast.Expr{
			&goast.SelectorExpr{
				X: &goast.CompositeLit{
					Type:   goast.NewIdent(string(arguments[0])),
					Lbrace: 1,
				},
				Sel: goast.NewIdent(string(arguments[1])),
			},
		},
	}

	// panic("TODO")

	exprType = n.Type
	return
}
