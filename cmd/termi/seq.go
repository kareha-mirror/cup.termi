package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func drawSeq(seq termi.Seq) {
	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)

	kindName := ""
	switch seq.Kind {
	case termi.SeqRune:
		kindName = "rune"

	case termi.SeqUp:
		kindName = "up"
	case termi.SeqDown:
		kindName = "down"
	case termi.SeqRight:
		kindName = "right"
	case termi.SeqLeft:
		kindName = "left"

	case termi.SeqBeginPaste:
		kindName = "begin paste"
	case termi.SeqEndPaste:
		kindName = "end paste"

	case termi.SeqUnknown:
		kindName = "unknown"
	}

	fmt.Printf("kind = %s\r\n", kindName)
	fmt.Printf("rune = 0x%08x\r\n", seq.Rune)
	fmt.Printf("raw = %v\r\n", []byte(seq.Raw))
}

func seqMain() {
	for {
		seq := termi.ReadSeq()
		if seq.Kind == termi.SeqQuit {
			return
		}

		drawSeq(seq)

		if seq.Kind == termi.SeqRune && seq.Rune == 'q' {
			return
		}
	}
}
