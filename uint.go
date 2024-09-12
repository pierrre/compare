package compare

import (
	"reflect"
	"strconv"
)

// UintValueComparator is a [ValueComparator] that handles uint values.
type UintValueComparator struct{}

// NewUintValueComparator creates a new [UintValueComparator].
func NewUintValueComparator() *UintValueComparator {
	return &UintValueComparator{}
}

// CompareValues implements [ValueComparator].
func (vc *UintValueComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() { //nolint:exhaustive // Only handles uint.
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	default:
		return false
	}
	u1 := v1.Uint()
	u2 := v2.Uint()
	if u1 != u2 {
		st.addDifferences(Difference{
			Message: msgUintNotEqual,
			V1:      strconv.FormatUint(u1, 10),
			V2:      strconv.FormatUint(u2, 10),
		})
	}
	return true
}
