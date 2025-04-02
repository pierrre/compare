package compare

import (
	"reflect"
	"strconv"
)

type UintValuesComparator struct{}

func NewUintValuesComparator() *UintValuesComparator {
	return &UintValuesComparator{}
}

func (vc *UintValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool) {
	switch v1.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return false, false
	}
	ui1 := v1.Uint()
	ui2 := v2.Uint()
	if ui1 == ui2 {
		return true, true
	}
	return st.Yield(&Difference{
		Message: "uint not equal",
		V1:      strconv.FormatUint(ui1, 10),
		V2:      strconv.FormatUint(ui2, 10),
	}), true
}
