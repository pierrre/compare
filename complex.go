package compare

import (
	"reflect"
	"strconv"
)

// ComplexValueComparator is a [ValueComparator] that handles complex values.
type ComplexValueComparator struct{}

// NewComplexValueComparator creates a new [ComplexValueComparator].
func NewComplexValueComparator() *ComplexValueComparator {
	return &ComplexValueComparator{}
}

// CompareValues implements [ValueComparator].
func (vc *ComplexValueComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() { //nolint:exhaustive // Only handles complex.
	case reflect.Complex64, reflect.Complex128:
	default:
		return false
	}
	c1 := v1.Complex()
	c2 := v2.Complex()
	if c1 != c2 {
		bitSize := v1.Type().Bits()
		st.addDifferences(Difference{
			Message: msgComplexNotEqual,
			V1:      strconv.FormatComplex(c1, 'g', -1, bitSize),
			V2:      strconv.FormatComplex(c2, 'g', -1, bitSize),
		})
	}
	return true
}
