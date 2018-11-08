package transpiler

import (
	"fmt"
	"reflect"
	"testing"
	"unicode/utf8"

	goast "go/ast"
	"go/token"

	"github.com/Konstantin8105/c4go/ast"
)

var chartests = []struct {
	in  int    // Integer Character Code
	out string // Output Character Literal
}{
	// NUL byte
	{0, "'\\x00'"},

	// ASCII control characters
	{7, "'\\a'"},
	{8, "'\\b'"},
	{9, "'\\t'"},
	{10, "'\\n'"},
	{11, "'\\v'"},
	{12, "'\\f'"},
	{13, "'\\r'"},

	// printable ASCII
	{32, "' '"},
	{34, "'\"'"},
	{39, "'\\''"},
	{65, "'A'"},
	{92, "'\\\\'"},
	{191, "'¿'"},

	// printable unicode
	{948, "'δ'"},
	{0x03a9, "'Ω'"},
	{0x2020, "'†'"},

	// non-printable unicode
	{0xffff, "'\\uffff'"},
	{utf8.MaxRune, "'\\U0010ffff'"},
}

func TestCharacterLiterals(t *testing.T) {
	for _, tt := range chartests {
		expected := &goast.BasicLit{Kind: token.CHAR, Value: tt.out}
		actual := transpileCharacterLiteral(&ast.CharacterLiteral{Value: tt.in})
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("input: %v", tt.in)
			t.Errorf("  expected: %v", expected)
			t.Errorf("  actual:   %v", actual)
		}
	}
}

func TestFormatFlag(t *testing.T) {
	tcs := []struct {
		in, out string
	}{
		{
			in:  "",
			out: "",
		},
		{
			in:  "%",
			out: "%",
		},
		{
			in:  "%34",
			out: "%34",
		},
		{
			in:  "%5.4",
			out: "%5.4",
		},
		{
			in:  "%5f",
			out: "%5f",
		},
		{
			in:  "%u %u",
			out: "%d %d",
		},
		{
			in:  "%u %2u",
			out: "%d %2d",
		},
		{
			in:  "%u",
			out: "%d",
		},
		{
			in:  "%5u  %2u",
			out: "%5d  %2d",
		},
		{
			in:  "%5u",
			out: "%5d",
		},
		{
			in:  "%lf",
			out: "%f",
		},
		{
			in:  "%12lf",
			out: "%12f",
		},
		{
			in:  "%12.4lf",
			out: "%12.4f",
		},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			act := ConvertToGoFlagFormat(tc.in)
			if act != tc.out {
				t.Fatalf("Not same '%v' != '%v'", act, tc.out)
			}
		})
	}
}
