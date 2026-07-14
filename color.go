package termi

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const ResetAttr = "\x1b[0m" // "\x1b[m"

type ColorMode int

const (
	ColorMode16 ColorMode = iota
	ColorMode256
	ColorModeTrue
)

func GetColorMode() ColorMode {
	if strings.Contains(os.Getenv("COLORTERM"), "truecolor") {
		return ColorModeTrue
	}
	if strings.Contains(os.Getenv("TERM"), "256color") {
		return ColorMode256
	}
	if os.Getenv("WT_SESSION") != "" {
		return ColorModeTrue
	}
	return ColorMode16
}

var colorMode = GetColorMode()

func SetColorMode(mode ColorMode) {
	colorMode = mode
}

type Color struct {
	def     bool
	indexed bool
	index   uint8
	r, g, b uint8
}

var DefaultColor = Color{true, false, 0, 0, 0, 0}

func NewRGB(r, g, b uint8) Color {
	return Color{false, false, 0, r, g, b}
}

func (c Color) Default() bool {
	return c.def
}

func (c Color) Indexed() bool {
	return c.indexed
}

func (c Color) Index() uint8 {
	return c.index
}

func (c Color) RGB() (uint8, uint8, uint8) {
	return c.r, c.g, c.b
}

var (
	// normal
	Black   = Color{false, true, 0, 0, 0, 0}
	Red     = Color{false, true, 1, 128, 0, 0}
	Green   = Color{false, true, 2, 0, 128, 0}
	Yellow  = Color{false, true, 3, 128, 128, 0}
	Blue    = Color{false, true, 4, 0, 0, 128}
	Magenta = Color{false, true, 5, 128, 0, 128}
	Cyan    = Color{false, true, 6, 0, 128, 128}
	White   = Color{false, true, 7, 192, 192, 192}

	// bright
	BrightBlack   = Color{false, true, 8, 128, 128, 128}
	BrightRed     = Color{false, true, 9, 255, 0, 0}
	BrightGreen   = Color{false, true, 10, 0, 255, 0}
	BrightYellow  = Color{false, true, 11, 255, 255, 0}
	BrightBlue    = Color{false, true, 12, 0, 0, 255}
	BrightMagenta = Color{false, true, 13, 255, 0, 255}
	BrightCyan    = Color{false, true, 14, 0, 255, 255}
	BrightWhite   = Color{false, true, 15, 255, 255, 255}
)

var palette = make([]Color, 256)

func Palette(i int) Color {
	if i < 0 || i >= 256 {
		return DefaultColor
	}
	return palette[i]
}

func SetPalette(i int, r, g, b uint8) {
	if i < 0 || i >= 256 {
		return
	}
	palette[i] = Color{false, true, uint8(i), r, g, b}
}

func init() {
	// normal
	palette[0] = Black
	palette[1] = Red
	palette[2] = Green
	palette[3] = Yellow
	palette[4] = Blue
	palette[5] = Magenta
	palette[6] = Cyan
	palette[7] = White

	// bright
	palette[8] = BrightBlack
	palette[9] = BrightRed
	palette[10] = BrightGreen
	palette[11] = BrightYellow
	palette[12] = BrightBlue
	palette[13] = BrightMagenta
	palette[14] = BrightCyan
	palette[15] = BrightWhite

	// cube colors
	level := []uint8{0, 95, 135, 175, 215, 255}
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				i := uint8(16 + r*36 + g*6 + b)
				palette[i] = Color{
					false, true, i, level[r], level[g], level[b],
				}
			}
		}
	}

	// grays
	for k := 0; k < 24; k++ {
		i := uint8(232 + k)
		v := uint8(8 + k*10)
		palette[i] = Color{false, true, i, v, v, v}
	}
}

func parseHexColor(s string) (Color, error) {
	if len(s) != 6 {
		return Color{}, fmt.Errorf("invalid hex color format")
	}

	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return Color{}, err
	}

	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return Color{}, err
	}

	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return Color{}, err
	}

	return Color{false, false, 0, uint8(r), uint8(g), uint8(b)}, nil
}

var normalColors = map[string]Color{
	"black":   Black,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"blue":    Blue,
	"magenta": Magenta,
	"cyan":    Cyan,
	"white":   White,
}

var brightColors = map[string]Color{
	"black":   BrightBlack,
	"red":     BrightRed,
	"green":   BrightGreen,
	"yellow":  BrightYellow,
	"blue":    BrightBlue,
	"magenta": BrightMagenta,
	"cyan":    BrightCyan,
	"white":   BrightWhite,
}

func parseNamedColor(s string) (Color, error) {
	lower := strings.ToLower(s)
	if strings.Contains(lower, "default") {
		return DefaultColor, nil
	}
	var m map[string]Color
	if strings.Contains(lower, "bright") {
		m = brightColors
	} else {
		m = normalColors
	}
	for name, c := range m {
		if strings.Contains(lower, name) {
			return c, nil
		}
	}
	return Color{}, fmt.Errorf("unknown color name")
}

func ParseColor(s string) (Color, error) {
	s = strings.TrimSpace(s)

	// try parse as RGB hex
	if len(s) == 6 {
		c, err := parseHexColor(s)
		if err == nil {
			return c, nil
		}
	}

	// try parse as palette index
	i, err := strconv.ParseUint(s, 10, 8)
	if err == nil {
		return palette[i], nil
	}

	// try parse as color name
	c, err := parseNamedColor(s)
	if err == nil {
		return c, nil
	}

	return Color{}, fmt.Errorf("invalid color format")
}

func (c Color) index16() uint8 {
	// 16 colors
	if c.indexed && c.index < 16 {
		return c.index
	}

	r := int(c.r)
	g := int(c.g)
	b := int(c.b)

	brightness := (299*r + 587*g + 114*b) / 1000

	var bright uint8 = 0
	if brightness >= 160 {
		bright = 8
	}

	mx := max(r, g, b)
	mn := min(r, g, b)

	if mx-mn >= 32 { // cube colors
		if r >= g && r >= b {
			if g+b < 256 {
				return 1 + bright // red
			}
			if g >= b {
				return 3 + bright // yellow
			}
			if b >= g {
				return 5 + bright // magenta
			}
		}
		if g >= b && g >= r {
			if b+r < 256 {
				return 2 + bright // green
			}
			if b >= r {
				return 6 + bright // cyan
			}
			if r >= b {
				return 3 + bright // yellow
			}
		}
		if b >= r && b >= g {
			if r+g < 256 {
				return 4 + bright // blue
			}
			if r >= g {
				return 5 + bright // magenta
			}
			if g >= r {
				return 6 + bright // cyan
			}
		}
	}

	// grays
	if brightness < 64 {
		return 0 // black
	}
	if brightness < 128 {
		return 0 + 8 // bright black
	}
	if brightness < 192 {
		return 7 // white
	}
	return 7 + 8 // bright white
}

func (c Color) index256() uint8 {
	// only use cube colors
	r := (int(c.r) + 21) * 5 / 255
	g := (int(c.g) + 21) * 5 / 255
	b := (int(c.b) + 21) * 5 / 255
	return uint8(16 + 36*r + 6*g + b)
}

func (c Color) cast16() Color {
	if c.indexed && c.index < 16 {
		return c
	}
	return palette[c.index16()]
}

func (c Color) cast256() Color {
	if c.indexed {
		return c
	}
	return palette[c.index256()]
}

func (c Color) fg16() string {
	var code uint8
	if c.index >= 8 {
		code = 90 + c.index - 8
	} else {
		code = 30 + c.index
	}
	return fmt.Sprintf("%d", code)
}

func (c Color) bg16() string {
	var code uint8
	if c.index >= 8 {
		code = 100 + c.index - 8
	} else {
		code = 40 + c.index
	}
	return fmt.Sprintf("%d", code)
}

func (c Color) fg256() string {
	return fmt.Sprintf("38;5;%d", c.index)
}

func (c Color) bg256() string {
	return fmt.Sprintf("48;5;%d", c.index)
}

func (c Color) fgTrue() string {
	return fmt.Sprintf("38;2;%d;%d;%d", c.r, c.g, c.b)
}

func (c Color) bgTrue() string {
	return fmt.Sprintf("48;2;%d;%d;%d", c.r, c.g, c.b)
}

func (c Color) Fg() string {
	fg := "39"
	if !c.def {
		switch colorMode {
		case ColorMode16:
			fg = c.cast16().fg16()
		case ColorMode256:
			fg = c.cast256().fg256()
		case ColorModeTrue:
			fg = c.fgTrue()
		}
	}
	return fmt.Sprintf("\x1b[%sm", fg)
}

func (c Color) Bg() string {
	bg := "49"
	if !c.def {
		switch colorMode {
		case ColorMode16:
			bg = c.cast16().bg16()
		case ColorMode256:
			bg = c.cast256().bg256()
		case ColorModeTrue:
			bg = c.bgTrue()
		}
	}
	return fmt.Sprintf("\x1b[%sm", bg)
}
