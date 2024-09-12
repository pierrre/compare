package compare

import (
	"fmt"
)

// Result is a list of [Difference].
type Result []Difference

// Format implements [fmt.Formatter].
//
// See [Difference.Format] for supported verb and flag.
func (r Result) Format(s fmt.State, verb rune) {
	if verb != 'v' {
		_, _ = fmt.Fprintf(s, "%%!%c(%T)", verb, r)
		return
	}
	if len(r) == 0 {
		_, _ = s.Write(resultNoneBytes)
		return
	}
	for i, d := range r {
		if i > 0 {
			_, _ = s.Write(resultNewLineBytes)
		}
		d.Format(s, verb)
	}
}

func (r Result) pathAppend(pe PathElem) {
	for i := range r {
		r[i].Path = append(r[i].Path, pe)
	}
}

var (
	resultNoneBytes    = []byte("<none>")
	resultNewLineBytes = []byte("\n")
)
