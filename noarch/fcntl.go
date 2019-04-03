package noarch

import "syscall"

type Flock = syscall.Flock_t

func Open(pathname []byte, flags int32, mode ...int32) int32 {
	if len(mode) == 0 {
		mode = append(mode, 644)
	}
	fd, err := syscall.Open(CStringToString(pathname), int(flags), uint32(mode[0]))
	if err != nil {
		return -1
	}
	return int32(fd)
}
