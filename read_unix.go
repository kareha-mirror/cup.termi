//go:build unix

package termi

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

var wakerR int
var wakerW int
var fds []unix.PollFd
var inputBuf [1]byte

func startRead() {
	var pipe [2]int
	err := unix.Pipe(pipe[:])
	if err != nil {
		panic(err)
	}
	wakerR = pipe[0]
	wakerW = pipe[1]

	fds = []unix.PollFd{
		{Fd: int32(syscall.Stdin), Events: unix.POLLIN},
		{Fd: int32(wakerR), Events: unix.POLLIN},
	}
}

func stopRead() {
	for {
		n, err := unix.Write(wakerW, []byte{1})
		if err != nil {
			panic(err)
		}
		if n == 1 {
			return
		}
	}
}

func finishRead() {
	unix.Close(wakerR)
	unix.Close(wakerW)
}

func read() (byte, error) {
	for {
		_, err := unix.Poll(fds, -1)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return 0, fmt.Errorf("failed to poll")
		}
		if fds[1].Revents&unix.POLLIN != 0 {
			var b [16]byte
			unix.Read(int(fds[1].Fd), b[:])
			return 0, fmt.Errorf("killed")
		}
		if fds[0].Revents&unix.POLLIN != 0 {
			for {
				n, err := os.Stdin.Read(inputBuf[:])
				if err != nil {
					return 0, err
				}
				if n == 1 {
					fireEscape(inputBuf[0])
					return inputBuf[0], nil
				}
			}
		}
		return 0, fmt.Errorf("invalid state")
	}
}
