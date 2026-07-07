package rutil

import (
	"testing"
)

func TestRuneIndex(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		start   int
		r       rune
		wantCol int
	}{
		{"hello world", "Hello, World!", 6, 'o', 8},
		{"valid hello", "Hello, World!", -2, 'o', 4},
		{"invalid hello", "Hello, World!", 14, 'o', -1},
		{"japanese hello", "こんにちは、世界！", 2, 'は', 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCol := RuneIndex(tt.line, tt.start, tt.r)
			if gotCol != tt.wantCol {
				t.Errorf(
					"RuneIndex(\"%s\", %d, '%c') = %d; wanted %d",
					tt.line, tt.start, tt.r, gotCol, tt.wantCol,
				)
			}
		})
	}
}

func TestLastRuneIndex(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		start   int
		r       rune
		wantCol int
	}{
		{"hello world", "Hello, World!", 6, 'o', 4},
		{"invalid hello", "Hello, World!", -2, 'o', -1},
		{"valid hello", "Hello, World!", 14, 'o', 8},
		{"japanese hello", "こんにちは、世界！", 6, 'は', 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCol := LastRuneIndex(tt.line, tt.start, tt.r)
			if gotCol != tt.wantCol {
				t.Errorf(
					"LastRuneIndex(\"%s\", %d, '%c') = %d; wanted %d",
					tt.line, tt.start, tt.r, gotCol, tt.wantCol,
				)
			}
		})
	}
}
