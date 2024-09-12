package compare

import (
	"reflect"
	"strconv"
)

// PointerComparator is a [ValueComparator] that handles pointer values.
type PointerComparator struct {
	ValueComparator
}

// NewPointerComparator creates a new [PointerComparator].
func NewPointerComparator(vc ValueComparator) *PointerComparator {
	return &PointerComparator{
		ValueComparator: vc,
	}
}

// CompareValues implements [ValueComparator].
func (vc *PointerComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	if v1.Kind() != reflect.Ptr {
		return false
	}
	p1 := v1.Pointer()
	p2 := v2.Pointer()
	if p1 == p2 {
		return true
	}
	v1 = v1.Elem()
	v2 = v2.Elem()
	v1Nil := !v1.IsValid()
	v2Nil := !v2.IsValid()
	if v1Nil && v2Nil {
		return true
	}
	if v1Nil != v2Nil {
		st.addDifferences(Difference{
			Message: msgOnlyOneIsNil,
			V1:      strconv.FormatBool(!v1Nil),
			V2:      strconv.FormatBool(!v2Nil),
		})
		return true
	}
	// TODO append path to differences.
	return vc.ValueComparator(st, v1, v2)
}
