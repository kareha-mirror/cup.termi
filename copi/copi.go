package copi

import (
	"sync"

	"golang.design/x/clipboard"
)

var (
	once    sync.Once
	initErr error
)

func ensureInitialized() error {
	once.Do(func() {
		initErr = clipboard.Init()
	})
	return initErr
}

func Read() (string, error) {
	if err := ensureInitialized(); err != nil {
		return "", err
	}
	b := clipboard.Read(clipboard.FmtText)
	return string(b), nil
}

func Write(s string) error {
	if err := ensureInitialized(); err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, []byte(s))
	return nil
}
