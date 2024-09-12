package compare

import (
	"reflect"
	"strconv"
)

// IntValueComparator is a [ValueComparator] that handles int values.
type IntValueComparator struct{}

// NewIntValueComparator creates a new [IntValueComparator].
func NewIntValueComparator() *IntValueComparator {
	return &IntValueComparator{}
}

// CompareValues implements [ValueComparator].
func (vc *IntValueComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() { //nolint:exhaustive // Only handles int.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return false
	}
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 != i2 {
		st.addDifferences(Difference{
			Message: msgIntNotEqual,
			V1:      strconv.FormatInt(i1, 10),
			V2:      strconv.FormatInt(i2, 10),
		})
	}
	return true
}
