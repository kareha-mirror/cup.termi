package rutil

import (
	"unicode/utf8"
)

func RuneAt(s string, col int) rune {
	for _, r := range s {
		if col == 0 {
			return r
		}
		col--
	}
	return utf8.RuneError
}
