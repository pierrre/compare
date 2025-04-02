package compare

import (
	"reflect"
)

// ComplexComparator is a [ValuesComparator] for complex values.
type ComplexComparator struct{}

// NewComplexComparator creates a new [ComplexComparator].
func NewComplexComparator() *ComplexComparator {
	return &ComplexComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *ComplexComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	switch v1.Kind() { //nolint:exhaustive // Only supports complex.
	case reflect.Complex64, reflect.Complex128:
	default:
		return nil, false
	}
	c1 := v1.Complex()
	c2 := v2.Complex()
	if c1 == c2 {
		return nil, true
	}
	return st.Result("complex not equal", c1, c2)
}
