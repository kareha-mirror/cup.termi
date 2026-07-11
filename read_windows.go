//go:build windows

package termi

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"golang.org/x/sys/windows"
)

func isWindowsTerminal() bool {
	return os.Getenv("WT_SESSION") != ""
}

var windowsTerminal = isWindowsTerminal()

var stopEvent windows.Handle

var readThread windows.Handle

func initRead() error {
	if windowsTerminal {
		var err error
		stopEvent, err = windows.CreateEvent(nil, 1, 0, nil)
		return err
	} else {
		return nil
	}
}

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

func stopRead() error {
	if windowsTerminal {
		return windows.SetEvent(stopEvent)
	} else {
		return CancelSynchronousIo(readThread)
	}
}

func finishRead() error {
	if windowsTerminal {
		return windows.CloseHandle(stopEvent)
	} else {
		return nil
	}
}

var readBuf [1]byte

func read() (byte, error) {
	for {
		var n uint32
		err := windows.ReadFile(
			windows.Handle(os.Stdin.Fd()), readBuf[:], &n, nil,
		)
		fmt.Fprintf(os.Stderr, "(%d,%v)", n, err) // XXX debug
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

func spawnReaderForWindowsTerminal() {
	keyWG.Add(1)
	go func() {
		defer keyWG.Done()

		input := windows.Handle(os.Stdin.Fd())
		handles := []windows.Handle{input, stopEvent}

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
				b, err := read()
				if err != nil {
					break loop
				}
				readCh <- b
			case windows.WAIT_OBJECT_0 + 1:
				break loop
			}
		}
		close(readDone)
		close(readCh)
	}()
}

func spawnReaderForOtherTerminals() {
	keyWG.Add(1)
	go func() {
		defer keyWG.Done()

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		h, err := windows.OpenThread(
			windows.THREAD_TERMINATE,
			false,
			windows.GetCurrentThreadId(),
		)
		if err != nil {
			// XXX inform error
			close(readDone)
			close(readCh)
			return
		}
		defer windows.CloseHandle(h)

		readThread = h

		for {
			b, err := read()
			if err != nil {
				break
			}
			readCh <- b
		}
		close(readDone)
		close(readCh)
	}()
}

func spawnReader() {
	if windowsTerminal {
		spawnReaderForWindowsTerminal()
	} else {
		spawnReaderForOtherTerminals()
	}
}
