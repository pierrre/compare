package compare

import (
	"reflect"
	"runtime"
)

// FuncComparator is a [ValuesComparator] for func values.
//
// It should be created with [NewFuncComparator].
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
	res, handled = compareNil(v1, v2)
	if handled {
		return res, true
	}
	p1 := v1.Pointer()
	p2 := v2.Pointer()
	if p1 == p2 {
		return nil, true
	}
	var fn1 string
	f1 := runtime.FuncForPC(p1)
	if f1 != nil {
		fn1 = f1.Name()
	}
	var fn2 string
	f2 := runtime.FuncForPC(p2)
	if f2 != nil {
		fn2 = f2.Name()
	}
	return Res("function pointer not equal", fn1, fn2), true
}

// Supports implements [SupportChecker].
func (vc *FuncComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.Func {
		res = vc
	}
	return res
}
