package compare

import (
	"reflect"
)

// UnsafePointerComparator is a [ValuesComparator] for [unsafe.Pointer] values.
//
// It should be created with [NewUnsafePointerComparator].
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
	return Res("unsafe pointer not equal", up1, up2), true
}

// Supports implements [SupportChecker].
func (vc *UnsafePointerComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.UnsafePointer {
		res = vc
	}
	return res
}
