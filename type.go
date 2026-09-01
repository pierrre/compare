package compare

import (
	"reflect"
)

// ByTypeComparators selects a [ValuesComparator] by [reflect.Type].
//
// It should be created with [NewByTypeComparators].
type ByTypeComparators map[reflect.Type]ValuesComparator

// NewByTypeComparators creates a new [ByTypeComparators].
func NewByTypeComparators() ByTypeComparators {
	return make(ByTypeComparators)
}

// CompareValues implements [ValuesComparator].
func (vcs ByTypeComparators) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if len(vcs) == 0 {
		return nil, false
	}
	vc := vcs[v1.Type()]
	if vc == nil {
		return nil, false
	}
	return vc.CompareValues(st, v1, v2)
}

func compareType(v1, v2 reflect.Value) (res Result, handled bool) {
	t1 := v1.Type()
	t2 := v2.Type()
	if t1 != t2 {
		return Res("type not equal", t1, t2), true
	}
	return nil, false
}
