package termi

type StringBuilder struct {
	buf []rune
}

func (b *StringBuilder) WriteRune(r rune) {
	b.buf = append(b.buf, r)
}

func (b *StringBuilder) WriteString(s string) {
	b.buf = append(b.buf, []rune(s)...)
}

func (b *StringBuilder) String() string {
	return string(b.buf)
}

func (b *StringBuilder) Reset() {
	b.buf = b.buf[:0]
}

func (b *StringBuilder) Len() int {
	return len(b.buf)
}

func (b *StringBuilder) RemoveTail() bool {
	if len(b.buf) == 0 {
		return false
	}
	b.buf = b.buf[:len(b.buf)-1]
	return true
}

func (b *StringBuilder) RemoveHead() bool {
	if len(b.buf) == 0 {
		return false
	}
	b.buf = b.buf[1:]
	return true
}
