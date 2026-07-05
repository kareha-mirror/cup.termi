//go:build unix

package suspend

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var signalWG sync.WaitGroup
var signalCh chan os.Signal
var signalDone chan struct{}

func Init() {
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

func Finish() {
	close(signalDone)
	signalWG.Wait()
}

func Suspend() {
	syscall.Kill(syscall.Getpid(), syscall.SIGTSTP)
}

func ForceSuspend() {
	syscall.Kill(syscall.Getpid(), syscall.SIGSTOP)
}
