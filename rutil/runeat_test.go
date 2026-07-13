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
		{"empty", "", 0, utf8.RuneError},
		{"found", "Hello, World!", 9, 'r'},
		{"found multibyte", "こんにちは、世界！", 4, 'は'},
		{"negative", "Hello, World!", -2, utf8.RuneError},
		{"over", "Hello, World!", 14, utf8.RuneError},
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
