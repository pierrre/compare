package compare

import (
	"reflect"
	"strconv"
)

type BoolValueComparator struct{}

func (vc *BoolValueComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	if v1.Kind() != reflect.Bool {
		return false
	}
	b1 := v1.Bool()
	b2 := v2.Bool()
	if b1 != b2 {
		st.addDifferences(Difference{
			Message: msgBoolNotEqual,
			V1:      strconv.FormatBool(b1),
			V2:      strconv.FormatBool(b2),
		})
	}
	return true
}
