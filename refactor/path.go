package compare

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/pierrre/go-libs/bytesutil"
)

// Path represents a path in a data structure.
type Path []PathElem

// Clone clones the path.
func (p Path) Clone() Path {
	return slices.Clone(p)
}

// Push pushes a [PathElem] to the path.
func (p *Path) Push(e PathElem) {
	*p = append(*p, e)
}

// Pop pops a [PathElem] from the path.
func (p *Path) Pop() PathElem {
	pv := *p
	i := len(pv) - 1
	e := pv[i]
	pv[i] = PathElem{}
	*p = pv[:i]
	return e
}

// AppendTo appends the text representation to a []byte.
func (p Path) AppendTo(b []byte) []byte {
	if len(p) == 0 {
		return append(b, '.')
	}
	for _, e := range p {
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
	Kind  PathElemKind
	Index int
	Key   reflect.Value
	Name  string
}

// AppendTo appends the text representation to a []byte.
func (e PathElem) AppendTo(b []byte) []byte {
	switch e.Kind {
	case PathElemKindIndex:
		b = append(b, '[')
		b = strconv.AppendInt(b, int64(e.Index), 10)
		b = append(b, ']')
	case PathElemKindKey:
		b = append(b, '[')
		b = fmt.Append(b, e.Key)
		b = append(b, ']')
	case PathElemKindName:
		b = append(b, '.')
		b = append(b, e.Name...)
	case PathElemKindPointer:
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

// PathElemKind represents the kind of a [PathElem].
type PathElemKind byte

// Values of PathElemKind.
const (
	PathElemKindIndex PathElemKind = iota + 1
	PathElemKindKey
	PathElemKindName
	PathElemKindPointer
)

type appender interface {
	AppendTo(b []byte) []byte
}

func stringFromAppender(a appender) string {
	bw := bytesWriterPool.Get()
	defer bytesWriterPool.Put(bw)
	*bw = a.AppendTo(*bw)
	return bw.String()
}

var bytesWriterPool = &bytesutil.WriterPool{
	MaxCap: -1,
}
