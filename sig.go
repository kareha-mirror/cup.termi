package termi

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Sig int

const (
	SigStop Sig = iota
	SigCont
)

var signalWG sync.WaitGroup
var signalCh chan os.Signal
var signalDone chan struct{}
var sigCh chan Sig

func StartSig() {
	signalCh = make(chan os.Signal, 1)
	signalDone = make(chan struct{})
	sigCh = make(chan Sig, 32)

	signal.Notify(signalCh, syscall.SIGTSTP, syscall.SIGCONT)

	signalWG.Add(1)
	go func() {
		defer signalWG.Done()
	loop:
		for {
			select {
			case s := <-signalCh:
				switch s {
				case syscall.SIGTSTP:
					sigCh <- SigStop
				case syscall.SIGCONT:
					sigCh <- SigCont
				}
			case <-signalDone:
				break loop
			}
		}
		signal.Stop(signalCh)
		signal.Reset(syscall.SIGTSTP, syscall.SIGCONT)
		close(signalCh)
		close(sigCh)
	}()
}

func StopSig() {
	close(signalDone)
	signalWG.Wait()
}

func Sigs() chan Sig {
	return sigCh
}
