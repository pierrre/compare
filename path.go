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
type PathElem struct {
	Struct *string
	Map    *string
	Index  *int
}

func (e PathElem) String() string {
	if e.Struct != nil {
		return "." + *e.Struct
	}
	if e.Map != nil {
		return "[" + *e.Map + "]"
	}
	if e.Index != nil {
		return "[" + strconv.Itoa(*e.Index) + "]"
	}
	return ""
}
