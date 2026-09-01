package compare

import (
	"math"
	"reflect"
)

// FloatComparator is a [ValuesComparator] for float values.
//
// It should be created with [NewFloatComparator].
type FloatComparator struct {
	// NaNEqual controls whether two NaN values compare equal.
	// Default: true.
	NaNEqual bool
	// SignedZeroEqual controls whether +0.0 and -0.0 compare equal.
	// Default: true.
	SignedZeroEqual bool
}

// NewFloatComparator creates a new [FloatComparator].
func NewFloatComparator() *FloatComparator {
	return &FloatComparator{
		NaNEqual:        true,
		SignedZeroEqual: true,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *FloatComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	switch v1.Kind() { //nolint:exhaustive // Only supports float.
	case reflect.Float32, reflect.Float64:
	default:
		return nil, false
	}
	f1 := v1.Float()
	f2 := v2.Float()
	if compareFloatEqual(vc.NaNEqual, vc.SignedZeroEqual, f1, f2) {
		return nil, true
	}
	return Res("float not equal", f1, f2), true
}

// Supports implements [SupportChecker].
func (vc *FloatComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	switch typ.Kind() { //nolint:exhaustive // Only supports float.
	case reflect.Float32, reflect.Float64:
		res = vc
	}
	return res
}

func compareFloatEqual(nanEqual, signedZeroEqual bool, f1, f2 float64) bool {
	if f1 == f2 {
		return signedZeroEqual || f1 != 0 || math.Signbit(f1) == math.Signbit(f2)
	}
	return nanEqual && math.IsNaN(f1) && math.IsNaN(f2)
}
