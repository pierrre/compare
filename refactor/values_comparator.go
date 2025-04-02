package compare

import (
	"reflect"
)

// ValuesComparator compares 2 values.
// The returned handled value is true if the comparator handled the value, and false otherwise.
//
// Implementations can assume that both values are valid and have the same type.
type ValuesComparator interface {
	CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool)
}

// ValuesComparatorFunc is a [ValuesComparator] function.
type ValuesComparatorFunc func(st *State, v1, v2 reflect.Value) (res Result, handled bool)

// CompareValues implements [ValuesComparator].
func (f ValuesComparatorFunc) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	return f(st, v1, v2)
}
