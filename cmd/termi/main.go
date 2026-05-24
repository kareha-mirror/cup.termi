package main

import (
	"tea.kareha.org/cup/termi"
)

func start() {
	termi.Raw()
}

func finish() {
	termi.Clear()
	termi.HomeCursor()
	termi.Cooked()
	termi.ShowCursor()
}

func draw() {
	w, h := termi.Size()

	termi.HideCursor()

	termi.Clear()
	termi.HomeCursor()

	termi.MoveCursor(0, 0)
	termi.Print("+")
	termi.MoveCursor(w-1, 0)
	termi.Print("+")
	termi.MoveCursor(0, h-1)
	termi.Print("+")
	termi.MoveCursor(w-1, h-1)
	termi.Print("+")

	termi.MoveCursor(1, 2)
	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			idx := 16*j + i
			termi.SetFgColor(termi.Palette[idx])
			termi.Printf(" %3d", idx)
		}
		termi.Print("\r\n")
	}
	termi.ResetColor()

	termi.ShowCursor()
}

func mainLoop() {
	for {
		key := termi.ReadKey()
		switch key.Kind {
		case termi.KeyRune:
			switch key.Rune {
			case termi.RuneEscape:
				return
			}
		}
		draw()
	}
}

func main() {
	start()
	defer finish()
	mainLoop()
}
