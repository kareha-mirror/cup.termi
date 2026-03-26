package termi

import (
	"fmt"
	"unicode"
)

func isWide(r rune) bool {
	return r >= 0x1100 && (r <= 0x115f || // Hangul Jamo
		r == 0x2329 || r == 0x232a ||
		(r >= 0x2e80 && r <= 0xa4cf) ||
		(r >= 0xac00 && r <= 0xd7a3) ||
		(r >= 0xf900 && r <= 0xfaff) ||
		(r >= 0xfe10 && r <= 0xfe19) ||
		(r >= 0xfe30 && r <= 0xfe6f) ||
		(r >= 0xff00 && r <= 0xff60) ||
		(r >= 0xffe0 && r <= 0xffe6))
}

func isEmoji(r rune) bool {
	return r >= 0x1f300 && r <= 0x1faff
}

const tabWidth = 4

func runeWidth(r rune, x int) int {
	// tab
	if r == '\t' {
		return tabWidth - (x % tabWidth)
	}

	// control code
	if r == 0 {
		return 0
	}
	if r < 32 || (r >= 0x7f && r < 0xa0) {
		return 0
	}

	// combining mark
	if unicode.Is(unicode.Mn, r) {
		return 0
	}

	// wide (loose CJK)
	if isWide(r) {
		return 2
	}

	// emoji
	if isEmoji(r) {
		return 2
	}

	return 1
}

func StringWidth(s string, col int) int {
	sum := 0
	i := 0
	for _, r := range s {
		if i >= col {
			break
		}
		w := runeWidth(r, sum)
		sum += w
		i++
	}
	return sum
}

func Print(s string) {
	fmt.Print(s)
}

func Printf(format string, a ...any) (n int, err error) {
	s := fmt.Sprintf(format, a...)
	Print(s)
	return len(s), nil
}

func Draw(s string) {
	x := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabWidth - (x % tabWidth)
			for i := 0; i < spaces; i++ {
				fmt.Print(" ")
			}
			x += spaces
		} else {
			fmt.Printf("%c", r)
			x += runeWidth(r, x)
		}
	}
}
