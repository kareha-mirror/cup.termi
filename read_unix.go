//go:build unix

package termi

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

var readWakerR int
var readWakerW int
var readFDs []unix.PollFd
var readBuf [1]byte

func startRead() error {
	var pipe [2]int
	err := unix.Pipe(pipe[:])
	if err != nil {
		return err
	}
	readWakerR = pipe[0]
	readWakerW = pipe[1]
	readFDs = []unix.PollFd{
		{Fd: int32(readWakerR), Events: unix.POLLIN},
		{Fd: int32(syscall.Stdin), Events: unix.POLLIN},
	}
	return nil
}

func stopRead() error {
	poison := []byte{0}
	for {
		n, err := unix.Write(readWakerW, poison)
		if err != nil {
			return err
		}
		if n < 1 {
			continue
		}
		return nil
	}
}

func finishRead() error {
	err := unix.Close(readWakerR)
	if err != nil {
		return err
	}
	return unix.Close(readWakerW)
}

func read() (byte, error) {
	for {
		_, err := unix.Poll(readFDs, -1)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return 0, fmt.Errorf("failed to poll")
		}
		if readFDs[0].Revents&unix.POLLIN != 0 {
			var sink [1]byte
			unix.Read(int(readFDs[0].Fd), sink[:])
			return 0, fmt.Errorf("killed")
		}
		if readFDs[1].Revents&unix.POLLIN != 0 {
			for {
				n, err := unix.Read(int(readFDs[1].Fd), readBuf[:])
				if err != nil {
					return 0, err
				}
				if n < 1 {
					continue
				}
				b := readBuf[0]
				checkEscape(b)
				return b, nil
			}
		}
		return 0, fmt.Errorf("invalid state")
	}
}
