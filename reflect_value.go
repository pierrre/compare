package compare

import (
	"reflect"
)

var reflectValueType = reflect.TypeFor[reflect.Value]()

// ReflectValueComparator is a [ValuesComparator] for [reflect.Value] values.
//
// It unwraps the values and compares the wrapped values.
//
// It should be created with [NewReflectValueComparator].
type ReflectValueComparator struct {
	ValuesComparator
}

// NewReflectValueComparator creates a new [ReflectValueComparator].
func NewReflectValueComparator(vc ValuesComparator) *ReflectValueComparator {
	return &ReflectValueComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *ReflectValueComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Struct || v1.Type() != reflectValueType {
		return nil, false
	}
	if !v1.CanInterface() || !v2.CanInterface() {
		return Res("not comparable", v1.Type(), v2.Type()), true
	}
	rv1, _ := reflect.TypeAssert[reflect.Value](v1)
	rv2, _ := reflect.TypeAssert[reflect.Value](v2)
	vl1 := rv1.IsValid()
	vl2 := rv2.IsValid()
	if !vl1 && !vl2 {
		return nil, true
	}
	if vl1 != vl2 {
		return Res("nil mismatch", vl1, vl2), true
	}
	res, handled = compareType(rv1, rv2)
	if handled {
		return res, true
	}
	res, _ = vc.ValuesComparator.CompareValues(st, rv1, rv2)
	return res, true
}

// Supports implements [SupportChecker].
func (vc *ReflectValueComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ == reflectValueType {
		res = vc
	}
	return res
}
