package compare

import (
	"reflect"
)

// SliceComparator is a [ValuesComparator] for slice values.
type SliceComparator struct {
	ValuesComparator
}

// NewSliceComparator creates a new [SliceComparator].
func NewSliceComparator(vc ValuesComparator) *SliceComparator {
	return &SliceComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *SliceComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Slice {
		return nil, false
	}
	res, handled = compareNilPointerLen(st, v1, v2)
	if handled {
		return res, true
	}
	return compareArray(vc.ValuesComparator, st, v1, v2), true
}
