package linux

import (
	"os"
	"testing"
)

func TestAssertFail(t *testing.T) {
	var res int
	osExit = func(code int) {
		res = code
	}
	defer func() {
		osExit = os.Exit
	}()

	_ = AssertFail([]byte(""), []byte(""), 10, []byte(""))
	if res != 134 {
		t.Fatalf("Another result")
	}
}
