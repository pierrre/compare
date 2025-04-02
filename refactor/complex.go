package compare

import (
	"reflect"
)

type ComplexValuesComparator struct{}

func NewComplexValuesComparator() *ComplexValuesComparator {
	return &ComplexValuesComparator{}
}

func (vc *ComplexValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Complex64, reflect.Complex128:
	default:
		return false
	}
	c1 := v1.Complex()
	c2 := v2.Complex()
	if c1 == c2 {
		return true
	}
	st.Yield(Difference{
		Message: "complex not equal",
		V1:      v1,
		V2:      v2,
	})
	return true
}
