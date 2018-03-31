package program

import (
	"fmt"

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
	//
	// warningMsg:= "// Warning (*ast.BinaryOperator):  /tmp/SQLITE/sqlite-amalgamation-3220000/shell.c:14132 :Cannot casting {char * -> NullPointerType *}. ..."
	// if strings.HasPrefix(message,warningMsg) {
	// 	fmt.Println(ast.Atos(n))
	// panic("found")
	// }

	return message
}
