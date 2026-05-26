package main

import (
	"fmt"
	"os"

	"tea.kareha.org/cup/termi"
)

func usage() {
	fmt.Printf("Usage: %s COMMAND\n", os.Args[0])
	fmt.Print("COMMAND: color key size\n")
}

func start() {
	termi.Raw()
}

func finish() {
	fmt.Print(termi.Clear())
	fmt.Print(termi.HomeCursor())
	termi.Cooked()
	fmt.Print(termi.ShowCursor())
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	start()
	defer finish()

	switch os.Args[1] {
	case "color":
		colorMain()
	case "key":
		keyMain()
	case "size":
		sizeMain()
	}
}
