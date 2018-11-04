package noarch

import (
	"golang.org/x/sys/unix"
)

type Winsize = unix.Winsize

func Ioctl(fd int, req int, w []Winsize) int {
	wb, err := unix.IoctlGetWinsize(fd, uint(req))
	if err != nil {
		return -1
	}
	w[0] = *wb
	return 0
}
