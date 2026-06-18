package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

func drawColor() {
	fmt.Print(termi.HomeCursor)

	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			idx := 16*j + i
			fmt.Print(termi.Palette(idx).Fg())
			fmt.Printf(" %3d", idx)
		}
		fmt.Print("\r\n")
	}
	fmt.Print(termi.ResetAttr)
}

func colorMain() {
	drawColor()

	<-termi.Keys()
}
