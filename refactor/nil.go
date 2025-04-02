package compare

import (
	"reflect"
)

func compareNil(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	n1 := v1.IsNil()
	n2 := v2.IsNil()
	if n1 && n2 {
		return nil, true
	}
	if n1 != n2 {
		return st.Result("nil mismatch", n1, n2)
	}
	return nil, false
}

func compareNilPointerLen(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	res, handled = compareNil(st, v1, v2)
	if handled {
		return res, true
	}
	if v1.Pointer() == v2.Pointer() {
		return nil, true
	}
	len1 := v1.Len()
	len2 := v2.Len()
	if len1 != len2 {
		return st.Result("length not equal", len1, len2)
	}
	return nil, false
}
