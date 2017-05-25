package compare

import (
	"bytes"
	"strconv"
)

// Path represents a field path.
type Path interface {
	PathString() string
	PathNext() Path
}

// PathString returns the string value for a Path.
func PathString(p Path) string {
	if p == nil {
		return "."
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	for p != nil {
		_, _ = buf.WriteString(p.PathString())
		p = p.PathNext()
	}
	s := buf.String()
	bufPool.Put(buf)
	return s
}

// StructPath is a Path for a struct field.
type StructPath struct {
	Field string
	Next  Path
}

// PathString implements PathItem
func (p StructPath) PathString() string {
	return "." + p.Field
}

// PathNext implements PathItem
func (p StructPath) PathNext() Path {
	return p.Next
}

func (p StructPath) String() string {
	return PathString(p)
}

// MapPath is a Path for a map key.
type MapPath struct {
	Key  string
	Next Path
}

// PathString implements PathItem
func (p MapPath) PathString() string {
	return "[" + p.Key + "]"
}

// PathNext implements PathItem
func (p MapPath) PathNext() Path {
	return p.Next
}

func (p MapPath) String() string {
	return PathString(p)
}

// IndexedPath is a Path for a slice/array index.
type IndexedPath struct {
	Index int
	Next  Path
}

// PathString implements PathItem
func (p IndexedPath) PathString() string {
	return "[" + strconv.Itoa(p.Index) + "]"
}

// PathNext implements PathItem
func (p IndexedPath) PathNext() Path {
	return p.Next
}

func (p IndexedPath) String() string {
	return PathString(p)
}
