package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func start() {
	termi.Raw()
}

func finish() {
	fmt.Print(termi.Clear())
	fmt.Print(termi.HomeCursor())
	termi.Cooked()
	fmt.Print(termi.ShowCursor())
}

func draw() {
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

	fmt.Print(termi.MoveCursor(0, 2))
	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			idx := 16*j + i
			fmt.Print(termi.Palette(idx).Fg())
			fmt.Printf(" %3d", idx)
		}
		fmt.Print("\r\n")
	}
	fmt.Print(termi.ResetAll)

	fmt.Print(termi.ShowCursor())
}

func mainLoop() {
	for {
		draw()

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

func main() {
	start()
	defer finish()
	mainLoop()
}
