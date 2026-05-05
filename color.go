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
	return ColorMode16
}

var colorMode = GetColorMode()

type Color8 uint8

const (
	Color8Black Color8 = iota
	Color8Red
	Color8Green
	Color8Yellow
	Color8Blue
	Color8Magenta
	Color8Cyan
	Color8White
)

type Color16 struct {
	c      Color8
	bright bool
}

func (c Color16) FgCode() uint8 {
	if c.bright {
		return 90 + uint8(c.c)
	} else {
		return 30 + uint8(c.c)
	}
}

func (c Color16) BgCode() uint8 {
	if c.bright {
		return 100 + uint8(c.c)
	} else {
		return 40 + uint8(c.c)
	}
}

type Color256 uint8

type Color struct {
	r, g, b uint8
}

func SetFgColor(c Color) {
	switch colorMode {
	case ColorMode16:
		SetFgColor16(ToColor16(c))
	case ColorMode256:
		SetFgColor256(ToColor256(c))
	case ColorModeTrue:
		SetFgColorTrue(c)
	}
}

func SetBgColor(c Color) {
	switch colorMode {
	case ColorMode16:
		SetBgColor16(ToColor16(c))
	case ColorMode256:
		SetBgColor256(ToColor256(c))
	case ColorModeTrue:
		SetBgColorTrue(c)
	}
}

func SetFgColor16(c Color16) {
	fmt.Printf("\x1b[%dm", c.FgCode())
}

func SetBgColor16(c Color16) {
	fmt.Printf("\x1b[%dm", c.BgCode())
}

func SetFgColor256(c Color256) {
	fmt.Printf("\x1b[38;5;%dm", c)
}

func SetBgColor256(c Color256) {
	fmt.Printf("\x1b[48;5;%dm", c)
}

func SetFgColorTrue(c Color) {
	fmt.Printf("\x1b[38;2;%d;%d;%dm", c.r, c.g, c.b)
}

func SetBgColorTrue(c Color) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm", c.r, c.g, c.b)
}

func ToColor16(c Color) Color16 {
	r := int(c.r)
	g := int(c.g)
	b := int(c.b)

	brightness := (r + g + b) / 3
	bright := brightness > 128

	if r >= g && r >= b {
		if g >= 128 {
			return Color16{Color8Yellow, bright}
		}
		if b >= 128 {
			return Color16{Color8Magenta, bright}
		}
		return Color16{Color8Red, bright}
	}
	if g >= r && g >= b {
		if r >= 128 {
			return Color16{Color8Yellow, bright}
		}
		if b >= 128 {
			return Color16{Color8Cyan, bright}
		}
		return Color16{Color8Green, bright}
	}
	if b >= r && b >= g {
		if r >= 128 {
			return Color16{Color8Magenta, bright}
		}
		if g >= 128 {
			return Color16{Color8Cyan, bright}
		}
		return Color16{Color8Blue, bright}
	}

	if brightness < 64 {
		return Color16{Color8Black, false}
	}
	if brightness >= 192 {
		return Color16{Color8White, true}
	}
	return Color16{Color8White, false}
}

func ToColor256(c Color) Color256 {
	r := (int(c.r) + 21) * 5 / 255
	g := (int(c.g) + 21) * 5 / 255
	b := (int(c.b) + 21) * 5 / 255
	return Color256(uint8(16 + 36*r + 6*g + b))
}

func ParseHexColor(s string) (Color, error) {
	if len(s) != 7 || s[0] != '#' {
		return Color{}, fmt.Errorf("invalid format")
	}

	r, err := strconv.ParseUint(s[1:3], 16, 8)
	if err != nil {
		return Color{}, err
	}

	g, err := strconv.ParseUint(s[3:5], 16, 8)
	if err != nil {
		return Color{}, err
	}

	b, err := strconv.ParseUint(s[5:7], 16, 8)
	if err != nil {
		return Color{}, err
	}

	return Color{
		r: uint8(r),
		g: uint8(g),
		b: uint8(b),
	}, nil
}
