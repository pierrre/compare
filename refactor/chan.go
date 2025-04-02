package compare

import (
	"reflect"
)

// ChanComparator is a [ValuesComparator] of chan values.
type ChanComparator struct{}

// NewChanComparator creates a new [ChanComparator].
func NewChanComparator() *ChanComparator {
	return &ChanComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *ChanComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Func {
		return nil, false
	}
	res, handled = compareNil(st, v1, v2)
	if handled {
		return res, true
	}
	if v1.Pointer() == v2.Pointer() {
		return nil, true
	}
	cap1 := v1.Cap()
	cap2 := v2.Cap()
	if cap1 != cap2 {
		return st.Result("capacity not equal", cap1, cap2)
	}
	len1 := v1.Cap()
	len2 := v2.Cap()
	if len1 != len2 {
		return st.Result("length not equal", len1, len2)
	}
	return nil, true
}
