package darwin

import "testing"

func TestAssertRtn(t *testing.T) {
	isTest = true
	defer func() {
		isTest = false
	}()

	_ = AssertRtn([]byte(""), []byte(""), 10, []byte(""))
}
