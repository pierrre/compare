package compare

import (
	"strconv"
	"strings"
)

// Path represents a field path, which is a list of PathElem.
//
// Elements are stored in reverse order, the first element is the deepest.
// It helps to prepend elements to the path efficiently.
type Path []PathElem

// String returns the string value for a Path.
func (p Path) String() string {
	if len(p) == 0 {
		return "."
	}
	ss := make([]string, len(p))
	for i, e := range p {
		ss[len(ss)-i-1] = e.String()
	}
	return strings.Join(ss, "")
}

// PathElem is a single element in a Path.
type PathElem interface {
	String() string
	pathElem()
}

// StructPathElem is a PathElem for a struct field.
type StructPathElem struct {
	Field string
}

// String returns the string representation.
func (e StructPathElem) String() string {
	return "." + e.Field
}

func (e StructPathElem) pathElem() {}

// MapPathElem is a PathElem for a map key.
type MapPathElem struct {
	Key string
}

// String returns the string representation.
func (e MapPathElem) String() string {
	return "[" + e.Key + "]"
}

func (e MapPathElem) pathElem() {}

// IndexedPathElem is a PathElem for a slice/array index.
type IndexedPathElem struct {
	Index int
}

// String returns the string representation.
func (e IndexedPathElem) String() string {
	return "[" + strconv.Itoa(e.Index) + "]"
}

func (e IndexedPathElem) pathElem() {}
