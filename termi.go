package termi

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var state *term.State

func Raw() {
	if state != nil {
		panic("invalid state")
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
	err := term.Restore(int(os.Stdin.Fd()), state)
	if err != nil {
		panic(err)
	}
	state = nil
}

const (
	SetAlternate   = "\x1b[?1049h"
	ResetAlternate = "\x1b[?1049l"

	Clear      = "\x1b[2J"
	HomeCursor = "\x1b[H"

	HideCursor = "\x1b[?25l"
	ShowCursor = "\x1b[?25h"

	SetInvert   = "\x1b[7m"
	ResetInvert = "\x1b[27m"

	SaveCursor = "\x1b[s"
	LoadCursor = "\x1b[u"

	ScrollReset = "\x1b[r"

	ClearTail = "\x1b[K"
)

func MoveCursor(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
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

func ScrollRange(top int, size int) string {
	return fmt.Sprintf("\x1b[%d;%dr", top+1, top+size)
}
