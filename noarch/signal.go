package noarch

import (
	"os"
	"syscall"
)

var signals map[int32]func(int32)

var c chan os.Signal

func init() {
	signals = map[int32]func(int32){}
	c = make(chan os.Signal, 1)
	go func() {
		for ch := range c {
			s := ch.(syscall.Signal)
			Raise(int32(s))
		}
	}()
}

func Raise(code int32) {
	if f, ok := signals[code]; ok {
		f(0)
	}
}

func Signal(code int32, f func(param int32)) (fr func(int32)) {
	fr, _ = signals[code]
	signals[code] = f
	return
}
