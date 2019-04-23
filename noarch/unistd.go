package noarch

import (
	"fmt"
	"io"
	"os"

	"syscall"

	"golang.org/x/sys/unix"
)

func Isatty(fd int32) int32 {
	_, err := unix.IoctlGetTermios(int(fd), syscall.TCGETS)
	// TODO need test
	if err != nil {
		return 0
	}
	return 1
}

func CloseOnExec(c int32) {
	syscall.CloseOnExec(int(c))
}

func Pipe(p []int32) int32 {
	pi := make([]int, len(p))
	for i := range p {
		pi[i] = int(p[i])
	}
	err := syscall.Pipe(pi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return -1
	}
	for i := range pi {
		p[i] = int32(pi[i])
	}
	return 0
}

func Read(fd int32, p []byte, num uint32) SsizeT {
	if num == 0 {
		return 0
	}
	p = p[:num]
	var n int
	var err error
	switch fd {
	case 0:
		n, err = os.Stdin.Read(p)
	case 1:
		n, err = os.Stdout.Read(p)
	case 2:
		n, err = os.Stderr.Read(p)
	default:
		n, err = syscall.Read(int(fd), p)
	}
	if err != nil && err != io.EOF {
		return SsizeT(-1)
	}
	return SsizeT(n)
}

func Write(fd int32, p []byte, num uint32) SsizeT {
	p = p[:num]
	n, err := syscall.Write(int(fd), p)
	if err != nil {
		return SsizeT(-1)
	}
	return SsizeT(n)
}

type SsizeT int32

func Ftruncate(fd int32, length int64) int32 {
	err := syscall.Ftruncate(int(fd), length)
	if err != nil {
		return -1
	}
	return 0
}

func Fstat(fd int32, stat []syscall.Stat_t) int32 {
	if len(stat) == 0 {
		return -1
	}
	// func Fstat(fd int, stat *Stat_t) (err error)
	err := syscall.Fstat(int(fd), &(stat[0]))
	if err != nil {
		return -1
	}
	return 0
}
