package linux

import "testing"

func TestAssertFail(t *testing.T) {
	isTest = true
	defer func() {
		isTest = false
	}()

	_ = AssertFail([]byte(""), []byte(""), 10, []byte(""))
}
