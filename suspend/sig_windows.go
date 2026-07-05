//go:build windows

package suspend

func Init() {
	sigCh = make(chan Sig, 32)
}

func Finish() {
	close(sigCh)
}

func Suspend() {
	// do nothing
}

func ForceSuspend() {
	// do nothing
}
