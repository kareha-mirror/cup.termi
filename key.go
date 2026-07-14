package termi

import (
	"sync"
	"time"
)

var EscapeTimeout = 100 * time.Millisecond

type KeyKind int

const (
	keyQuit KeyKind = iota // private

	KeyRune

	KeyUp
	KeyDown
	KeyRight
	KeyLeft

	KeyBeginPaste
	KeyEndPaste

	KeyUnknown

	// w32-input-mode
	KeyEscapeDown
	KeyEscapeUp
)

type Key struct {
	Kind KeyKind
	Rune rune
	Raw  string
}

const (
	RuneEscape    rune = 0x1b
	RuneEnter     rune = '\r'
	RuneNewline   rune = '\n'
	RuneBackspace rune = '\b'
	RuneDelete    rune = 0x7f
)

var (
	keyWG   sync.WaitGroup
	keyChan chan Key
)

func InitKey() error {
	err := initRead()
	if err != nil {
		return err
	}

	internalInitKey()
	keyChan = make(chan Key, 32)

	spawnReader()

	keyWG.Add(1)
	go func() {
		defer keyWG.Done()
		for {
			key := readKey()
			if key.Kind == keyQuit {
				break
			}
			keyChan <- key
		}
		close(keyChan)
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
	return keyChan
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
