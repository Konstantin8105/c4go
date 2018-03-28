package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestDiff(t *testing.T) {
	tcs := []struct {
		a string
		b string

		expect string
	}{
		{
			a:      "a",
			b:      "b",
			expect: "*   1 \"a\"                                     \"b\"",
		},
		{
			a:      "a",
			b:      "a",
			expect: "1 \"a\"                                     \"a\"",
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Test %d : %s", i, tc.a+tc.b), func(t *testing.T) {
			act := ShowDiff(tc.a, tc.b)
			if strings.TrimSpace(act) != strings.TrimSpace(tc.expect) {
				t.Errorf("Not correct result.\nExpected:%s\nActual:%s\n",
					tc.expect, act)
				t.Errorf("Length\nExpected:%d\nActual:%d\n",
					len(tc.expect), len(act))
			}
		})
	}
}
