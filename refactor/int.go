package compare

import (
	"reflect"
)

type IntValuesComparator struct{}

func NewIntValuesComparator() *IntValuesComparator {
	return &IntValuesComparator{}
}

func (vc *IntValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return false
	}
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 == i2 {
		return true
	}
	st.Yield(Difference{
		Message: "int not equal",
		V1:      v1,
		V2:      v2,
	})
	return true
}
