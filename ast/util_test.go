package ast

import (
	"fmt"
	"testing"
)

func TestRemoveQuotes(t *testing.T) {
	tcs := []struct {
		in  string
		out string
	}{
		{"  ", ""},
		{"  \"\"  ", ""},
		{"''", ""},

		{"  \"string\"  ", "string"},
		{"'string'", "string"},
	}
	for index, tc := range tcs {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			a := removeQuotes(tc.in)
			if a != tc.out {
				t.Errorf("Not acceptable : `%v` `%v`", a, tc.out)
			}
		})
	}
}

func TestAtof(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Not acceptable result")
		}
	}()
	_ = atof("Some not float64")
}

func TestUnquote(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Not acceptable result")
		}
	}()
	_ = unquote("\"")
}

func TestTypesTree(t *testing.T) {
	_ = typesTree(nil, 0)
}
