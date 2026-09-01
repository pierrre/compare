package compare

import (
	"bytes"
	"reflect"
)

var bytesType = reflect.TypeFor[[]byte]()

// BytesComparator is a [ValuesComparator] for byte slice values ([]byte).
//
// It should be created with [NewBytesComparator].
type BytesComparator struct{}

// NewBytesComparator creates a new [BytesComparator].
func NewBytesComparator() *BytesComparator {
	return &BytesComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *BytesComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Slice || v1.Type() != bytesType {
		return nil, false
	}
	b1 := v1.Bytes()
	b2 := v2.Bytes()
	if bytes.Equal(b1, b2) {
		return nil, true
	}
	return Res("bytes not equal", b1, b2), true
}

// Supports implements [SupportChecker].
func (vc *BytesComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ == bytesType {
		res = vc
	}
	return res
}
