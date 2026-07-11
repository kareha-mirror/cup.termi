//go:build windows

package termi

import (
	"strconv"
	"strings"
	"time"
)

var readCh chan Key
var readDone chan struct{}
var keyBuf []Key

func internalInitKey() {
	readCh = make(chan Key, 32)
	readDone = make(chan struct{})
}

func readKeyElem() (Key, bool) {
	if len(keyBuf) > 0 {
		k := keyBuf[0]
		keyBuf = keyBuf[1:]
		return k, true
	}

	select {
	case k := <-readCh:
		return k, true
	case <-readDone:
		return Key{}, false
	}
}

func readKeyElemTimeout(d time.Duration) (Key, bool) {
	if len(keyBuf) > 0 {
		k := keyBuf[0]
		keyBuf = keyBuf[1:]
		return k, true
	}

	select {
	case k := <-readCh:
		return k, true
	case <-readDone:
		return Key{}, false
	case <-time.After(d):
		return Key{}, false
	}
}

func keysToString(keys []Key) string {
	b := strings.Builder{}
	for _, key := range keys {
		if key.Kind != KeyRune {
			continue
		}
		b.WriteRune(key.Rune)
	}
	return b.String()
}

func readKey() Key {
	k, ok := readKeyElem()
	if !ok {
		return Key{KeyQuit, 0, ""}
	}
	if k.Kind != KeyRune {
		return k
	}
	if k.Rune != 0x1b { // Escape
		return k
	}

	key := []Key{k}

	if EscapeTimeout <= 0 {
		k, ok = readKeyElem()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
	} else {
		var ok bool
		k, ok = readKeyElemTimeout(EscapeTimeout)
		if !ok {
			return key[0]
		}
	}

	key = append(key, k)
	if k.Kind != KeyRune {
		keyBuf = append(keyBuf, key[1:]...)
		return key[0]
	}
	if k.Rune != '[' {
		keyBuf = append(keyBuf, key[1:]...)
		return key[0]
	}

	params := strings.Builder{}
	for {
		k, ok = readKeyElem()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		if k.Kind != KeyRune {
			break
		}
		key = append(key, k)
		if k.Rune < 0x30 || k.Rune > 0x3f {
			break
		}
		params.WriteRune(k.Rune)
	}

	//inter := strings.Builder{}
	for {
		if k.Rune < 0x20 || k.Rune > 0x2f {
			break
		}
		//inter.WriteRune(k.Rune)
		k, ok = readKeyElem()
		if !ok {
			return Key{KeyQuit, 0, ""}
		}
		if k.Kind != KeyRune {
			break
		}
		key = append(key, k)
	}

	if k.Kind != KeyRune {
		keyBuf = append(keyBuf, key[1:]...)
		return key[0]
	}
	if k.Rune < 0x40 || k.Rune > 0x7e { // final
		keyBuf = append(keyBuf, key[1:]...)
		return key[0]
	}

	switch k.Rune {
	case '_':
		parts := strings.Split(params.String(), ";")
		if len(parts) != 6 {
			return Key{KeyUnknown, 0, keysToString(key)}
		}
		uc, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return Key{KeyUnknown, 0, keysToString(key)}
		}
		if uc == 0 {
			return Key{KeyUnknown, 0, keysToString(key)}
		}
		kd, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return Key{KeyUnknown, 0, keysToString(key)}
		}
		if uc == 0x1b {
			if kd == 1 {
				return Key{KeyEscapeDown, 0, keysToString(key)}
			} else {
				return Key{KeyEscapeUp, 0, keysToString(key)}
			}
		}
		if kd != 1 {
			return Key{KeyUnknown, 0, ""}
		}
		return Key{KeyRune, rune(uc), ""}
	}

	return Key{KeyUnknown, 0, keysToString(key)}
}

var prevEscape = false

func checkEscape(r rune) {
	escape := r == 0x1b
	if escape == prevEscape {
		return
	}
	if escapeListener != nil {
		(*escapeListener)(escape)
	}
	prevEscape = escape
}
