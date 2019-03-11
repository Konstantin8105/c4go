package noarch

import (
	"syscall"

	"github.com/pkg/term/termios"
)

type Termios = syscall.Termios

func Tcgetattr(fd int32, t []Termios) int32 {
	if err := termios.Tcgetattr(uintptr(fd), &t[0]); err != nil {
		return -1
	}
	return 0
}

func Tcsetattr(fd int32, opt int32, t []Termios) int32 {
	if err := termios.Tcsetattr(uintptr(fd), uintptr(opt), &t[0]); err != nil {
		return -1
	}
	return 0
}

func Tcsendbreak(fd int32, dur int32) int32 {
	if err := termios.Tcsendbreak(uintptr(fd), uintptr(dur)); err != nil {
		return -1
	}
	return 0
}

func Tcdrain(fd int) int {
	if err := termios.Tcdrain(uintptr(fd)); err != nil {
		return -1
	}
	return 0
}

func Tcflush(fd int, dur int) int {
	if err := termios.Tcflush(uintptr(fd), uintptr(dur)); err != nil {
		return -1
	}
	return 0
}

// TODO in pkg/term/termios
// func Tcflow(fd int, dur int) int {
// 	if err := termios.Tcflow(uintptr(fd), uintptr(dur)); err != nil {
// 		return -1
// 	}
// 	return 0
// }

func Cfmakeraw(t []Termios) {
	termios.Cfmakeraw(&t[0])
}

func Cfgetispeed(t []Termios) uint32 {
	return termios.Cfgetispeed(&t[0])
}

func Cfgetospeed(t []Termios) uint32 {
	return termios.Cfgetospeed(&t[0])
}
