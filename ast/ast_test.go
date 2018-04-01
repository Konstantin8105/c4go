package ast

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Konstantin8105/c4go/util"
)

func formatMultiLine(o interface{}) string {
	s := fmt.Sprintf("%#v", o)
	s = strings.Replace(s, "{", "{\n", -1)
	s = strings.Replace(s, ", ", "\n", -1)

	return s
}

func runNodeTests(t *testing.T, tests map[string]Node) {
	i := 1
	for line, expected := range tests {
		testName := fmt.Sprintf("Example%d", i)
		i++

		t.Run(testName, func(t *testing.T) {
			// Append the name of the struct onto the front. This would make the
			// complete line it would normally be parsing.
			name := reflect.TypeOf(expected).Elem().Name()
			actual, err := Parse(name + " " + line)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("%s", util.ShowDiff(formatMultiLine(expected),
					formatMultiLine(actual)))
			}
			if err != nil {
				t.Errorf("Error parsing %v", err)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	cond := &ConditionalOperator{}
	cond.AddChild(&ImplicitCastExpr{})
	cond.AddChild(&ImplicitCastExpr{})
	s := Atos(cond)
	if len(s) == 0 {
		t.Fatalf("Cannot convert AST tree : %#v", cond)
	}
	lines := strings.Split(s, "\n")
	var amount int
	for _, l := range lines {
		if strings.Contains(l, "ImplicitCastExpr") {
			amount++
		}
	}
	if amount != 2 {
		t.Error("Not correct design of output")
	}
}

var lines = []string{
// c4go ast sqlite3.c | head -5000 | sed 's/^[ |`-]*//' | sed 's/<<<NULL>>>/NullStmt/g' | gawk 'length > 0 {print "`" $0 "`,"}'
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, line := range lines {
			Parse(line)
		}
	}
}

func TestPanicCheck(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("panic for parsing string line is not acceptable")
		}
	}()
	_, err := Parse("Some strange line")
	if err == nil {
		t.Errorf("Haven`t error for strange string line not acceptable")
	}
	// Correct node of AST:
	// GotoStmt 0x7fb9cc1994d8 <line:18893:9, col:14>  'end_getDigits' 0x7fb9cc199490
	// Modify for panic in ast regexp
	//
	n, err := Parse("GotoStmt 0x7fb9cc1994d8 <lin8893:9, col:14> ts' 99490")
	if err == nil {
		t.Errorf("Haven`t error for guarantee panic line not acceptable\n%v",
			Atos(n))
	}
}

func TestNullStmt(t *testing.T) {
	n, err := Parse("NullStmt")
	if n != nil || err != nil {
		t.Errorf("Not acceptable for NullStmt")
	}
}
