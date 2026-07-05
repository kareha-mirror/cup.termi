package rutil

func RuneIndex(line string, start int, r rune) int {
	col := 0
	for _, ru := range line {
		if col < start {
			col++
			continue
		}
		if ru == r {
			return col
		}
		col++
	}
	return -1
}

func LastRuneIndex(line string, start int, r rune) int {
	found := -1
	col := 0
	for _, ru := range line {
		if col > start {
			break
		}
		if ru == r {
			found = col
		}
		col++
	}
	return found
}
