package compare

import (
	"reflect"
)

// IntComparator is a [ValuesComparator] for int values.
type IntComparator struct{}

// NewIntComparator creates a new [IntComparator].
func NewIntComparator() *IntComparator {
	return &IntComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *IntComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	switch v1.Kind() { //nolint:exhaustive // Only supports int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return nil, false
	}
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 == i2 {
		return nil, true
	}
	return st.Result("int not equal", i1, i2)
}
