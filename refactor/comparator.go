package compare

import (
	"reflect"
)

// Comparator is a value comparator.
type Comparator struct {
	ValuesComparator ValuesComparator
}

// Compare compares 2 values.
func (c *Comparator) Compare(iv1, iv2 any) Result {
	v1 := reflect.ValueOf(iv1)
	v2 := reflect.ValueOf(iv2)
	vl1 := v1.IsValid()
	vl2 := v2.IsValid()
	if !vl1 && !vl2 {
		return nil
	}
	if vl1 != vl2 {
		return Result{{
			Message: "nil mismatch",
			V1:      !vl1,
			V2:      !vl2,
		}}
	}
	st := &State{} // TODO recycle
	res, handled := compareType(st, v1, v2)
	if handled {
		return res
	}
	res, _ = c.ValuesComparator.CompareValues(st, v1, v2) // TODO check handled ?
	return res
}
