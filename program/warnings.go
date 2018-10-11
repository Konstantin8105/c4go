package program

import (
	"fmt"
	"os"
	"strings"

	"github.com/Konstantin8105/c4go/ast"
)

// GenerateWarningMessage - generate warning message
func (p *Program) GenerateWarningMessage(e error, n ast.Node) string {
	message := "// Warning "
	if e == nil || len(e.Error()) == 0 {
		return ""
	}
	if n != nil {
		message += fmt.Sprintf("(%T): %s:", n, n.Position().GetSimpleLocation())
	}
	message += fmt.Sprintf("%s", e.Error())

	// example for debugging search ploblem ast node
	// warningMsg := "noarch.InterfaceTo[5]intSlice"
	// if strings.Contains(message, warningMsg) {
	// 	fmt.Println(ast.Atos(n))
	// 	panic("found")
	// }

	message = strings.Replace(message, os.Getenv("GOPATH"), "$GOPATH", -1)
	return message
}
