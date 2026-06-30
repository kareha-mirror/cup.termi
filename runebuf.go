package termi

type RuneBuf struct {
	buf []rune
	s   string
}

func (b *RuneBuf) WriteRune(r rune) {
	b.buf = append(b.buf, r)
	b.s = ""
}

func (b *RuneBuf) WriteString(s string) {
	b.buf = append(b.buf, []rune(s)...)
	b.s = ""
}

func (b *RuneBuf) String() string {
	if b.s == "" {
		b.s = string(b.buf)
	}
	return b.s
}

func (b *RuneBuf) Reset() {
	b.buf = b.buf[:0]
	b.s = ""
}

func (b *RuneBuf) RuneCount() int {
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
	b.s = ""
	return true
}

func (b *RuneBuf) RemoveTail() bool {
	if len(b.buf) == 0 {
		return false
	}
	b.buf = b.buf[:len(b.buf)-1]
	b.s = ""
	return true
}

func (b *RuneBuf) Body(from, to int) RuneBuf {
	if from < 0 || from > len(b.buf)-1 {
		return RuneBuf{}
	}
	if to < 0 || to > len(b.buf) {
		return RuneBuf{}
	}
	if to < from {
		return RuneBuf{}
	}
	return RuneBuf{buf: append([]rune{}, b.buf[from:to]...)}
}
