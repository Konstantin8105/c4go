package noarch

var signals map[uint]func(int)

func init() {
	signals = map[uint]func(int){}
}

func Raise(code uint) {
	if f, ok := signals[code]; ok {
		f(0)
	}
}

func Signal(code uint, f func(param int)) (fr func(int)) {
	fr, _ = signals[code]
	signals[code] = f
	return
}
