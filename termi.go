package termi

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var state *term.State

func Raw() {
	if state != nil {
		term.Restore(int(os.Stdin.Fd()), state)
		state = nil
	}
	s, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	state = s
}

func Cooked() {
	if state == nil {
		panic("invalid state")
	}
	term.Restore(int(os.Stdin.Fd()), state)
	state = nil
}

func Clear() string {
	return "\x1b[2J"
}

func HomeCursor() string {
	return "\x1b[H"
}

func MoveCursor(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

func HideCursor() string {
	return "\x1b[?25l"
}

func ShowCursor() string {
	return "\x1b[?25h"
}

func Size() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w, h = 80, 24
	}
	if w < 1 {
		w = 80
	}
	if h < 1 {
		h = 24
	}
	return w, h
}

func EnableInvert() string {
	return "\x1b[7m"
}

func DisableInvert() string {
	return "\x1b[27m"
}

func SaveCursor() string {
	return "\x1b[s"
}

func LoadCursor() string {
	return "\x1b[u"
}

func ScrollRange(top int, size int) string {
	return fmt.Sprintf("\x1b[%d;%dr", top+1, top+size)
}

func ScrollReset() string {
	return "\x1b[r"
}

func ClearTail() string {
	return "\x1b[K"
}
