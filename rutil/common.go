package rutil

import (
	"strings"
)

func FilterByPrefix(list []string, prefix string) []string {
	result := []string{}
	for _, elem := range list {
		if strings.HasPrefix(elem, prefix) {
			result = append(result, elem)
		}
	}
	return result
}

func CommonPrefix(list []string) string {
	if len(list) == 0 {
		return ""
	}

	prefix := []rune(list[0])

	for _, s := range list[1:] {
		r := []rune(s)

		n := len(prefix)
		if len(r) < n {
			n = len(r)
		}

		i := 0
		for i < n && prefix[i] == r[i] {
			i++
		}

		prefix = prefix[:i]
		if len(prefix) == 0 {
			break
		}
	}

	return string(prefix)
}
