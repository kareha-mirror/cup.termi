package termi

import (
	"unicode"
)

var TabWidth = 4

func isControl(r rune) bool {
	return r < 32 || (r >= 0x7f && r < 0xa0)
}

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

func RuneWidth(r rune) int {
	// control code
	if isControl(r) {
		return 2
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

func RuneWidthWithTab(r rune, x int) int {
	// tab
	if r == '\t' {
		return TabWidth - (x % TabWidth)
	}

	// others
	return RuneWidth(r)
}

func StringWidth(s string, col int) int {
	sum := 0
	i := 0
	for _, r := range s {
		if i >= col {
			break
		}
		w := RuneWidthWithTab(r, sum)
		sum += w
		i++
	}
	return sum
}

func Render(s string) string {
	buf := []rune{}
	x := 0
	for _, r := range s {
		if r == '\t' {
			spaces := TabWidth - (x % TabWidth)
			for i := 0; i < spaces; i++ {
				buf = append(buf, ' ')
			}
			x += spaces
		} else if r < 0x20 {
			buf = append(buf, '^')
			buf = append(buf, r+'@')
		} else if r == 0x7f {
			buf = append(buf, '^')
			buf = append(buf, '?')
		} else if r >= 0x80 && r < 0xa0 {
			buf = append(buf, '^')
			buf = append(buf, '=')
		} else {
			buf = append(buf, r)
			x += RuneWidthWithTab(r, x)
		}
	}
	return string(buf)
}

func Wrap(s string, w int, tail bool) []string {
	if s == "" {
		return []string{""}
	}
	lines := []string{}
	runes := []rune{}
	sum := 0
	for _, r := range s {
		rw := RuneWidthWithTab(r, sum)
		sum += rw
		if sum > w {
			lines = append(lines, string(runes))
			runes = runes[:0]
			sum = rw
		}
		runes = append(runes, r)
		if sum >= w || (r == '\t' && sum+TabWidth > w) {
			lines = append(lines, string(runes))
			runes = runes[:0]
			sum = 0
		}
	}
	if len(runes) > 0 {
		lines = append(lines, string(runes))
	} else if tail && sum < 1 {
		lines = append(lines, "")
	}
	return lines
}
