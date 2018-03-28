package types

import (
	"github.com/Konstantin8105/c4go/program"
)

// ResolveTypeForBinaryOperator determines the result Go type when performing a
// binary expression.
func ResolveTypeForBinaryOperator(p *program.Program, operator, leftType, rightType string) string {
	if operator == "==" ||
		operator == "!=" ||
		operator == ">" ||
		operator == ">=" ||
		operator == "<" ||
		operator == "<=" {
		return "bool"
	}

	return leftType
}
