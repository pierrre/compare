package compare

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"reflect"

	"github.com/pierrre/go-libs/bufpool"
	"github.com/pierrre/go-libs/strconvio"
)

type Comparator struct {
	ValuesComparator ValuesComparator
}

func (c *Comparator) Compare(iv1, iv2 any) iter.Seq[Difference] {
	return func(yield func(Difference) bool) {
		v1 := reflect.ValueOf(iv1)
		v2 := reflect.ValueOf(iv2)
		st := &State{
			yield:   yield,
			YieldOK: true,
		}
		c.ValuesComparator.CompareValues(st, v1, v2)
	}
}

type State struct {
	yield   func(Difference) bool
	YieldOK bool
	Path    Path
}

func (st *State) Yield(diff Difference) {
	diff.Path = st.Path
	st.YieldOK = st.yield(diff)
}

type ValuesComparator interface {
	CompareValues(st *State, v1, v2 reflect.Value) (handled bool)
}

type ValuesComparatorFunc func(st *State, v1, v2 reflect.Value) (handled bool)

func (f ValuesComparatorFunc) CompareValues(st *State, v1, v2 reflect.Value) (handled bool) {
	return f(st, v1, v2)
}

type Difference struct {
	Path    Path
	Message string
	V1, V2  reflect.Value
}

func (d Difference) WriteToBuffer(buf *bytes.Buffer) {
	d.Path.WriteToBuffer(buf)
	_, _ = buf.WriteString(": ")
	_, _ = buf.WriteString(d.Message)
	_ = buf.WriteByte('\n')
	_, _ = fmt.Fprint(buf, d.V1)
	_ = buf.WriteByte('\n')
	_, _ = fmt.Fprint(buf, d.V2)
}

func (d Difference) String() string {
	return StringWithBuffer(d)
}

func (d Difference) WriteTo(w io.Writer) (n int64, err error) {
	return WriteToWithBuffer(w, d)
}

type Path []PathElem

func (p *Path) Push(e PathElem) {
	*p = append(*p, e)
}

func (p *Path) Pop() {
	pv := *p
	l := len(pv)
	pv[l-1] = PathElem{}
	*p = pv[:l-1]
}

func (p Path) WriteToBuffer(buf *bytes.Buffer) {
	if len(p) == 0 {
		_ = buf.WriteByte('.')
		return
	}
	for _, e := range p {
		e.WriteToBuffer(buf)
	}
}

func (p Path) String() string {
	return StringWithBuffer(p)
}

func (p Path) WriteTo(w io.Writer) (n int64, err error) {
	return WriteToWithBuffer(w, p)
}

type PathElem struct {
	Kind  PathElemKind
	Index int
	Key   reflect.Value
	Name  string
}

func (e PathElem) WriteToBuffer(buf *bytes.Buffer) {
	switch e.Kind {
	case PathElemKindIndex:
		_ = buf.WriteByte('[')
		_, _ = strconvio.WriteInt(buf, int64(e.Index), 10)
		_ = buf.WriteByte(']')
	case PathElemKindKey:
		_ = buf.WriteByte('[')
		_, _ = fmt.Fprint(buf, e.Key)
		_ = buf.WriteByte(']')
	case PathElemKindName:
		_ = buf.WriteByte('.')
		_, _ = buf.WriteString(e.Name)
	case PathElemKindPointer:
		_ = buf.WriteByte('*')
	}
}

func (e PathElem) String() string {
	return StringWithBuffer(e)
}

func (e PathElem) WriteTo(w io.Writer) (n int64, err error) {
	return WriteToWithBuffer(w, e)
}

type PathElemKind byte

const (
	PathElemKindIndex PathElemKind = iota + 1
	PathElemKindKey
	PathElemKindName
	PathElemKindPointer
)

type WriterToBuffer interface {
	WriteToBuffer(buf *bytes.Buffer)
}

func WriteToWithBuffer(w io.Writer, wb WriterToBuffer) (n int64, err error) {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	wb.WriteToBuffer(buf)
	return buf.WriteTo(w)
}

func StringWithBuffer(wb WriterToBuffer) string {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	wb.WriteToBuffer(buf)
	return buf.String()
}

var bufPool bufpool.Pool
