package noarch

import "golang.org/x/sys/unix"

func OpenM(pathname []byte, flags int, mode int) int {
	fd, err := unix.Open(CStringToString(pathname), flags, uint32(mode))
	if err != nil {
		return -1
	}
	return fd
}

func Open(pathname []byte, flags int) int {
	return OpenM(pathname, flags, 644)
}
