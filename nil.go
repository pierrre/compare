package compare

import (
	"reflect"
)

func compareNil(v1, v2 reflect.Value) (res Result, handled bool) {
	n1 := v1.IsNil()
	n2 := v2.IsNil()
	if n1 && n2 {
		return nil, true
	}
	if n1 != n2 {
		return Res("nil mismatch", n1, n2), true
	}
	return nil, false
}

func compareNilLenPointer(v1, v2 reflect.Value) (res Result, handled bool) {
	res, handled = compareNil(v1, v2)
	if handled {
		return res, true
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return Res("length not equal", len1, len2), true
	}
	if v1.Pointer() == v2.Pointer() {
		return nil, true
	}
	return nil, false
}

func isNilableKind(k reflect.Kind) bool {
	switch k { //nolint:exhaustive // Only nil-able kinds are handled.
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	}
	return false
}
