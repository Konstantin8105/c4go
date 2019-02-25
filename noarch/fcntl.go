package noarch

import "golang.org/x/sys/unix"

type Flock = unix.Flock_t

func Open(pathname []byte, flags int, mode ...int) int {
	if len(mode) == 0 {
		mode = append(mode, 644)
	}
	fd, err := unix.Open(CStringToString(pathname), flags, uint32(mode[0]))
	if err != nil {
		return -1
	}
	return fd
}
