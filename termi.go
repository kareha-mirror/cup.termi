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

const Clear = "\x1b[2J"

const HomeCursor = "\x1b[H"

func MoveCursor(x, y int) string {
	return fmt.Sprintf("\x1b[%d;%dH", y+1, x+1)
}

const HideCursor = "\x1b[?25l"

const ShowCursor = "\x1b[?25h"

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

const SetInvert = "\x1b[7m"

const ResetInvert = "\x1b[27m"

const SaveCursor = "\x1b[s"

const LoadCursor = "\x1b[u"

func ScrollRange(top int, size int) string {
	return fmt.Sprintf("\x1b[%d;%dr", top+1, top+size)
}

const ScrollReset = "\x1b[r"

const ClearTail = "\x1b[K"
