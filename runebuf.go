package termi

type RuneBuf struct {
	buf []rune
}

func (b *RuneBuf) WriteRune(r rune) {
	b.buf = append(b.buf, r)
}

func (b *RuneBuf) WriteString(s string) {
	b.buf = append(b.buf, []rune(s)...)
}

func (b *RuneBuf) String() string {
	return string(b.buf)
}

func (b *RuneBuf) Reset() {
	b.buf = b.buf[:0]
}

func (b *RuneBuf) Len() int {
	return len(b.buf)
}

func (b *RuneBuf) Head() (rune, bool) {
	if len(b.buf) < 1 {
		return 0, false
	}
	return b.buf[0], true
}

func (b *RuneBuf) Tail() (rune, bool) {
	if len(b.buf) < 1 {
		return 0, false
	}
	return b.buf[len(b.buf)-1], true
}

func (b *RuneBuf) RemoveHead() bool {
	if len(b.buf) == 0 {
		return false
	}
	b.buf = b.buf[1:]
	return true
}

func (b *RuneBuf) RemoveTail() bool {
	if len(b.buf) == 0 {
		return false
	}
	b.buf = b.buf[:len(b.buf)-1]
	return true
}

func (b *RuneBuf) Substring(from, to int) (string, bool) {
	if from < 0 || from > len(b.buf)-1 {
		return "", false
	}
	if to < 0 || to > len(b.buf) {
		return "", false
	}
	if to < from {
		return "", false
	}
	return string(b.buf[from:to]), true
}
