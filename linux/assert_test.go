package linux

import (
	"io/ioutil"
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

	var r func() (*os.File, []byte)
	os.Stderr, r = repl(os.Stderr)

	_ = AssertFail([]byte(""), []byte(""), 10, []byte(""))
	if res != 134 {
		t.Fatalf("Another result")
	}

	var b []byte
	os.Stderr, b = r()
	t.Log(string(b))

	if len(b) == 0 {
		t.Fatalf("Haven't error")
	}
}

func repl(o *os.File) (*os.File, func() (*os.File, []byte)) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	name := tmp.Name()

	var io *os.File = o

	return tmp, func() (*os.File, []byte) {
		err = tmp.Close()
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		return io, b
	}
}
