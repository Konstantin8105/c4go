package types_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
)

type resolveTestCase struct {
	cType  string
	goType string
}

var resolveTestCases = []resolveTestCase{
	{"int", "int"},
	{"char *[13]", "[][]byte"},
	{"__uint16_t", "uint16"},
	{"void *", "interface{}"},
	{"unsigned short int", "uint16"},
	{"div_t", "noarch.DivT"},
	{"ldiv_t", "noarch.LdivT"},
	{"lldiv_t", "noarch.LldivT"},
	{"int [2]", "[]int"},
	{"int [2][3]", "[][]int"},
	{"int [2][3][4]", "[][][]int"},
	{"int [2][3][4][5]", "[][][][]int"},
	{"int (*[2])(int, int)", "[2]func(int,int)(int)"},
	{"int (*(*(*)))(int, int)", "[][]func(int,int)(int)"},
}

func TestResolve(t *testing.T) {
	p := program.NewProgram()

	for i, testCase := range resolveTestCases {
		t.Run(fmt.Sprintf("Test %d : %s", i, testCase.cType), func(t *testing.T) {
			goType, err := types.ResolveType(p, testCase.cType)
			if err != nil {
				t.Fatal(err)
			}

			goType = strings.Replace(goType, " ", "", -1)
			testCase.goType = strings.Replace(testCase.goType, " ", "", -1)

			if goType != testCase.goType {
				t.Errorf("Expected '%s' -> '%s', got '%s'",
					testCase.cType, testCase.goType, goType)
			}
		})
	}
}

func TestResolveError(t *testing.T) {
	tcs := []string{"w:w", "", "const"}
	for i, tc := range tcs {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			var p program.Program
			if _, err := types.ResolveType(&p, tc); err == nil {
				t.Fatalf("Not acceptable")
			}
		})
	}
}
