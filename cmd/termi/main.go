package main

import (
	"fmt"
	"os"

	"tea.kareha.org/cup/termi"
	"tea.kareha.org/cup/termi/suspend"
)

func usage() {
	fmt.Printf("Usage: %s COMMAND\n", os.Args[0])
	fmt.Print("COMMAND:\n")
	fmt.Print("  color: show 256 color table\n")
	fmt.Print("  key: key tester\n")
	fmt.Print("  size: detect screen size and show corners\n")
}

func start() {
	termi.Raw()
	fmt.Print(termi.SetAlternate)
	fmt.Print(termi.Clear)
	fmt.Print(termi.HideCursor)
	fmt.Print(termi.HomeCursor)

	termi.InitKey()
	suspend.Init()
}

func finish() {
	suspend.Finish()
	termi.FinishKey()

	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	fmt.Print(termi.ResetAlternate)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	command := os.Args[1]

	start()
	defer finish()

	switch command {
	case "color":
		colorMain()
	case "key":
		keyMain()
	case "size":
		sizeMain()
	}
}
