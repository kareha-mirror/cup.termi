//go:build windows

package termi

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func isWindowsTerminal() bool {
	return os.Getenv("WT_SESSION") != "" || os.Getenv("KAKIKO_RUNNING") != ""
}

var windowsTerminal = isWindowsTerminal()

var stopEvent windows.Handle

var readThread windows.Handle

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procReadConsoleInputW = kernel32.NewProc("ReadConsoleInputW")
)

type KeyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

const EventTypeKey = 0x0001

type InputRecord struct {
	EventType uint16
	_         uint16
	KeyEvent  KeyEventRecord
}

func ReadConsoleInput(
	h windows.Handle,
	rec *InputRecord,
	length uint32,
	read *uint32,
) error {
	r1, _, e := procReadConsoleInputW.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(rec)),
		uintptr(length),
		uintptr(unsafe.Pointer(read)),
	)
	if r1 == 0 {
		return e
	}
	return nil
}

func initRead() error {
	var err error
	stopEvent, err = windows.CreateEvent(nil, 1, 0, nil)
	return err
}

func stopRead() error {
	return windows.SetEvent(stopEvent)
}

func finishRead() error {
	return windows.CloseHandle(stopEvent)
}

var readBuf [1]byte

func read() (byte, bool, error) {
	if windowsTerminal {
		for {
			n, err := os.Stdin.Read(readBuf[:])
			if err != nil {
				return 0, false, err
			}
			if n == 1 {
				b := readBuf[0]
				checkEscape(b)
				return b, true, nil
			}
		}
	} else {
		for {
			var rec InputRecord
			var n uint32
			err := ReadConsoleInput(windows.Handle(os.Stdin.Fd()), &rec, 1, &n)
			if err != nil {
				return 0, false, err
			}
			if n == 1 {
				if rec.EventType != EventTypeKey {
					return 0, false, nil
				}
				if rec.KeyEvent.KeyDown == 0 {
					return 0, false, nil
				}
				if rec.KeyEvent.UnicodeChar == 0 {
					return 0, false, nil
				}
				b := byte(rec.KeyEvent.UnicodeChar)
				checkEscape(b)
				return b, true, nil
			}
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
				break loop
			case windows.WAIT_OBJECT_0 + 1:
				b, ok, err := read()
				if err != nil {
					break loop
				}
				if !ok {
					continue
				}
				readCh <- b
			}
		}
		close(readDone)
		close(readCh)
	}()
}
