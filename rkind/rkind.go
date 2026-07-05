package rkind

func IsBlank(r rune) bool {
	return r == ' ' || r == '\t'
}

func IsWord(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		r == '_' ||
		(r >= '0' && r <= '9') ||
		(r >= 'A' && r <= 'Z')
}

func IsSymbol(r rune) bool {
	return !IsBlank(r) && !IsWord(r) && r < 0x80
}

func IsExtraSymbol(r rune) bool {
	switch r {
	case '。', '、', '？', '！', '「', '」':
		return true
	default:
		return false
	}
}

func IsOther(r rune) bool {
	return r >= 0x80 && !IsExtraSymbol(r)
}

type RuneKind int

const (
	Blank RuneKind = iota
	Word
	Symbol
	ExtraSymbol
	Other
)

func Kind(r rune) RuneKind {
	// not optimized
	// don't optimize yet
	switch {
	case IsBlank(r):
		return Blank
	case IsWord(r):
		return Word
	case IsSymbol(r):
		return Symbol
	case IsExtraSymbol(r):
		return ExtraSymbol
	default:
		return Other
	}
}
