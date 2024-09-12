package compare

import (
	"fmt"

	"github.com/pierrre/go-libs/strconvio"
	"github.com/pierrre/go-libs/unsafeio"
)

// Path represents a field path, which is a list of [PathElem].
//
// Elements are stored in reverse order, the first element is the deepest.
// It helps to prepend elements to the path efficiently.
type Path []PathElem

// Format implements [fmt.Formatter].
//
// It only supports the 'v' verb.
func (p Path) Format(s fmt.State, verb rune) {
	if len(p) == 0 {
		_, _ = unsafeio.WriteString(s, ".")
		return
	}
	for i := len(p) - 1; i >= 0; i-- {
		p[i].Format(s, verb)
	}
}

// PathElem is a single element in a [Path].
type PathElem struct {
	Struct *string `json:"struct,omitempty"`
	Map    *string `json:"map,omitempty"`
	Index  *int    `json:"index,omitempty"`
}

// Format implements [fmt.Formatter].
//
// It only supports the 'v' verb.
func (e PathElem) Format(s fmt.State, verb rune) {
	switch {
	case e.Struct != nil:
		_, _ = unsafeio.WriteString(s, ".")
		_, _ = unsafeio.WriteString(s, *e.Struct)
	case e.Map != nil:
		_, _ = unsafeio.WriteString(s, "[")
		_, _ = unsafeio.WriteString(s, *e.Map)
		_, _ = unsafeio.WriteString(s, "]")
	case e.Index != nil:
		_, _ = unsafeio.WriteString(s, "[")
		_, _ = strconvio.WriteInt(s, int64(*e.Index), 10)
		_, _ = unsafeio.WriteString(s, "]")
	}
}

func toPtr[V any](v V) *V {
	return &v
}
