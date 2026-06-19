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

func (cp ColorPair) Elements() (Color, Color) {
	return cp.fg, cp.bg
}

func ParseColorPair(s string) (ColorPair, error) {
	// split
	parts := strings.Split(s, ",")
	if len(parts[0]) < 1 {
		parts[0] = "default"
	}
	if len(parts) < 2 {
		parts = append(parts, "default")
	}

	// parse
	fg, err := ParseColor(parts[0])
	if err != nil {
		return ColorPair{}, err
	}
	bg, err := ParseColor(parts[1])
	if err != nil {
		return ColorPair{}, err
	}

	return ColorPair{fg, bg}, nil
}

func (cp ColorPair) Seq() string {
	fg := "39"
	if !cp.fg.def {
		switch colorMode {
		case ColorMode16:
			fg = cp.fg.cast16().fg16()
		case ColorMode256:
			fg = cp.fg.cast256().fg256()
		case ColorModeTrue:
			fg = cp.fg.fgTrue()
		}
	}

	bg := "49"
	if !cp.bg.def {
		switch colorMode {
		case ColorMode16:
			bg = cp.bg.cast16().bg16()
		case ColorMode256:
			bg = cp.bg.cast256().bg256()
		case ColorModeTrue:
			bg = cp.bg.bgTrue()
		}
	}

	return fmt.Sprintf("\x1b[%s;%sm", fg, bg)
}
