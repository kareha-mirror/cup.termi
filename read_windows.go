//go:build windows

package termi

import (
	"fmt"
	"os"
	"sync/atomic"
)

var readAlive = atomic.Bool{}

func initRead() error {
	readAlive.Store(true)
	return nil
}

func stopRead() error {
	readAlive.Store(false)
	return nil
}

func finishRead() error {
	return nil
}

var readBuf [1]byte

func read() (byte, error) {
	for {
		n, err := os.Stdin.Read(readBuf[:])
		if err != nil {
			return 0, err
		}
		if n == 1 {
			if !readAlive.Load() {
				return 0, fmt.Errorf("read stopped")
			}
			b := readBuf[0]
			checkEscape(b)
			return b, nil
		}
	}
}

func waitKey() {
	//keyWG.Wait() // XXX
}
