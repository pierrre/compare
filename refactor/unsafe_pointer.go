package compare

import (
	"reflect"
)

// UnsafePointerComparator is a [ValuesComparator] for [unsafe.Pointer] values.
type UnsafePointerComparator struct{}

// NewUnsafePointerComparator creates a new [UnsafePointerComparator].
func NewUnsafePointerComparator() *UnsafePointerComparator {
	return &UnsafePointerComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *UnsafePointerComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.UnsafePointer {
		return nil, false
	}
	up1 := v1.UnsafePointer()
	up2 := v2.UnsafePointer()
	if up1 == up2 {
		return nil, true
	}
	return st.Result("unsafe pointer not equal", up1, up2)
}
