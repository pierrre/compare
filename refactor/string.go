package compare

import "reflect"

type StringValuesComparator struct{}

func NewStringValuesComparator() *StringValuesComparator {
	return &StringValuesComparator{}
}

func (vc *StringValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	if v1.Kind() != reflect.String {
		return false
	}
	s1 := v1.String()
	s2 := v2.String()
	if s1 == s2 {
		return true
	}
	st.Yield(Difference{
		Message: "string not equal",
		V1:      v1,
		V2:      v2,
	})
	return true
}
