package compare

import (
	"reflect"
	"strconv"
)

type IntValuesComparator struct{}

func NewIntValuesComparator() *IntValuesComparator {
	return &IntValuesComparator{}
}

func (vc *IntValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool) {
	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return false, false
	}
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 == i2 {
		return true, true
	}
	return st.Yield(&Difference{
		Message: "int not equal",
		V1:      strconv.FormatInt(i1, 10),
		V2:      strconv.FormatInt(i2, 10),
	}), true
}
