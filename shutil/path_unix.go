//go:build unix

package shutil

import (
	"os"
)

const fallbackPath = "/bin/sh"

func Path() string {
	if path := os.Getenv("SHELL"); path != "" {
		return path
	}
	return fallbackPath
}
