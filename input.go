package termi

import (
	"io"
	"os"
	"time"
	"unicode/utf8"
)

var EscapeTimeout = 100 * time.Millisecond

type SeqKind int

const (
	SeqRune SeqKind = iota

	SeqUp
	SeqDown
	SeqRight
	SeqLeft

	SeqBeginPaste
	SeqEndPaste

	SeqUnknown
)

type Seq struct {
	Kind SeqKind
	Rune rune
	Raw  string
}

const RuneEscape rune = 0x1b
const RuneEnter rune = '\r'
const RuneBackspace rune = '\b'
const RuneDelete rune = 0x7f

var ch = make(chan byte, 32)
var done = make(chan struct{})
var buf []byte = make([]byte, 0)

func readElem() byte {
	b := make([]byte, 1)
	_, err := io.ReadFull(os.Stdin, b)
	if err != nil {
		panic(err)
	}
	fireEscape(b[0])
	return b[0]
}

func readByte() byte {
	if len(buf) > 0 {
		b := buf[0]
		buf = buf[1:]
		return b
	}
	select {
	case b := <-ch:
		return b
	case <-done:
		return 0
	}
}

func readByteTimeout(d time.Duration) (byte, bool) {
	if len(buf) > 0 {
		b := buf[0]
		buf = buf[1:]
		return b, true
	}

	select {
	case b := <-ch:
		return b, true
	case <-done:
		return 0, false
	case <-time.After(d):
		return 0, false
	}
}

func Init() {
	go func() {
		for {
			select {
			case ch <- readElem():
			case <-done:
				return
			}
		}
	}()
}

func Finish() {
	close(done)
	close(ch)
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

func ReadSeq() Seq {
	b := readByte()
	if b != 0x1b { // Escape
		expected := runeSize(b)
		if expected == -1 {
			panic("Invalid UTF-8 head")
		}
		full := make([]byte, expected)
		full[0] = b
		if expected > 1 {
			for i := 1; i < len(full); i++ {
				full[i] = readByte()
			}
		}
		r, size := utf8.DecodeRune(full)
		if r == utf8.RuneError && size == 1 {
			panic("Invalid UTF-8 body")
		}
		return Seq{SeqRune, r, ""}
	}

	seq := []byte{b}

	skip := func(b byte) {
		if b >= 0x40 && b <= 0x7e {
			return
		}
		for {
			b = readByte()
			seq = append(seq, b)
			if b >= 0x40 && b <= 0x7e {
				return
			}
		}
	}

	if EscapeTimeout <= 0 {
		b = readByte()
	} else {
		var ok bool
		b, ok = readByteTimeout(EscapeTimeout)
		if !ok {
			return Seq{SeqRune, rune(seq[0]), string(seq)}
		}
	}

	seq = append(seq, b)
	if b != '[' {
		buf = append(buf, seq[1:]...)
		return Seq{SeqRune, rune(seq[0]), ""}
	}

	b = readByte()
	seq = append(seq, b)
	switch b {
	case 'A':
		return Seq{SeqUp, 0, string(seq)}
	case 'B':
		return Seq{SeqDown, 0, string(seq)}
	case 'C':
		return Seq{SeqRight, 0, string(seq)}
	case 'D':
		return Seq{SeqLeft, 0, string(seq)}
	}

	if b == '2' {
		b = readByte()
		seq = append(seq, b)
		if b != '0' {
			skip(b)
			return Seq{SeqUnknown, 0, string(seq)}
		}

		b = readByte()
		seq = append(seq, b)
		if b != '0' && b != '1' {
			skip(b)
			return Seq{SeqUnknown, 0, string(seq)}
		}

		b2 := readByte()
		seq = append(seq, b2)
		if b2 != '~' {
			skip(b2)
			return Seq{SeqUnknown, 0, string(seq)}
		}

		if b == '0' {
			return Seq{SeqBeginPaste, 0, string(seq)}
		} else {
			return Seq{SeqEndPaste, 0, string(seq)}
		}
	}

	skip(b)
	return Seq{SeqUnknown, 0, string(seq)}
}

type EscapeListener *func(bool)

var escapeListeners = make([]EscapeListener, 0)

func AddEscapeListener(f EscapeListener) {
	escapeListeners = append(escapeListeners, f)
}

func RemoveEscapeListener(f EscapeListener) bool {
	for i := 0; i < len(escapeListeners); i++ {
		if escapeListeners[i] == f {
			if i+1 < len(escapeListeners) {
				escapeListeners = append(
					escapeListeners[:i], escapeListeners[i+1:]...,
				)
			} else {
				escapeListeners = escapeListeners[:i]
			}
			return true
		}
	}
	return false
}

var prevEsc = false

func fireEscape(b byte) {
	esc := b == 0x1b
	if esc == prevEsc {
		return
	}
	for _, f := range escapeListeners {
		(*f)(esc)
	}
	prevEsc = esc
}
