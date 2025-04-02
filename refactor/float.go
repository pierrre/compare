package compare

import (
	"reflect"
)

// FloatComparator is a [ValuesComparator] for float values.
type FloatComparator struct{}

// NewFloatComparator creates a new [FloatComparator].
func NewFloatComparator() *FloatComparator {
	return &FloatComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *FloatComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	switch v1.Kind() { //nolint:exhaustive // Only supports float.
	case reflect.Float32, reflect.Float64:
	default:
		return nil, false
	}
	f1 := v1.Float()
	f2 := v2.Float()
	if f1 == f2 {
		return nil, true
	}
	return st.Result("float not equal", f1, f2)
}
