package rutil

import (
	"testing"
)

func TestByteIndex(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		col       int
		wantIndex int
	}{
		{"empty", "", 0, 0},
		{"empty negative", "", -2, 0},
		{"empty over", "", 2, 0},
		{"found", "Hello, World!", 9, 9},
		{"found multibyte", "こんにちは、世界！", 4, 12},
		{"negative", "Hello, World!", -2, 0},
		{"over", "Hello, World!", 14, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIndex := ByteIndex(tt.s, tt.col)
			if gotIndex != tt.wantIndex {
				t.Errorf(
					"ByteIndex(\"%s\", %d) = %d; wanted %d",
					tt.s, tt.col, gotIndex, tt.wantIndex,
				)
			}
		})
	}
}

func TestHead(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		end      int
		wantHead string
	}{
		{"empty", "", 0, ""},
		{"empty negative", "", -2, ""},
		{"empty over", "", 2, ""},
		{"found", "Hello, World!", 9, "Hello, Wo"},
		{"found multibyte", "こんにちは、世界！", 4, "こんにち"},
		{"negative", "Hello, World!", -2, ""},
		{"over", "Hello, World!", 14, "Hello, World!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHead := Head(tt.s, tt.end)
			if gotHead != tt.wantHead {
				t.Errorf(
					"Head(\"%s\", %d) = \"%s\"; wanted \"%s\"",
					tt.s, tt.end, gotHead, tt.wantHead,
				)
			}
		})
	}
}

func TestTail(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		start    int
		wantTail string
	}{
		{"empty", "", 0, ""},
		{"empty negative", "", -2, ""},
		{"empty over", "", 2, ""},
		{"found", "Hello, World!", 9, "rld!"},
		{"found multibyte", "こんにちは、世界！", 4, "は、世界！"},
		{"negative", "Hello, World!", -2, "Hello, World!"},
		{"over", "Hello, World!", 14, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTail := Tail(tt.s, tt.start)
			if gotTail != tt.wantTail {
				t.Errorf(
					"Tail(\"%s\", %d) = \"%s\"; wanted \"%s\"",
					tt.s, tt.start, gotTail, tt.wantTail,
				)
			}
		})
	}
}

func TestBody(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		start    int
		end      int
		wantBody string
	}{
		{"empty", "", 0, 0, ""},
		{"found", "Hello, World!", 3, 9, "lo, Wo"},
		{"found multibyte", "こんにちは、世界！", 2, 5, "にちは"},
		{"negative", "Hello, World!", -4, -2, ""},
		{"over", "Hello, World!", 14, 16, ""},

		{"negative start", "Hello", -2, 2, "He"},
		{"reversed", "Hello", 4, 2, ""},
		{"same", "Hello", 2, 2, ""},
		{"end over", "Hello", 2, 20, "llo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody := Body(tt.s, tt.start, tt.end)
			if gotBody != tt.wantBody {
				t.Errorf(
					"Body(\"%s\", %d, %d) = \"%s\"; wanted \"%s\"",
					tt.s, tt.start, tt.end, gotBody, tt.wantBody,
				)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		col      int
		wantHead string
		wantTail string
	}{
		{"empty", "", 0, "", ""},
		{"empty negative", "", -2, "", ""},
		{"empty over", "", 2, "", ""},
		{"found", "Hello, World!", 9, "Hello, Wo", "rld!"},
		{"found multibyte", "こんにちは、世界！", 4, "こんにち", "は、世界！"},
		{"negative", "Hello, World!", -2, "", "Hello, World!"},
		{"over", "Hello, World!", 14, "Hello, World!", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHead, gotTail := Split(tt.s, tt.col)
			if gotHead != tt.wantHead || gotTail != tt.wantTail {
				t.Errorf(
					"Split(\"%s\", %d) = \"%s\", \"%s\"; wanted \"%s\", \"%s\"",
					tt.s, tt.col, gotHead, gotTail, tt.wantHead, tt.wantTail,
				)
			}
		})
	}
}

func TestSplitBody(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		start    int
		end      int
		wantHead string
		wantBody string
		wantTail string
	}{
		{"empty", "", 0, 0, "", "", ""},
		{"found", "Hello, World!", 3, 9, "Hel", "lo, Wo", "rld!"},
		{
			"found multibyte",
			"こんにちは、世界！",
			2, 5,
			"こん", "にちは", "、世界！",
		},
		{"negative", "Hello, World!", -4, -2, "", "", "Hello, World!"},
		{"over", "Hello, World!", 14, 16, "Hello, World!", "", ""},

		{"negative start", "Hello", -2, 2, "", "He", "llo"},
		{"reversed", "Hello", 4, 2, "Hell", "", "o"},
		{"same", "Hello", 2, 2, "He", "", "llo"},
		{"end over", "Hello", 2, 20, "He", "llo", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHead, gotBody, gotTail := SplitBody(tt.s, tt.start, tt.end)
			if gotHead != tt.wantHead ||
				gotBody != tt.wantBody ||
				gotTail != tt.wantTail {
				t.Errorf(
					"SplitBody(\"%s\", %d, %d) = \"%s\", \"%s\", \"%s\"; wanted \"%s\", \"%s\", \"%s\"",
					tt.s, tt.start, tt.end, gotHead, gotBody, gotTail, tt.wantHead, tt.wantBody, tt.wantTail,
				)
			}
		})
	}
}
