package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// EqualComparator is a [ValuesComparator] that compares values with a method that reports equality (for example Equal(T) bool).
//
// It should be created with [NewEqualComparator].
type EqualComparator struct {
	// Names is the list of method names to look for, tried in order.
	// Default: ["Equal", "Eq"].
	Names []string
}

// NewEqualComparator creates a new [EqualComparator].
func NewEqualComparator() *EqualComparator {
	return &EqualComparator{
		Names: []string{"Equal", "Eq"},
	}
}

// CompareValues implements [ValuesComparator].
func (vc *EqualComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
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
	eq := m.Func.Call([]reflect.Value{v1, v2})[0].Bool()
	if eq {
		return nil, true
	}
	i1, _ := reflectutil.TryValueInterface(v1)
	i2, _ := reflectutil.TryValueInterface(v2)
	return Res("method "+m.Name+" returned false", i1, i2), true
}

// Supports implements [SupportChecker].
func (vc *EqualComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() != reflect.Interface && typ.NumMethod() != 0 {
		_, ok := vc.findMethod(typ)
		if ok {
			res = vc
		}
	}
	return res
}

// findMethod returns the first method of the value's type whose name is in Names and whose signature matches func(T, T) bool.
func (vc *EqualComparator) findMethod(typ reflect.Type) (reflect.Method, bool) {
	methods := reflectutil.GetMethods(typ)
	for _, name := range vc.Names {
		m, ok := methods.GetByName(name)
		if !ok {
			continue
		}
		mt := m.Type // receiver is the first argument: func(T, T) bool
		if mt.NumIn() != 2 || mt.In(0) != typ || mt.In(1) != typ ||
			mt.NumOut() != 1 || mt.Out(0).Kind() != reflect.Bool {
			continue
		}
		return m, true
	}
	return reflect.Method{}, false
}
