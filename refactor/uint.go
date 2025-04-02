package compare

import (
	"reflect"
)

type UintValuesComparator struct{}

func NewUintValuesComparator() *UintValuesComparator {
	return &UintValuesComparator{}
}

func (vc *UintValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		return false
	}
	ui1 := v1.Uint()
	ui2 := v2.Uint()
	if ui1 == ui2 {
		return true
	}
	st.Yield(Difference{
		Message: "uint not equal",
		V1:      v1,
		V2:      v2,
	})
	return true
}
