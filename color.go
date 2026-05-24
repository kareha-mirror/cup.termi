package termi

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ResetColor() {
	fmt.Print("\x1b[0m")
}

func DefaultColor() {
	fmt.Print("\x1b[39m") // default fg
	fmt.Print("\x1b[49m") // default bg
}

type ColorMode uint8

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
	return ColorMode16
}

var colorMode = GetColorMode()

func SetColorMode(mode ColorMode) {
	colorMode = mode
}

type Color struct {
	indexed bool
	index   uint8
	r, g, b uint8
}

func (c Color) IsIndex() bool {
	return c.indexed
}

func (c Color) Index() uint8 {
	return c.index
}

func (c Color) R() uint8 {
	return c.r
}

func (c Color) G() uint8 {
	return c.g
}

func (c Color) B() uint8 {
	return c.b
}

// normal
var Black = Color{true, 0, 0, 0, 0}
var Red = Color{true, 1, 128, 0, 0}
var Green = Color{true, 2, 0, 128, 0}
var Yellow = Color{true, 3, 128, 128, 0}
var Blue = Color{true, 4, 0, 0, 128}
var Magenta = Color{true, 5, 128, 0, 128}
var Cyan = Color{true, 6, 0, 128, 128}
var White = Color{true, 7, 192, 192, 192}

// bright
var BrightBlack = Color{true, 8, 128, 128, 128}
var BrightRed = Color{true, 9, 255, 0, 0}
var BrightGreen = Color{true, 10, 0, 255, 0}
var BrightYellow = Color{true, 11, 255, 255, 0}
var BrightBlue = Color{true, 12, 0, 0, 255}
var BrightMagenta = Color{true, 13, 255, 0, 255}
var BrightCyan = Color{true, 14, 0, 255, 255}
var BrightWhite = Color{true, 15, 255, 255, 255}

var Palette = make([]Color, 256)

func init() {
	// normal
	Palette[0] = Black
	Palette[1] = Red
	Palette[2] = Green
	Palette[3] = Yellow
	Palette[4] = Blue
	Palette[5] = Magenta
	Palette[6] = Cyan
	Palette[7] = White

	// bright
	Palette[8] = BrightBlack
	Palette[9] = BrightRed
	Palette[10] = BrightGreen
	Palette[11] = BrightYellow
	Palette[12] = BrightBlue
	Palette[13] = BrightMagenta
	Palette[14] = BrightCyan
	Palette[15] = BrightWhite

	level := []uint8{0, 95, 135, 175, 215, 255}

	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				i := r*36 + g*6 + b
				Palette[16+i] = Color{
					true,
					uint8(16 + i),
					level[r],
					level[g],
					level[b],
				}
			}
		}
	}

	for k := 0; k < 24; k++ {
		v := 8 + k*10
		Palette[232+k] = Color{
			true,
			uint8(232 + k),
			uint8(v),
			uint8(v),
			uint8(v),
		}
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

	return Color{
		false,
		0,
		uint8(r),
		uint8(g),
		uint8(b),
	}, nil
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
	// try pase as RGB hex
	if len(s) == 6 {
		c, err := parseHexColor(s)
		if err == nil {
			return c, nil
		}
	}

	// try parse as palette index
	i, err := strconv.ParseUint(s, 10, 8)
	if err == nil {
		return Palette[i], nil
	}

	// try parse as color name
	c, err := parseNamedColor(s)
	if err == nil {
		return c, nil
	}

	return Color{}, fmt.Errorf("invalid color format")
}

func getColor16Index(c Color) uint8 {
	r := int(c.r)
	g := int(c.g)
	b := int(c.b)

	mx := max(r, g, b)
	mn := min(r, g, b)
	brightness := (r + g + b) / 3

	if mx-mn < 32 { // gray
		if brightness < 64 {
			return 0 // black
		}
		if brightness >= 192 {
			return 7 + 8 // white
		}
		return 7 // white
	}

	var bright uint8 = 0
	if brightness > 128 {
		bright = 8
	}

	if r >= g && r >= b {
		if g >= 128 {
			return 3 + bright // yellow
		}
		if b >= 128 {
			return 5 + bright // magenta
		}
		return 1 + bright // red
	}
	if g >= r && g >= b {
		if r >= 128 {
			return 3 + bright // yellow
		}
		if b >= 128 {
			return 6 + bright // cyan
		}
		return 2 + bright // green
	}
	// b >= r && b >= g
	if r >= 128 {
		return 5 + bright // magenta
	}
	if g >= 128 {
		return 6 + bright // cyan
	}
	return 4 + bright // blue
}

func getColor256Index(c Color) uint8 {
	r := (int(c.r) + 21) * 5 / 255
	g := (int(c.g) + 21) * 5 / 255
	b := (int(c.b) + 21) * 5 / 255
	return uint8(16 + 36*r + 6*g + b)
}

func CastColor(c Color) Color {
	switch colorMode {
	case ColorMode16:
		if c.indexed && c.index < 16 {
			return c
		}
		i := getColor16Index(c)
		return Palette[i]
	case ColorMode256:
		if c.indexed {
			return c
		}
		i := getColor256Index(c)
		return Palette[i]
	default: // ColorModeTrue
		return c
	}
}

func SetFgColor(c Color) {
	c = CastColor(c)
	switch colorMode {
	case ColorMode16:
		setFgColor16(c)
	case ColorMode256:
		setFgColor256(c)
	case ColorModeTrue:
		setFgColorTrue(c)
	}
}

func SetBgColor(c Color) {
	c = CastColor(c)
	switch colorMode {
	case ColorMode16:
		setBgColor16(c)
	case ColorMode256:
		setBgColor256(c)
	case ColorModeTrue:
		setBgColorTrue(c)
	}
}

func setFgColor16(c Color) {
	var code uint8
	if c.index >= 8 {
		code = 90 + c.index - 8
	} else {
		code = 30 + c.index
	}
	fmt.Printf("\x1b[%dm", code)
}

func setBgColor16(c Color) {
	var code uint8
	if c.index >= 8 {
		code = 100 + c.index - 8
	} else {
		code = 40 + c.index
	}
	fmt.Printf("\x1b[%dm", code)
}

func setFgColor256(c Color) {
	fmt.Printf("\x1b[38;5;%dm", c.index)
}

func setBgColor256(c Color) {
	fmt.Printf("\x1b[48;5;%dm", c.index)
}

func setFgColorTrue(c Color) {
	fmt.Printf("\x1b[38;2;%d;%d;%dm", c.r, c.g, c.b)
}

func setBgColorTrue(c Color) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm", c.r, c.g, c.b)
}
