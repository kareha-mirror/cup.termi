package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func drawSize() {
	w, h := termi.Size()

	fmt.Print(termi.HideCursor())

	fmt.Print(termi.Clear())
	fmt.Print(termi.HomeCursor())

	fmt.Print(termi.MoveCursor(0, 0))
	fmt.Print("+")
	fmt.Print(termi.MoveCursor(w-1, 0))
	fmt.Print("+")
	fmt.Print(termi.MoveCursor(0, h-1))
	fmt.Print("+")
	fmt.Print(termi.MoveCursor(w-1, h-1))
	fmt.Print("+")

	fmt.Print(termi.ShowCursor())
}

func sizeMain() {
	for {
		drawSize()

		key := termi.ReadKey()
		switch key.Kind {
		case termi.KeyRune:
			switch key.Rune {
			case termi.RuneEscape:
				return
			}
		}
	}
}
