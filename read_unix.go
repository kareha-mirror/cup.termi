//go:build unix

package termi

import (
	"context"
	"errors"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func newInput() *os.File {
	fd, err := unix.Dup(syscall.Stdin)
	if err != nil {
		panic(err)
	}
	err = unix.SetNonblock(fd, true)
	if err != nil {
		panic(err)
	}
	return os.NewFile(uintptr(fd), "(input)")
}

func setBlocking() {
	err := unix.SetNonblock(syscall.Stdin, false)
	if err != nil {
		panic(err)
	}
}

var in *os.File
var inputBuf [1]byte

func read(ctx context.Context) (byte, error) {
	for {
		select {
		case <-ctx.Done():
			in.Close()
			return 0, ctx.Err()
		default:
			n, err := in.Read(inputBuf[:])
			if errors.Is(err, unix.EAGAIN) ||
				errors.Is(err, unix.EWOULDBLOCK) {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			if err != nil {
				return 0, err
			}
			if n == 1 {
				fireEscape(inputBuf[0])
				return inputBuf[0], nil
			}
		}
	}
}
