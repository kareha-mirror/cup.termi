//go:build windows

package pity

import (
	"errors"

	"golang.org/x/sys/windows"

	"github.com/charmbracelet/x/conpty"
)

type Pity struct {
	f *conpty.ConPty
	h windows.Handle
}

func Start(cmd string, args ...string) (*Pity, error) {
	f, err := conpty.New(0, 0, 0)
	if err != nil {
		return nil, err
	}
	_, h, err := f.Spawn(cmd, args, nil)
	if err != nil {
		return nil, err
	}
	return &Pity{
		f: f,
		h: windows.Handle(h),
	}, nil
}

func (pty *Pity) SetSize(w, h int) error {
	return pty.f.Resize(w, h)
}

func (p *Pity) Read(b []byte) (n int, err error) {
	return p.f.Read(b)
}

func (p *Pity) Write(b []byte) (n int, err error) {
	return p.f.Write(b)
}

func (p *Pity) Close() error {
	err := windows.CloseHandle(p.h)
	return errors.Join(err, p.f.Close())
}

func (p *Pity) Wait() error {
	if _, err :=
		windows.WaitForSingleObject(p.h, windows.INFINITE); err != nil {
		return err
	}

	var code uint32
	if err := windows.GetExitCodeProcess(p.h, &code); err != nil {
		return err
	}

	if code != 0 {
		return &ExitError{code: int(code)}
	}
	return nil
}
