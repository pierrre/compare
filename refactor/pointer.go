package compare

import "reflect"

type PointerValuesComparator struct {
	ValuesComparator
}

func NewPointerValuesComparator(vc ValuesComparator) *PointerValuesComparator {
	return &PointerValuesComparator{
		ValuesComparator: vc,
	}
}

func (vc *PointerValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) bool {
	if v1.Kind() != reflect.Pointer {
		return false
	}
	p1 := v1.Pointer()
	p2 := v2.Pointer()
	if p1 == p2 {
		return true
	}
	if p1 == 0 || p2 == 0 {
		st.Yield(Difference{
			Message: "pointer nil mismatch",
			V1:      v1,
			V2:      v2,
		})
		return true
	}
	st.Path = append(st.Path, PathElem{
		Kind: PathElemKindPointer,
	})
	defer func() {
		st.Path = st.Path[:len(st.Path)-1]
	}()
	return vc.ValuesComparator.CompareValues(st, v1.Elem(), v2.Elem())
}
