package noarch

import "golang.org/x/sys/unix"

type Winsize = unix.Winsize

func Ioctl(fd int, req int, w *Winsize) int {
	var err error
	w, err = unix.IoctlGetWinsize(fd, uint(req))
	if err != nil {
		return -1
	}
	return 0
}
