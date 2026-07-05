//go:build unix

package stdio

import (
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

type Stdio struct {
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
}

func Dup() (Stdio, error) {
	stdin, err := unix.Dup(syscall.Stdin)
	if err != nil {
		return Stdio{}, err
	}
	stdout, err := unix.Dup(syscall.Stdout)
	if err != nil {
		return Stdio{}, err
	}
	stderr, err := unix.Dup(syscall.Stderr)
	if err != nil {
		return Stdio{}, err
	}
	return Stdio{
		Stdin:  os.NewFile(uintptr(stdin), "(stdin)"),
		Stdout: os.NewFile(uintptr(stdout), "(stdout)"),
		Stderr: os.NewFile(uintptr(stderr), "(stderr)"),
	}, nil
}

func (s *Stdio) AttachTo(cmd *exec.Cmd) {
	cmd.Stdin = s.Stdin
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr
}

func (s *Stdio) Close() {
	s.Stdin.Close()
	s.Stdout.Close()
	s.Stderr.Close()
}
