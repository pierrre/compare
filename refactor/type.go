package compare

import (
	"reflect"
)

func compareType(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	t1 := v1.Type()
	t2 := v1.Type()
	if t1 != t2 {
		return st.Result("type not equal", t1, t2)
	}
	return nil, false
}
