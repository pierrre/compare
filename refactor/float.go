package compare

import (
	"reflect"
)

type FloatValuesComparator struct{}

func NewFloatValuesComparator() *FloatValuesComparator {
	return &FloatValuesComparator{}
}

func (vc *FloatValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Float32, reflect.Float64:
	default:
		return false
	}
	f1 := v1.Float()
	f2 := v2.Float()
	if f1 == f2 {
		return true
	}
	st.Yield(Difference{
		Message: "float not equal",
		V1:      v1,
		V2:      v2,
	})
	return true
}
