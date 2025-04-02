package compare

import (
	"reflect"
)

type BoolValuesComparator struct{}

func NewBoolValuesComparator() *BoolValuesComparator {
	return &BoolValuesComparator{}
}

func (vc *BoolValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	if v1.Kind() != reflect.Bool {
		return false
	}
	b1 := v1.Bool()
	b2 := v2.Bool()
	if b1 == b2 {
		return true
	}
	st.Yield(Difference{
		Message: "bool not equal",
		V1:      v1,
		V2:      v2,
	})
	return true
}
