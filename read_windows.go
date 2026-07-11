//go:build windows

package termi

import (
	"fmt"
	"os"
	"sync/atomic"
	"syscall"

	"golang.org/x/sys/windows"
)

var stopping = atomic.Bool{}

var stopEvent windows.Handle

var readThread windows.Handle

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procCancelSynchronousIo = kernel32.NewProc("CancelSynchronousIo")
)

func CancelSynchronousIo(h windows.Handle) error {
	r1, _, e1 := procCancelSynchronousIo.Call(uintptr(h))
	if r1 == 0 {
		if e1 != windows.ERROR_SUCCESS {
			return e1
		}
		return syscall.EINVAL
	}
	return nil
}

func initRead() error {
	stopping.Store(false)
	var err error
	stopEvent, err = windows.CreateEvent(nil, 1, 0, nil)
	return err
}

func stopRead() error {
	stopping.Store(true)
	return windows.SetEvent(stopEvent)
}

func finishRead() error {
	return windows.CloseHandle(stopEvent)
}

var readBuf [1]byte

func read() (byte, error) {
	for {
		n, err := os.Stdin.Read(readBuf[:])
		if err != nil {
			return 0, err
		}
		if n == 1 {
			b := readBuf[0]
			checkEscape(b)
			return b, nil
		}
	}
}

func spawnReader() {
	keyWG.Add(1)
	go func() {
		defer keyWG.Done()

		input := windows.Handle(os.Stdin.Fd())
		handles := []windows.Handle{stopEvent, input}

	loop:
		for {
			n, err := windows.WaitForMultipleObjects(
				handles, false, windows.INFINITE,
			)
			if err != nil {
				break
			}
			switch n {
			case windows.WAIT_OBJECT_0:
				fmt.Fprint(os.Stderr, "(stopEvent)") // XXX
				break loop
			case windows.WAIT_OBJECT_0 + 1:
				fmt.Fprint(os.Stderr, "(input)") // XXX
				if stopping.Load() {
					break loop
				}
				b, err := read()
				if err != nil {
					break loop
				}
				readCh <- b
			}
		}
		close(readDone)
		close(readCh)
	}()
}
