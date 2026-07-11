package termi

import (
	"sync"
	"time"
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

	KeyEscapeDown
	KeyEscapeUp
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
var keyCh chan Key

func InitKey() error {
	err := initRead()
	if err != nil {
		return err
	}

	internalInitKey()
	keyCh = make(chan Key, 32)

	spawnReader()

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
	keyWG.Wait()
	return finishRead()
}

func Keys() chan Key {
	return keyCh
}

//
// Escape Listener
//

type EscapeListener *func(bool)

var escapeListener EscapeListener

func SetEscapeListener(f EscapeListener) {
	escapeListener = f
}
