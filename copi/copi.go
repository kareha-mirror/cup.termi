package copi

import (
	"golang.design/x/clipboard"
)

var initialized = false

func ensureInitialized() error {
	if initialized {
		return nil
	}
	if err := clipboard.Init(); err != nil {
		return err
	}
	initialized = true
	return nil
}

func Read() (string, error) {
	err := ensureInitialized()
	if err != nil {
		return "", err
	}
	b := clipboard.Read(clipboard.FmtText)
	return string(b), nil
}

func Write(s string) error {
	err := ensureInitialized()
	if err != nil {
		return err
	}
	clipboard.Write(clipboard.FmtText, []byte(s))
	return nil
}
