//go:build windows

package termi

import (
	"os"

	"golang.org/x/sys/windows"
)

var stopEvent windows.Handle

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

func waitKey() {
	keyWG.Wait()
}

var readThread windows.Handle

func spawnReader() {
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
