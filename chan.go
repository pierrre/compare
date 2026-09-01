package compare

import (
	"reflect"
)

// ChanComparator is a [ValuesComparator] of chan values.
//
// It should be created with [NewChanComparator].
type ChanComparator struct{}

// NewChanComparator creates a new [ChanComparator].
func NewChanComparator() *ChanComparator {
	return &ChanComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *ChanComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Chan {
		return nil, false
	}
	res, handled = compareNil(v1, v2)
	if handled {
		return res, true
	}
	if v1.Pointer() == v2.Pointer() {
		return nil, true
	}
	cap1 := v1.Cap()
	cap2 := v2.Cap()
	if cap1 != cap2 {
		return Res("capacity not equal", cap1, cap2), true
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return Res("length not equal", len1, len2), true
	}
	return nil, true
}

// Supports implements [SupportChecker].
func (vc *ChanComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.Chan {
		res = vc
	}
	return res
}
