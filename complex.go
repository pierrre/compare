package compare

import (
	"reflect"
)

// ComplexComparator is a [ValuesComparator] for complex values.
//
// It should be created with [NewComplexComparator].
type ComplexComparator struct {
	// NaNEqual controls whether the real and imaginary parts with NaN values compare equal.
	// Default: true
	NaNEqual bool
	// SignedZeroEqual controls whether real and imaginary parts with +0.0 and -0.0 compare equal.
	// Default: true
	SignedZeroEqual bool
}

// NewComplexComparator creates a new [ComplexComparator].
func NewComplexComparator() *ComplexComparator {
	return &ComplexComparator{
		NaNEqual:        true,
		SignedZeroEqual: true,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *ComplexComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	switch v1.Kind() { //nolint:exhaustive // Only supports complex.
	case reflect.Complex64, reflect.Complex128:
	default:
		return nil, false
	}
	c1 := v1.Complex()
	c2 := v2.Complex()
	if compareFloatEqual(vc.NaNEqual, vc.SignedZeroEqual, real(c1), real(c2)) &&
		compareFloatEqual(vc.NaNEqual, vc.SignedZeroEqual, imag(c1), imag(c2)) {
		return nil, true
	}
	return Res("complex not equal", c1, c2), true
}

// Supports implements [SupportChecker].
func (vc *ComplexComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	switch typ.Kind() { //nolint:exhaustive // Only supports complex.
	case reflect.Complex64, reflect.Complex128:
		res = vc
	}
	return res
}
