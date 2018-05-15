package util

import (
	"fmt"
	"testing"
)

func TestUcfirst(t *testing.T) {
	tcs := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "A"},
		{"w", "W"},
		{"wa", "Wa"},
	}

	for index, tc := range tcs {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			a := Ucfirst(tc.in)
			if a != tc.out {
				t.Errorf("Result is not same: `%s` `%s`", a, tc.out)
			}
		})
	}
}
