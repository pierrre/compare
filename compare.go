package compare

import (
	"reflect"
	"strconv"

	"github.com/pierrre/go-libs/reflectutil"
)

// CompareComparator is a [ValuesComparator] that compares values with a method that returns a comparison result (for example Compare(T) int).
//
// It should be created with [NewCompareComparator].
type CompareComparator struct {
	// CompareNames is the list of method names to look for, tried in order.
	// Default: ["Compare", "Cmp"].
	CompareNames []string
}

// NewCompareComparator creates a new [CompareComparator].
func NewCompareComparator() *CompareComparator {
	return &CompareComparator{
		CompareNames: []string{"Compare", "Cmp"},
	}
}

// CompareValues implements [ValuesComparator].
func (vc *CompareComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() == reflect.Interface {
		return nil, false
	}
	typ := v1.Type()
	if typ.NumMethod() == 0 {
		return nil, false
	}
	if isNilableKind(v1.Kind()) && (v1.IsNil() || v2.IsNil()) {
		return nil, false
	}
	if !v1.CanInterface() || !v2.CanInterface() {
		return nil, false
	}
	m, ok := vc.findMethod(typ)
	if !ok {
		return nil, false
	}
	cmp := m.Func.Call([]reflect.Value{v1, v2})[0].Int()
	if cmp == 0 {
		return nil, true
	}
	i1, _ := reflectutil.TryValueInterface(v1)
	i2, _ := reflectutil.TryValueInterface(v2)
	return Res("method "+m.Name+" returned "+strconv.FormatInt(cmp, 10), i1, i2), true
}

// Supports implements [SupportChecker].
func (vc *CompareComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() != reflect.Interface && typ.NumMethod() != 0 {
		_, ok := vc.findMethod(typ)
		if ok {
			res = vc
		}
	}
	return res
}

// findMethod returns the first method of the value's type whose name is in CompareNames and whose signature matches func(T, T) int.
func (vc *CompareComparator) findMethod(typ reflect.Type) (reflect.Method, bool) {
	methods := reflectutil.GetMethods(typ)
	for _, name := range vc.CompareNames {
		m, ok := methods.GetByName(name)
		if !ok {
			continue
		}
		mt := m.Type // receiver is the first argument: func(T, T) int
		if mt.NumIn() != 2 || mt.In(0) != typ || mt.In(1) != typ ||
			mt.NumOut() != 1 || mt.Out(0).Kind() != reflect.Int {
			continue
		}
		return m, true
	}
	return reflect.Method{}, false
}
