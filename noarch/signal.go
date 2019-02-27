package noarch

import (
	"os"
	"syscall"
)

var signals map[int]func(int)

var c chan os.Signal

func init() {
	signals = map[int]func(int){}
	c = make(chan os.Signal, 1)
	go func() {
		for ch := range c {
			s := ch.(syscall.Signal)
			Raise(int(s))
		}
	}()
}

func Raise(code int) {
	if f, ok := signals[code]; ok {
		f(0)
	}
}

func Signal(code int, f func(param int)) (fr func(int)) {
	fr, _ = signals[code]
	signals[code] = f
	return
}
