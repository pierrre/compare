package compare

import (
	"reflect"
	"strconv"
)

type FloatValuesComparator struct{}

func NewFloatValuesComparator() *FloatValuesComparator {
	return &FloatValuesComparator{}
}

func (vc *FloatValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool) {
	switch v1.Kind() {
	case reflect.Float32, reflect.Float64:
	default:
		return false, false
	}
	f1 := v1.Float()
	f2 := v2.Float()
	if f1 == f2 {
		return true, true
	}
	bitSize := v1.Type().Bits()
	return st.Yield(&Difference{
		Message: "float not equal",
		V1:      strconv.FormatFloat(f1, 'g', -1, bitSize),
		V2:      strconv.FormatFloat(f2, 'g', -1, bitSize),
	}), true
}
