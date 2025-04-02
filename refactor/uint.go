package compare

import (
	"reflect"
)

// UintComparator is a [ValuesComparator] for uint values.
type UintComparator struct{}

// NewUintComparator creates a new [UintComparator].
func NewUintComparator() *UintComparator {
	return &UintComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *UintComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	switch v1.Kind() { //nolint:exhaustive // Only supports uint.
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		return nil, false
	}
	ui1 := v1.Uint()
	ui2 := v2.Uint()
	if ui1 == ui2 {
		return nil, true
	}
	return st.Result("uint not equal", ui1, ui2)
}
