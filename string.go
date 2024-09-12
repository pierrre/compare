package compare

import (
	"reflect"
	"strconv"
)

// StringComparator is a [ValueComparator] that handles string values.
type StringComparator struct{}

// NewStringComparator creates a new [StringComparator].
func NewStringComparator() *StringComparator {
	return &StringComparator{}
}

// CompareValues implements [ValueComparator].
func (vc *StringComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	if v1.Kind() != reflect.String {
		return false
	}
	s1 := v1.String()
	s2 := v2.String()
	if s1 != s2 {
		st.addDifferences(Difference{
			Message: msgStringNotEqual,
			V1:      strconv.Quote(s1),
			V2:      strconv.Quote(s2),
		})
	}
	return true
}
