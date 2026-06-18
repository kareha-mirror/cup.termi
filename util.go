package termi

import (
	"os"
	"strings"
)

func IsTmux() bool {
	return os.Getenv("TMUX") != "" ||
		strings.Contains(os.Getenv("TERM"), "tmux")
}
