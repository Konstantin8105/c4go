package ast

import (
	"fmt"
	"testing"
)

func TestC4GoErrorNode1(t *testing.T) {
	var c C4goErrorNode
	tcs := []func(){
		func() { c.AddChild(c) },
		func() { _ = c.Address() },
		func() { _ = c.Children() },
		func() { _ = c.Position() },
	}

	for index, tc := range tcs {
		t.Run(fmt.Sprintf("%v", index), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Cannot found panic")
				}
			}()
			tc()
		})
	}
}
