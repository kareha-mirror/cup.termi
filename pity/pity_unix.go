//go:build unix

package pity

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Pity struct {
	c *exec.Cmd
	f *os.File
}

func Start(cmd string, args ...string) (*Pity, error) {
	c := exec.Command(cmd, args...)
	f, err := pty.Start(c)
	if err != nil {
		return nil, err
	}
	return &Pity{
		c: c,
		f: f,
	}, nil
}

func (p *Pity) SetSize(w, h int) error {
	return pty.Setsize(p.f, &pty.Winsize{
		Rows: uint16(h),
		Cols: uint16(w),
	})
}

func (p *Pity) Read(b []byte) (n int, err error) {
	return p.f.Read(b)
}

func (p *Pity) Write(b []byte) (n int, err error) {
	return p.f.Write(b)
}

func (p *Pity) Close() error {
	return p.f.Close()
}

func (p *Pity) Wait() error {
	err := p.c.Wait()
	if err == nil {
		return nil
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return &ExitError{code: ee.ExitCode()}
	}
	return err
}
