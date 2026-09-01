package compare

import (
	"reflect"
)

// StringComparator is a [ValuesComparator] for string values.
//
// It should be created with [NewStringComparator].
type StringComparator struct{}

// NewStringComparator creates a new [StringComparator].
func NewStringComparator() *StringComparator {
	return &StringComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *StringComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.String {
		return nil, false
	}
	s1 := v1.String()
	s2 := v2.String()
	if s1 == s2 {
		return nil, true
	}
	return Res("string not equal", s1, s2), true
}

// Supports implements [SupportChecker].
func (vc *StringComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.String {
		res = vc
	}
	return res
}
