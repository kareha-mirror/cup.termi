package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func drawKey(key termi.Key) {
	fmt.Print(termi.HideCursor())

	fmt.Print(termi.Clear())
	fmt.Print(termi.HomeCursor())

	kindName := ""
	switch key.Kind {
	case termi.KeyRune:
		kindName = "rune"

	case termi.KeyUp:
		kindName = "up"
	case termi.KeyDown:
		kindName = "down"
	case termi.KeyRight:
		kindName = "right"
	case termi.KeyLeft:
		kindName = "left"

	case termi.KeyBeginPaste:
		kindName = "begin paste"
	case termi.KeyEndPaste:
		kindName = "end paste"

	case termi.KeyUnknown:
		kindName = "unknown"
	}

	fmt.Printf("kind = %s\r\n", kindName)
	fmt.Printf("rune = 0x%08x\r\n", key.Rune)
	fmt.Printf("raw = %v\r\n", []byte(key.Raw))

	fmt.Print(termi.ShowCursor())
}

func keyMain() {
	for {
		key := termi.ReadKey()

		drawKey(key)

		switch key.Kind {
		case termi.KeyRune:
			switch key.Rune {
			case termi.RuneEscape:
				return
			}
		}
	}
}
