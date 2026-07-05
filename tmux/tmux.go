package tmux

import (
	"os"
	"strings"
)

func Exists() bool {
	return os.Getenv("TMUX") != "" ||
		strings.Contains(os.Getenv("TERM"), "tmux")
}
