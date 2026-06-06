package main

import (
	"fmt"
	"os"

	"tea.kareha.org/cup/termi"
)

func usage() {
	fmt.Printf("Usage: %s COMMAND\n", os.Args[0])
	fmt.Print("COMMAND:\n")
	fmt.Print("  color: show 256 color table\n")
	fmt.Print("  seq: key / sequence tester\n")
	fmt.Print("  size: detect screen size and show corners\n")
}

func start() {
	termi.Raw()
	fmt.Print(termi.Clear)
	fmt.Print(termi.HideCursor)
	fmt.Print(termi.HomeCursor)

	termi.Init()
}

func finish() {
	termi.Finish()

	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
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
	case "seq":
		seqMain()
	case "size":
		sizeMain()
	}
}
