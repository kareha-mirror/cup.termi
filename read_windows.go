//go:build windows

package termi

import (
	"os"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var stopEvent windows.Handle

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

var surrogatePending uint16
var surrogateHasHigh bool

func isHighSurrogate(r uint16) bool {
	return 0xd800 <= r && r <= 0xdbff
}

func isLowSurrogate(r uint16) bool {
	return 0xdc00 <= r && r <= 0xdfff
}

func read() (Key, bool, error) {
	for {
		var rec InputRecord
		var n uint32
		err := ReadConsoleInput(windows.Handle(os.Stdin.Fd()), &rec, 1, &n)
		if err != nil {
			return Key{}, false, err
		}
		if n == 1 {
			if rec.EventType != EventTypeKey {
				return Key{}, false, nil
			}
			if rec.KeyEvent.KeyDown == 0 {
				return Key{}, false, nil
			}
			c := rec.KeyEvent.UnicodeChar
			if c == 0 {
				return Key{}, false, nil
			}
			if isHighSurrogate(c) {
				surrogatePending = c
				surrogateHasHigh = true
				return Key{}, false, nil
			}
			if surrogateHasHigh && isLowSurrogate(c) {
				r := utf16.DecodeRune(rune(surrogatePending), rune(c))
				checkEscape(r)
				return Key{KeyRune, r, ""}, true, nil
			} else {
				r := rune(c)
				checkEscape(r)
				return Key{KeyRune, r, ""}, true, nil
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
