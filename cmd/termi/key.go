package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
	"tea.kareha.org/cup/termi/suspend"
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
					suspend.Suspend()
				}
			}
			drawKey(key)
		case sig := <-suspend.Sigs():
			if sig == suspend.SigStop {
				termi.FinishKey()
				fmt.Print(termi.ResetAlternate)
				termi.Cooked()
				fmt.Print(termi.ShowCursor)

				suspend.ForceSuspend()
				for {
					sig := <-suspend.Sigs()
					if sig == suspend.SigCont {
						fmt.Print(termi.HideCursor)
						termi.Raw()
						fmt.Print(termi.SetAlternate)
						termi.InitKey()
						break
					}
				}
			}
		}
	}
}
