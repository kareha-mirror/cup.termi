package termi

import (
	"strconv"
	"strings"
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

func InitKey() error {
	err := initRead()
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

func FinishKey() error {
	err := stopRead()
	if err != nil {
		return err
	}
	waitKey()
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

	params := strings.Builder{}
	for {
		b, ok = readByte()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		key = append(key, b)
		if b < 0x30 || b > 0x3f {
			break
		}
		params.WriteRune(rune(b))
	}

	//inter := strings.Builder{}
	for {
		if b < 0x20 || b > 0x2f {
			break
		}
		//inter.WriteRune(rune(b))
		b, ok = readByte()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		key = append(key, b)
	}

	if b < 0x40 || b > 0x7e { // final
		keyBuf = append(keyBuf, key[1:]...)
		return Key{KeyRune, rune(key[0]), ""}
	}

	switch b {
	case 'A':
		return Key{KeyUp, 0, string(key)}
	case 'B':
		return Key{KeyDown, 0, string(key)}
	case 'C':
		return Key{KeyRight, 0, string(key)}
	case 'D':
		return Key{KeyLeft, 0, string(key)}
	case '~':
		p := params.String()
		switch p {
		case "200":
			return Key{KeyBeginPaste, 0, string(key)}
		case "201":
			return Key{KeyEndPaste, 0, string(key)}
		default:
			return Key{KeyUnknown, 0, string(key)}
		}
	case '_':
		parts := strings.Split(params.String(), ";")
		if len(parts) != 6 {
			return Key{KeyUnknown, 0, string(key)}
		}
		uc, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return Key{KeyUnknown, 0, string(key)}
		}
		if uc == 0 {
			return Key{KeyUnknown, 0, string(key)}
		}
		kd, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return Key{KeyUnknown, 0, string(key)}
		}
		if kd != 1 {
			return Key{KeyUnknown, 0, ""}
		}
		return Key{KeyRune, rune(uc), ""}
	}

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
