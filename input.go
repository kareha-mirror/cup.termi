package termi

import (
	"io"
	"os"
	"unicode/utf8"
)

type KeyKind int

const (
	KeyRune KeyKind = iota
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
)

type Key struct {
	Kind KeyKind
	Rune rune
	Raw  string
}

const RuneEscape rune = 0x1b
const RuneEnter rune = '\r'
const RuneBackspace rune = '\b'
const RuneDelete rune = 0x7f

var buf []byte = make([]byte, 0)

func readByte() byte {
	if len(buf) > 0 {
		b := buf[0]
		buf = buf[1:]
		return b
	}

	b := make([]byte, 1)
	_, err := io.ReadFull(os.Stdin, b)
	if err != nil {
		panic(err)
	}
	return b[0]
}

func runeSize(b byte) int {
	switch {
	case b&0x80 == 0:
		return 1
	case b&0xe0 == 0xc0:
		return 2
	case b&0xf0 == 0xe0:
		return 3
	case b&0xf8 == 0xf0:
		return 4
	default:
		return -1 // invalid
	}
}

func ReadKey() Key {
	b := readByte()
	if b != 0x1b { // Escape
		expected := runeSize(b)
		if expected == -1 {
			panic("Invalid UTF-8 head")
		}
		full := make([]byte, expected)
		full[0] = b
		if expected > 1 {
			_, err := io.ReadFull(os.Stdin, full[1:])
			if err != nil {
				panic(err)
			}
		}
		r, size := utf8.DecodeRune(full)
		if r == utf8.RuneError && size == 1 {
			panic("Invalid UTF-8 body")
		}
		return Key{KeyRune, r, ""}
	}

	seq := []byte{b}

	b = readByte()
	seq = append(seq, b)
	if b != '[' {
		buf = append(buf, seq[1:]...)
		return Key{KeyRune, rune(seq[0]), ""}
	}

	b = readByte()
	seq = append(seq, b)
	switch b {
	case 'A':
		return Key{KeyUp, 0, string(seq)}
	case 'B':
		return Key{KeyDown, 0, string(seq)}
	case 'C':
		return Key{KeyRight, 0, string(seq)}
	case 'D':
		return Key{KeyLeft, 0, string(seq)}
	}

	buf = append(buf, seq[1:]...)
	return Key{KeyRune, rune(seq[0]), ""}
}
