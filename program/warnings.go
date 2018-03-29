package program

import (
	"fmt"

	"github.com/Konstantin8105/c4go/ast"
)

// GenerateErrorMessage - generate error message
func (p *Program) GenerateErrorMessage(e error, n ast.Node) string {
	message := "// Error "
	if e == nil {
		return ""
	}
	if n != nil {
		message += fmt.Sprintf("(%T): %s:", n, n.Position().GetSimpleLocation())
	}
	message += fmt.Sprintf("%s", e.Error())
	return message
}

// GenerateWarningMessage - generate warning message
func (p *Program) GenerateWarningMessage(e error, n ast.Node) string {
	message := "// Warning "
	if e == nil {
		return ""
	}
	if n != nil {
		message += fmt.Sprintf("(%T): %s:", n, n.Position().GetSimpleLocation())
	}
	message += fmt.Sprintf("%s", e.Error())
	return message
}

// GenerateWarningOrErrorMessage - generate error if it happen
func (p *Program) GenerateWarningOrErrorMessage(e error, n ast.Node, isError bool) string {
	if isError {
		return p.GenerateErrorMessage(e, n)
	}

	return p.GenerateWarningMessage(e, n)
}
