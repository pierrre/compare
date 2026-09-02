package compare

import (
	"reflect"
)

// SliceComparator is a [ValuesComparator] for slice values.
//
// It should be created with [NewSliceComparator].
type SliceComparator struct {
	ValuesComparator
	// MaxDifferences is the maximum number of different items to report.
	// If reached, the comparison is stopped for the current slice.
	// Setting it to 0 disables it.
	// Default: 10.
	MaxDifferences int
}

// NewSliceComparator creates a new [SliceComparator].
func NewSliceComparator(vc ValuesComparator) *SliceComparator {
	return &SliceComparator{
		ValuesComparator: vc,
		MaxDifferences:   10,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *SliceComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Slice {
		return nil, false
	}
	res, handled = compareNilLenPointer(v1, v2)
	if handled {
		return res, true
	}
	return compareArray(vc.ValuesComparator, st, v1, v2, vc.MaxDifferences), true
}

// Supports implements [SupportChecker].
func (vc *SliceComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.Slice {
		res = vc
	}
	return res
}
