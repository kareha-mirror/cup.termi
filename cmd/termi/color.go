package main

import (
	"fmt"

	"tea.kareha.org/cup/termi"
)

var colorLabel = true
var colorBg = true
var colorComplement = true

func drawColor() {
	fmt.Print(termi.HomeCursor)

	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			idx := 16*j + i
			r, g, b := termi.Palette(idx).RGB()
			complement := int(r)+int(g)+int(b) < 192
			if colorBg {
				if colorLabel {
					if colorComplement && complement {
						fmt.Print(termi.Palette(15).Fg())
					} else {
						fmt.Print(termi.Palette(0).Fg())
					}
					fmt.Printf("%s %3d", termi.Palette(idx).Bg(), idx)
				} else {
					fmt.Printf("%s    ", termi.Palette(idx).Bg())
				}
			} else {
				if colorComplement && complement {
					fmt.Print(termi.Palette(15).Bg())
				} else {
					fmt.Print(termi.Palette(0).Bg())
				}
				fmt.Printf("%s %3d", termi.Palette(idx).Fg(), idx)
			}
		}
		fmt.Print("\r\n")
	}

	fmt.Print(termi.ResetAttr)

	fmt.Print("\r\n")
	fmt.Print("q: quit\r\n")
	fmt.Print("t: toggle label\r\n")
	fmt.Print("b: toggle background\r\n")
	fmt.Print("c: toggle complement\r\n")
}

func colorMain() {
	for {
		drawColor()

		key := <-termi.Keys()
		switch key.Kind {
		case termi.KeyRune:
			switch key.Rune {
			case 'q':
				return
			case 't':
				colorLabel = !colorLabel
			case 'b':
				colorBg = !colorBg
			case 'c':
				colorComplement = !colorComplement
			}
		}
	}
}
