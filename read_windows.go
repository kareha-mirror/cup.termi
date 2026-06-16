//go:build windows

package termi

import (
	"os"
)

func startRead() error {
	return nil
}

func stopRead() error {
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
			b := readBuf[0]
			checkEscape(b)
			return b, nil
		}
	}
}
