//go:build windows

package termi

import (
	"context"
	"os"
)

func newInput() *os.File {
	return os.Stdin
}

func setBlocking() {
	// do nothing
}

var in *os.File
var inputBuf [1]byte

func read(ctx context.Context) (byte, error) {
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
