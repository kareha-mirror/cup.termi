package pity

import (
	"fmt"
)

type ExitError struct {
	code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit status %d", e.code)
}

func (e *ExitError) ExitCode() int {
	return e.code
}
