package compare

import (
	"reflect"
)

// PointerComparator is a [ValuesComparator] for pointer values.
type PointerComparator struct {
	ValuesComparator
}

// NewPointerComparator creates a new [PointerComparator].
func NewPointerComparator(vc ValuesComparator) *PointerComparator {
	return &PointerComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *PointerComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Pointer {
		return nil, false
	}
	res, handled = compareNil(st, v1, v2)
	if handled {
		return res, true
	}
	if v1.Pointer() == v2.Pointer() {
		return nil, true
	}
	res, _ = vc.ValuesComparator.CompareValues(st, v1.Elem(), v2.Elem())
	if len(res) != 0 {
		res.AppendPathElem(PathElem{
			Pointer: true,
		})
	}
	return res, true
}
