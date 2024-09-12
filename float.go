package compare

import (
	"reflect"
	"strconv"
)

// FloatValueComparator is a [ValueComparator] that handles float values.
type FloatValueComparator struct{}

// NewFloatValueComparator creates a new [FloatValueComparator].
func NewFloatValueComparator() *FloatValueComparator {
	return &FloatValueComparator{}
}

// CompareValues implements [ValueComparator].
func (vc *FloatValueComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() { //nolint:exhaustive // Only handles float.
	case reflect.Float32, reflect.Float64:
	default:
		return false
	}
	f1 := v1.Float()
	f2 := v2.Float()
	if f1 != f2 {
		bitSize := v1.Type().Bits()
		st.addDifferences(Difference{
			Message: msgFloatNotEqual,
			V1:      strconv.FormatFloat(f1, 'g', -1, bitSize),
			V2:      strconv.FormatFloat(f2, 'g', -1, bitSize),
		})
	}
	return true
}
