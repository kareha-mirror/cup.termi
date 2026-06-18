package termi

type Sig int

const (
	SigStop Sig = iota
	SigCont
)

var sigCh chan Sig

func Sigs() chan Sig {
	return sigCh
}
