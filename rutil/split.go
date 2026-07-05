package rutil

func Index(s string, col int) int {
	for i, _ := range s {
		if col < 1 {
			return i
		}
		col--
	}
	return len(s)
}

func Head(s string, end int) string {
	return s[:Index(s, end)]
}

func Tail(s string, start int) string {
	i := Index(s, start)
	if i >= len(s) {
		return ""
	}
	return s[i:]
}

func Body(s string, start, end int) string {
	i := Index(s, start)
	if i >= len(s) {
		return ""
	}
	j := Index(s[i:], end-start)
	return s[i : i+j]
}

func Split(s string, col int) (string, string) {
	i := Index(s, col)
	if i >= len(s) {
		return s, ""
	}
	return s[:i], s[i:]
}

func SplitBody(s string, start, end int) (string, string, string) {
	i := Index(s, start)
	if i >= len(s) {
		return s, "", ""
	}
	body := s[i:]
	j := Index(body, end-start)
	if j >= len(body) {
		return s[:i], body, ""
	}
	return s[:i], body[:j], body[j:]
}
