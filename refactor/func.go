package compare

import (
	"reflect"
)

// FuncComparator is a [ValuesComparator] for func values.
type FuncComparator struct{}

// NewFuncComparator creates a new [FuncComparator].
func NewFuncComparator() *FuncComparator {
	return &FuncComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *FuncComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Func {
		return nil, false
	}
	res, handled = compareNil(st, v1, v2)
	if handled {
		return res, true
	}
	p1 := v1.Pointer()
	p2 := v2.Pointer()
	if p1 == p2 {
		return nil, true
	}
	return st.Result("function pointer not equal", p1, p2) // TODO: show function name ?
}
