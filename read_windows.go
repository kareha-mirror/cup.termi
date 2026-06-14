//go:build windows

package termi

import (
	"os"
)

func startRead() {
	in = os.Stdin
}

func stopRead() {
	// do nothing
}

func finishRead() {
	// do nothing
}

var in *os.File
var inputBuf [1]byte

func read() (byte, error) {
	for {
		n, err := in.Read(inputBuf[:])
		if err != nil {
			return 0, err
		}
		if n == 1 {
			fireEscape(inputBuf[0])
			return inputBuf[0], nil
		}
	}
}
