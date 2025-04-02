package compare

import (
	"reflect"
)

// BoolComparator is a [ValuesComparator] for bool values.
type BoolComparator struct{}

// NewBoolComparator creates a new [BoolComparator].
func NewBoolComparator() *BoolComparator {
	return &BoolComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *BoolComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Bool {
		return nil, false
	}
	b1 := v1.Bool()
	b2 := v2.Bool()
	if b1 == b2 {
		return nil, true
	}
	return st.Result("bool not equal", b1, b2)
}
