//go:build windows

package termi

import (
	"os"
	"os/exec"
)

type Stdio struct {
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
}

func DupStdio() (Stdio, error) {
	return Stdio{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}

func (s *Stdio) AttachTo(cmd *exec.Cmd) {
	cmd.Stdin = s.Stdin
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr
}

func (s *Stdio) Close() {
	// do nothing
}
