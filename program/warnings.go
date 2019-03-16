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
		message += fmt.Sprintf("(%T): %s:", n, n.Position().GetSimpleLocation())
	}
	message += fmt.Sprintf("%s", e.Error())
	message = PathSimplification(message)
	return message
}

func PathSimplification(message string) string {
	message = strings.Replace(message, os.Getenv("GOPATH"), "$GOPATH", -1)
	message = strings.Replace(message, "$GOPATH/src/github.com/Konstantin8105/c4go", "C4GO", -1)
	return message
}
