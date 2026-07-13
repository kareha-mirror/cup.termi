package rutil

func ByteIndex(s string, col int) int {
	for i := range s {
		if col < 1 {
			return i
		}
		col--
	}
	return len(s)
}

func Head(s string, end int) string {
	return s[:ByteIndex(s, end)]
}

func Tail(s string, start int) string {
	i := ByteIndex(s, start)
	if i >= len(s) {
		return ""
	}
	return s[i:]
}

func Body(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end <= start {
		return ""
	}

	i := ByteIndex(s, start)
	if i >= len(s) {
		return ""
	}
	j := ByteIndex(s[i:], end-start)
	return s[i : i+j]
}

func Split(s string, col int) (string, string) {
	i := ByteIndex(s, col)
	if i >= len(s) {
		return s, ""
	}
	return s[:i], s[i:]
}

func SplitBody(s string, start, end int) (string, string, string) {
	if start < 0 {
		start = 0
	}
	if end <= start {
		head, tail := Split(s, start)
		return head, "", tail
	}

	i := ByteIndex(s, start)
	if i >= len(s) {
		return s, "", ""
	}
	body := s[i:]
	j := ByteIndex(body, end-start)
	return s[:i], body[:j], body[j:]
}
