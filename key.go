package termi

import (
	"sync"
	"time"
	"unicode/utf8"
)

var EscapeTimeout = 100 * time.Millisecond

type KeyKind int

const (
	KeyRune KeyKind = iota

	KeyUp
	KeyDown
	KeyRight
	KeyLeft

	KeyBeginPaste
	KeyEndPaste

	KeyUnknown
	KeyQuit
)

type Key struct {
	Kind KeyKind
	Rune rune
	Raw  string
}

const RuneEscape rune = 0x1b
const RuneEnter rune = '\r'
const RuneNewline rune = '\n'
const RuneBackspace rune = '\b'
const RuneDelete rune = 0x7f

var keyWG sync.WaitGroup
var readCh chan byte
var readDone chan struct{}
var keyBuf []byte
var keyCh chan Key

func readByte() (byte, bool) {
	if len(keyBuf) > 0 {
		b := keyBuf[0]
		keyBuf = keyBuf[1:]
		return b, true
	}
	select {
	case b := <-readCh:
		return b, true
	case <-readDone:
		return 0, false
	}
}

func readByteTimeout(d time.Duration) (byte, bool) {
	if len(keyBuf) > 0 {
		b := keyBuf[0]
		keyBuf = keyBuf[1:]
		return b, true
	}

	select {
	case b := <-readCh:
		return b, true
	case <-readDone:
		return 0, false
	case <-time.After(d):
		return 0, false
	}
}

func StartKey() error {
	err := startRead()
	if err != nil {
		return err
	}

	readCh = make(chan byte, 32)
	readDone = make(chan struct{})
	keyBuf = make([]byte, 0)
	keyCh = make(chan Key, 32)

	keyWG.Add(1)
	go func() {
		defer keyWG.Done()
		for {
			b, err := read()
			if err != nil {
				break
			}
			readCh <- b
		}
		close(readDone)
		close(readCh)
	}()
	keyWG.Add(1)
	go func() {
		defer keyWG.Done()
		for {
			key := readKey()
			if key.Kind == KeyQuit {
				break
			}
			keyCh <- key
		}
		close(keyCh)
	}()

	return nil
}

func StopKey() error {
	err := stopRead()
	if err != nil {
		return err
	}
	keyWG.Wait()
	return finishRead()
}

func Keys() chan Key {
	return keyCh
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

func readKey() Key {
	b, ok := readByte()
	if !ok {
		return Key{KeyQuit, 0, ""}
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
					return Key{KeyQuit, 0, ""}
				}
				full[i] = b
			}
		}
		r, size := utf8.DecodeRune(full)
		if r == utf8.RuneError && size == 1 {
			panic("Invalid UTF-8 body")
		}
		return Key{KeyRune, r, ""}
	}

	key := []byte{b}

	skip := func(b byte) {
		if b >= 0x40 && b <= 0x7e {
			return
		}
		for {
			b, ok = readByte()
			if !ok {
				return
			}
			key = append(key, b)
			if b >= 0x40 && b <= 0x7e {
				return
			}
		}
	}

	if EscapeTimeout <= 0 {
		b, ok = readByte()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
	} else {
		var ok bool
		b, ok = readByteTimeout(EscapeTimeout)
		if !ok {
			return Key{KeyRune, rune(key[0]), string(key)}
		}
	}

	key = append(key, b)
	if b != '[' {
		keyBuf = append(keyBuf, key[1:]...)
		return Key{KeyRune, rune(key[0]), ""}
	}

	b, ok = readByte()
	if !ok {
		return Key{KeyQuit, 0, ""}
	}
	key = append(key, b)
	switch b {
	case 'A':
		return Key{KeyUp, 0, string(key)}
	case 'B':
		return Key{KeyDown, 0, string(key)}
	case 'C':
		return Key{KeyRight, 0, string(key)}
	case 'D':
		return Key{KeyLeft, 0, string(key)}
	}

	if b == '2' {
		b, ok = readByte()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		key = append(key, b)
		if b != '0' {
			skip(b)
			return Key{KeyUnknown, 0, string(key)}
		}

		b, ok = readByte()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		key = append(key, b)
		if b != '0' && b != '1' {
			skip(b)
			return Key{KeyUnknown, 0, string(key)}
		}

		b2, ok := readByte()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		key = append(key, b2)
		if b2 != '~' {
			skip(b2)
			return Key{KeyUnknown, 0, string(key)}
		}

		if b == '0' {
			return Key{KeyBeginPaste, 0, string(key)}
		} else {
			return Key{KeyEndPaste, 0, string(key)}
		}
	}

	skip(b)
	return Key{KeyUnknown, 0, string(key)}
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
