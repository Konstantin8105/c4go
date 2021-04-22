package util

import (
	"fmt"
	"testing"
)

func TestPanicIfNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Cannot check panic")
		}
	}()
	PanicIfNil(nil, "Check on nil : PanicIfNil")
}

func TestPanicOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Cannot check panic")
		}
	}()
	PanicOnError(fmt.Errorf("some error"), "Check on nil : PanicOnError")
}
