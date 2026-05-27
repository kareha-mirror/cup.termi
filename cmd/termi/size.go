package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func drawSize() {
	w, h := termi.Size()

	fmt.Print(termi.MoveCursor(0, 0))
	fmt.Print("+")
	fmt.Print(termi.MoveCursor(w-1, 0))
	fmt.Print("+")
	fmt.Print(termi.MoveCursor(0, h-1))
	fmt.Print("+")
	fmt.Print(termi.MoveCursor(w-1, h-1))
	fmt.Print("+")
}

func sizeMain() {
	drawSize()

	termi.ReadSeq()
}
