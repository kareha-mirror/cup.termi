package termi

import (
	"fmt"
	"strings"
)

type ColorPair struct {
	fg, bg Color
}

func NewColorPair(fg, bg Color) ColorPair {
	return ColorPair{fg, bg}
}

func (cp ColorPair) Pair() (Color, Color) {
	return cp.fg, cp.bg
}

func splitColorPairString(s string) (string, string) {
	parts := strings.Split(s, ",")
	if len(parts[0]) < 1 {
		parts[0] = "default"
	}
	if len(parts) < 2 {
		parts = append(parts, "default")
	}
	return parts[0], parts[1]
}

func ParseColorPair(s string) (ColorPair, error) {
	fgStr, bgStr := splitColorPairString(s)
	fg, err := ParseColor(fgStr)
	if err != nil {
		return ColorPair{}, err
	}
	bg, err := ParseColor(bgStr)
	if err != nil {
		return ColorPair{}, err
	}
	return ColorPair{fg, bg}, nil
}

func (cp ColorPair) Seq() string {
	var fg, bg string
	if cp.fg.def {
		fg = "39"
	} else {
		switch colorMode {
		case ColorMode16:
			fg = cp.fg.cast().fg16()
		case ColorMode256:
			fg = cp.fg.cast().fg256()
		case ColorModeTrue:
			fg = cp.fg.cast().fgTrue()
		default: // unknown color mode
			fg = "39"
		}
	}
	if cp.bg.def {
		bg = "49"
	} else {
		switch colorMode {
		case ColorMode16:
			bg = cp.bg.cast().bg16()
		case ColorMode256:
			bg = cp.bg.cast().bg256()
		case ColorModeTrue:
			bg = cp.bg.cast().bgTrue()
		default: // unknown color mode
			bg = "49"
		}
	}
	return fmt.Sprintf("\x1b[%s;%sm", fg, bg)
}
