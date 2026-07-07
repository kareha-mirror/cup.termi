package rutil

import (
	"testing"
	"unicode/utf8"
)

func TestRuneAt(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		col      int
		wantRune rune
	}{
		{"null string", "", 0, utf8.RuneError},
		{"hello world", "Hello, World!", 9, 'r'},
		{"japanese hello", "こんにちは、世界！", 4, 'は'},
		{"invalid hello", "Hello, World!", -2, utf8.RuneError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRune := RuneAt(tt.s, tt.col)
			if gotRune != tt.wantRune {
				t.Errorf(
					"RuneAt(\"%s\", %d) = '%c'; wanted '%c'",
					tt.s, tt.col, gotRune, tt.wantRune,
				)
			}
		})
	}
}
