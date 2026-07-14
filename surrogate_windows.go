//go:build windows

package termi

type surrogateState struct {
	pending uint16
	hasHigh bool
}

var winSurrogate surrogateState
var keySurrogate surrogateState

func isHighSurrogate(r uint16) bool {
	return 0xd800 <= r && r <= 0xdbff
}

func isLowSurrogate(r uint16) bool {
	return 0xdc00 <= r && r <= 0xdfff
}
