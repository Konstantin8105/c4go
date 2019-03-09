package noarch

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func Isatty(fd int32) int32 {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	// TODO need test
	if err != nil {
		return 0
	}
	return 1
}

func CloseOnExec(c int32) {
	unix.CloseOnExec(int(c))
}

func Pipe(p []int32) int32 {
	pi := make([]int, len(p))
	for i := range p {
		pi[i] = int(p[i])
	}
	err := unix.Pipe(pi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return -1
	}
	for i := range pi {
		p[i] = int32(pi[i])
	}
	return 0
}

func Read(fd int32, p []byte, num uint) SsizeT {
	n, err := syscall.Read(int(fd), p)
	if err != nil {
		return SsizeT(-1)
	}
	return SsizeT(n)
}

func Write(fd int32, p []byte, num uint) SsizeT {
	n, err := syscall.Write(int(fd), p)
	if err != nil {
		return SsizeT(-1)
	}
	return SsizeT(n)
}

type SsizeT int32
