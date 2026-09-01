package compare

import (
	"fmt"
	"slices"
	"strconv"
)

// Path represents a path in a data structure.
//
// Elements are stored in reverse order; the first element is the deepest.
// It helps to prepend elements to the path efficiently.
type Path []PathElem

// AppendTo appends the text representation to a []byte.
func (p Path) AppendTo(b []byte) []byte {
	if len(p) == 0 {
		return append(b, '.')
	}
	for _, e := range slices.Backward(p) {
		b = e.AppendTo(b)
	}
	return b
}

// AppendText implements [encoding.TextAppender].
func (p Path) AppendText(b []byte) ([]byte, error) {
	return p.AppendTo(b), nil
}

func (p Path) String() string {
	return stringFromAppender(p)
}

// PathElem represents an element in a [Path].
type PathElem struct {
	Index   *int
	Key     *any
	Name    *string
	Pointer bool
}

// AppendTo appends the text representation to a []byte.
func (e PathElem) AppendTo(b []byte) []byte {
	switch {
	case e.Index != nil:
		b = append(b, '[')
		b = strconv.AppendInt(b, int64(*e.Index), 10)
		b = append(b, ']')
	case e.Key != nil:
		b = append(b, '[')
		b = fmt.Append(b, *e.Key)
		b = append(b, ']')
	case e.Name != nil:
		b = append(b, '.')
		b = append(b, *e.Name...)
	case e.Pointer:
		b = append(b, '*')
	}
	return b
}

// AppendText implements [encoding.TextAppender].
func (e PathElem) AppendText(b []byte) ([]byte, error) {
	return e.AppendTo(b), nil
}

func (e PathElem) String() string {
	return stringFromAppender(e)
}
