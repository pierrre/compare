package compare

import (
	"reflect"
)

// ArrayComparator is a [ValuesComparator] for array values.
type ArrayComparator struct {
	ValuesComparator
}

// NewArrayComparator creates a new [ArrayComparator].
func NewArrayComparator(vc ValuesComparator) *ArrayComparator {
	return &ArrayComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *ArrayComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Array {
		return nil, false
	}
	return compareArray(vc.ValuesComparator, st, v1, v2), true
}

func compareArray(vc ValuesComparator, st *State, v1, v2 reflect.Value) (res Result) {
	for i := range v1.Len() {
		r := compareArrayIndex(vc, st, v1, v2, i)
		if len(r) > 0 {
			// TODO: max diff
			res = append(res, r...)
		}
	}
	return res
}

func compareArrayIndex(vc ValuesComparator, st *State, v1, v2 reflect.Value, i int) (res Result) {
	st.Path.Push(PathElem{
		Index: i,
	})
	defer st.Path.Pop()
	res, _ = vc.CompareValues(st, v1.Index(i), v2.Index(i))
	return res
}
