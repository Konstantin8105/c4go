package types_test

import (
	"testing"

	"github.com/Konstantin8105/c4go/program"
	"github.com/Konstantin8105/c4go/types"
)

type sizeofTestCase struct {
	cType   string
	size    int
	isError bool
}

var sizeofTestCases = []sizeofTestCase{
	{"int", 4, true},
	{"int [2]", 4 * 2, true},
	{"int [2][3]", 4 * 2 * 3, true},
	{"int [2][3][4]", 4 * 2 * 3 * 4, true},
	{"int *[2]", 8 * 2, true},
	{"int *[2][3]", 8 * 2 * 3, true},
	{"int *[2][3][4]", 8 * 2 * 3 * 4, true},
	{"int *", 8, true},
	{"int **", 8, true},
	{"int ***", 8, true},
	{"char *const", 8, true},
	{"char *const [3]", 24, true},
	{"struct c [2]", 0, true},
}

func TestSizeOf(t *testing.T) {
	p := program.NewProgram()

	for _, testCase := range sizeofTestCases {
		t.Run(testCase.cType, func(t *testing.T) {
			size, err := types.SizeOf(p, testCase.cType)
			if !((err != nil && testCase.isError == false) ||
				(err == nil && testCase.isError == true)) {
				t.Fatal(err)
			}

			if size != testCase.size {
				t.Errorf("Expected '%s' -> '%d', got '%d'",
					testCase.cType, testCase.size, size)
			}
		})
	}
}
