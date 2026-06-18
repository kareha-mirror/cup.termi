package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func drawKey(key termi.Key) {
	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)

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
}

func keyMain() {
	for {
		select {
		case key := <-termi.Keys():
			if key.Kind == termi.KeyRune {
				switch key.Rune {
				case 'q':
					return
				case '\x1a': // Ctrl-Z
					termi.Suspend()
				}
			}
			drawKey(key)
		case sig := <-termi.Sigs():
			if sig == termi.SigStop {
				termi.StopInput()
				fmt.Print(termi.ResetAlternate)
				termi.Cooked()
				fmt.Print(termi.ShowCursor)

				termi.ForceSuspend()
				for {
					sig := <-termi.Sigs()
					if sig == termi.SigCont {
						fmt.Print(termi.HideCursor)
						termi.Raw()
						fmt.Print(termi.SetAlternate)
						termi.StartInput()
						break
					}
				}
			}
		}
	}
}
