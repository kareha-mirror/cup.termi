package termi

import (
	"sync"
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
	SeqQuit
)

type Seq struct {
	Kind SeqKind
	Rune rune
	Raw  string
}

const RuneEscape rune = 0x1b
const RuneEnter rune = '\r'
const RuneNewline rune = '\n'
const RuneBackspace rune = '\b'
const RuneDelete rune = 0x7f

var inputWG sync.WaitGroup
var inputCh chan byte
var inputDone chan struct{}
var inputBuf []byte

func readByte() (byte, bool) {
	if len(inputBuf) > 0 {
		b := inputBuf[0]
		inputBuf = inputBuf[1:]
		return b, true
	}
	select {
	case b := <-inputCh:
		return b, true
	case <-inputDone:
		return 0, false
	}
}

func readByteTimeout(d time.Duration) (byte, bool) {
	if len(inputBuf) > 0 {
		b := inputBuf[0]
		inputBuf = inputBuf[1:]
		return b, true
	}

	select {
	case b := <-inputCh:
		return b, true
	case <-inputDone:
		return 0, false
	case <-time.After(d):
		return 0, false
	}
}

func StartInput() error {
	err := startRead()
	if err != nil {
		return err
	}

	inputCh = make(chan byte, 32)
	inputDone = make(chan struct{})
	inputBuf = make([]byte, 0)

	inputWG.Add(1)
	go func() {
		defer inputWG.Done()
		for {
			b, err := read()
			if err != nil {
				break
			}
			inputCh <- b
		}
		close(inputDone)
		close(inputCh)
	}()

	return nil
}

func StopInput() error {
	err := stopRead()
	if err != nil {
		return err
	}
	inputWG.Wait()
	return finishRead()
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
	b, ok := readByte()
	if !ok {
		return Seq{SeqQuit, 0, ""}
	}
	if b != 0x1b { // Escape
		expected := runeSize(b)
		if expected == -1 {
			panic("Invalid UTF-8 head")
		}
		full := make([]byte, expected)
		full[0] = b
		if expected > 1 {
			for i := 1; i < len(full); i++ {
				b, ok = readByte()
				if !ok {
					return Seq{SeqQuit, 0, ""}
				}
				full[i] = b
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
			b, ok = readByte()
			if !ok {
				return
			}
			seq = append(seq, b)
			if b >= 0x40 && b <= 0x7e {
				return
			}
		}
	}

	if EscapeTimeout <= 0 {
		b, ok = readByte()
		if !ok {
			return Seq{SeqQuit, 0, ""}
		}
	} else {
		var ok bool
		b, ok = readByteTimeout(EscapeTimeout)
		if !ok {
			return Seq{SeqRune, rune(seq[0]), string(seq)}
		}
	}

	seq = append(seq, b)
	if b != '[' {
		inputBuf = append(inputBuf, seq[1:]...)
		return Seq{SeqRune, rune(seq[0]), ""}
	}

	b, ok = readByte()
	if !ok {
		return Seq{SeqQuit, 0, ""}
	}
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
		b, ok = readByte()
		if !ok {
			return Seq{SeqQuit, 0, ""}
		}
		seq = append(seq, b)
		if b != '0' {
			skip(b)
			return Seq{SeqUnknown, 0, string(seq)}
		}

		b, ok = readByte()
		if !ok {
			return Seq{SeqQuit, 0, ""}
		}
		seq = append(seq, b)
		if b != '0' && b != '1' {
			skip(b)
			return Seq{SeqUnknown, 0, string(seq)}
		}

		b2, ok := readByte()
		if !ok {
			return Seq{SeqQuit, 0, ""}
		}
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

//
// Escape Listener
//

type EscapeListener *func(bool)

var escapeListener EscapeListener

func SetEscapeListener(f EscapeListener) {
	escapeListener = f
}

var prevEscape = false

func checkEscape(b byte) {
	escape := b == 0x1b
	if escape == prevEscape {
		return
	}
	if escapeListener != nil {
		(*escapeListener)(escape)
	}
	prevEscape = escape
}
