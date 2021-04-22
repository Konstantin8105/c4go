package program

import (
	"fmt"
	"os"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
)

var WarningMessage string = "// Warning "

// GenerateWarningMessage - generate warning message
func (p *Program) GenerateWarningMessage(e error, n ast.Node) string {
	message := WarningMessage
	if e == nil || len(e.Error()) == 0 {
		return ""
	}
	if n != nil {
		var prefix string
		if fd, ok := n.(*ast.FunctionDecl); ok {
			prefix = fmt.Sprintf("n:%s,t1:%s,t2:%s", fd.Name, fd.Type, fd.Type2)
		}
		if prefix == "" {
			message += fmt.Sprintf("(%T): %s:", n, n.Position().GetSimpleLocation())
		} else {
			message += fmt.Sprintf("(%T): {prefix: %s}. %s:", n, prefix, n.Position().GetSimpleLocation())
		}
	}
	message += fmt.Sprintf("%s", e.Error())
	message = PathSimplification(message)
	return message
}

func PathSimplification(message string) string {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		message = strings.Replace(message, gopath, "GOPATH", -1)
		message = strings.Replace(message, "GOPATH/src/github.com/Konstantin8105/c4go", "C4GO", -1)
	}
	return message
}
