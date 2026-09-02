package compare

import (
	"fmt"
)

// Result represents the result of the comparison of 2 values.
// It's a list of [Difference].
type Result []Difference

// Res creates a new [Result] with a single [Difference].
func Res(msg string, v1, v2 any) Result {
	return Result{{
		Message: msg,
		V1:      v1,
		V2:      v2,
	}}
}

// AppendPathElem appends a [PathElem] to the [Path] of all [Difference]s of the [Result].
func (r Result) AppendPathElem(e PathElem) {
	for i := range r {
		r[i].Path = append(r[i].Path, e)
	}
}

// AppendTo appends the text representation to a []byte.
func (r Result) AppendTo(b []byte) []byte {
	for _, d := range r {
		b = d.AppendTo(b)
		b = append(b, '\n')
	}
	return b
}

// AppendText implements [encoding.TextAppender].
func (r Result) AppendText(b []byte) ([]byte, error) {
	return r.AppendTo(b), nil
}

func (r Result) String() string {
	return stringFromAppender(r)
}

// Difference represents a difference between 2 values.
type Difference struct {
	Path    Path
	Message string
	V1, V2  any
}

// AppendTo appends the text representation to a []byte.
func (d Difference) AppendTo(b []byte) []byte {
	b = d.Path.AppendTo(b)
	b = append(b, ": "...)
	b = append(b, d.Message...)
	b = append(b, "\nv1 = "...)
	b = fmt.Append(b, d.V1)
	b = append(b, "\nv2 = "...)
	b = fmt.Append(b, d.V2)
	return b
}

// AppendText implements [encoding.TextAppender].
func (d Difference) AppendText(b []byte) ([]byte, error) {
	return d.AppendTo(b), nil
}

func (d Difference) String() string {
	return stringFromAppender(d)
}
