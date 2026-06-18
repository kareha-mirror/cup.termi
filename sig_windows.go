//go:build windows

package termi

func StartSig() {
	sigCh = make(chan Sig, 32)
}

func StopSig() {
	close(sigCh)
}

func Suspend() {
	// do nothing
}

func ForceSuspend() {
	// do nothing
}
