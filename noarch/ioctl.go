package noarch

import (
	"golang.org/x/sys/unix"
)

type Winsize = unix.Winsize

func Ioctl(fd int32, req int32, w []Winsize) int32 {
	wb, err := unix.IoctlGetWinsize(int(fd), uint(req))
	if err != nil {
		return -1
	}
	w[0] = *wb
	return 0
}
