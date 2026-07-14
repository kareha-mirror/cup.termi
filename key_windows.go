//go:build windows

package termi

import (
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

var (
	readChan chan Key
	readDone chan struct{}
	keyBuf   []Key
)

func internalInitKey() {
	readChan = make(chan Key, 32)
	readDone = make(chan struct{})
	keyBuf = nil
	keySurrogate = surrogateState{}
}

func readWinKey() (Key, bool) {
	if len(keyBuf) > 0 {
		k := keyBuf[0]
		keyBuf = keyBuf[1:]
		return k, true
	}
	select {
	case k := <-readChan:
		return k, true
	case <-readDone:
		return Key{}, false
	}
}

func readWinKeyTimeout(d time.Duration) (Key, bool) {
	if len(keyBuf) > 0 {
		k := keyBuf[0]
		keyBuf = keyBuf[1:]
		return k, true
	}
	select {
	case k := <-readChan:
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
	k, ok := readWinKey()
	if !ok {
		return Key{keyQuit, 0, ""}
	}
	if k.Kind != KeyRune {
		return k
	}
	if k.Rune != 0x1b { // Escape
		return k
	}

	keys := []Key{k}

	if EscapeTimeout <= 0 {
		k, ok = readWinKey()
		if !ok {
			return Key{keyQuit, 0, ""}
		}
	} else {
		var ok bool
		k, ok = readWinKeyTimeout(EscapeTimeout)
		if !ok {
			return keys[0]
		}
	}

	keys = append(keys, k)
	if k.Kind != KeyRune {
		keyBuf = append(keyBuf, keys[1:]...)
		return keys[0]
	}
	if k.Rune != '[' {
		keyBuf = append(keyBuf, keys[1:]...)
		return keys[0]
	}

	params := strings.Builder{}
	for {
		k, ok = readWinKey()
		if !ok {
			return Key{keyQuit, 0, ""}
		}
		if k.Kind != KeyRune {
			break
		}
		keys = append(keys, k)
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
		k, ok = readWinKey()
		if !ok {
			return Key{keyQuit, 0, ""}
		}
		if k.Kind != KeyRune {
			break
		}
		keys = append(keys, k)
	}

	if k.Kind != KeyRune {
		keyBuf = append(keyBuf, keys[1:]...)
		return keys[0]
	}
	if k.Rune < 0x40 || k.Rune > 0x7e { // final
		keyBuf = append(keyBuf, keys[1:]...)
		return keys[0]
	}

	switch k.Rune {
	case '_':
		parts := strings.Split(params.String(), ";")
		if len(parts) != 6 {
			return Key{KeyUnknown, 0, keysToString(keys)}
		}
		uc, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			return Key{KeyUnknown, 0, keysToString(keys)}
		}
		if uc == 0 {
			return Key{KeyUnknown, 0, keysToString(keys)}
		}
		kd, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return Key{KeyUnknown, 0, keysToString(keys)}
		}
		if uc == 0x1b {
			if kd == 1 {
				return Key{KeyEscapeDown, 0, keysToString(keys)}
			} else {
				return Key{KeyEscapeUp, 0, keysToString(keys)}
			}
		}
		if kd != 1 {
			return Key{KeyUnknown, 0, ""}
		}
		c := uint16(uc)
		if isHighSurrogate(c) {
			keySurrogate.pending = c
			keySurrogate.hasHigh = true
			return Key{KeyUnknown, 0, ""}
		}
		if keySurrogate.hasHigh && isLowSurrogate(c) {
			r := utf16.DecodeRune(rune(keySurrogate.pending), rune(c))
			return Key{KeyRune, r, ""}
		} else {
			return Key{KeyRune, rune(c), ""}
		}
	}

	return Key{KeyUnknown, 0, keysToString(keys)}
}
