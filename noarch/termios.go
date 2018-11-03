package noarch

import (
	"syscall"

	"github.com/pkg/term/termios"
)

func Tcgetattr(fd int, t *syscall.Termios) int {
	if err := termios.Tcgetattr(uintptr(fd), t); err != nil {
		return -1
	}
	return 0
}

func Tcsetattr(fd int, opt int, t *syscall.Termios) int {
	if err := termios.Tcsetattr(uintptr(fd), uintptr(opt), t); err != nil {
		return -1
	}
	return 0
}

func Tcsendbreak(fd int, dur int) int {
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
