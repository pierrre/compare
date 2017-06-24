package compare

import (
	"fmt"
	"reflect"
)

// Func represents a comparison function.
// It is guaranteed that both values: are valid, of the same type, and can be converted to interface{}.
// If the returned value "stop" is true, the comparison will stop.
type Func func(v1, v2 reflect.Value) (r Result, stop bool)

var funcs []Func

// RegisterFunc registers a Func.
// It allows to handle manually the comparison for certain values.
func RegisterFunc(f Func) {
	funcs = append(funcs, f)
}

func compareFuncs(v1, v2 reflect.Value) (Result, bool) {
	if !v1.CanInterface() || !v2.CanInterface() {
		return nil, false
	}
	for _, f := range funcs {
		if r, stop := f(v1, v2); stop {
			return r, true
		}
	}
	return nil, false
}

func init() {
	RegisterFunc(compareMethodEqual)
	RegisterFunc(compareMethodCmp)
	RegisterFunc(compareValue)
}

var methodEqualNames []string

// RegisterMethodEqual registers an equal method.
// This methods must be callable as "v1.METHOD(v2) bool".
func RegisterMethodEqual(name string) {
	methodEqualNames = append(methodEqualNames, name)
}

func init() {
	RegisterMethodEqual("Equal")
	RegisterMethodEqual("Eq")
}

func compareMethodEqual(v1, v2 reflect.Value) (Result, bool) {
	for _, name := range methodEqualNames {
		r, stop := compareMethodEqualName(v1, v2, name)
		if stop {
			return r, true
		}
	}
	return nil, false
}

func compareMethodEqualName(v1, v2 reflect.Value, name string) (Result, bool) {
	m := v1.MethodByName(name)
	if !m.IsValid() {
		return nil, false
	}
	t := m.Type()
	if t.NumIn() != 1 || t.In(0) != v2.Type() || t.NumOut() != 1 || t.Out(0) != reflect.TypeOf(true) {
		return nil, false

	}
	if m.Call([]reflect.Value{v2})[0].Interface().(bool) {
		return nil, true
	}
	return Result{Difference{
		Message: fmt.Sprintf(msgMethodNotEqual, name),
		V1:      v1.Interface(),
		V2:      v2.Interface(),
	}}, true
}

func compareMethodCmp(v1, v2 reflect.Value) (Result, bool) {
	m := v1.MethodByName("Cmp")
	if !m.IsValid() {
		return nil, false
	}
	t := m.Type()
	if t.NumIn() != 1 || t.In(0) != v2.Type() || t.NumOut() != 1 || t.Out(0) != reflect.TypeOf(int(1)) {
		return nil, false
	}
	c := m.Call([]reflect.Value{v2})[0].Interface().(int)
	if c == 0 {
		return nil, true
	}
	return Result{Difference{
		Message: fmt.Sprintf(msgMethodCmpNotEqual, c),
		V1:      v1.Interface(),
		V2:      v2.Interface(),
	}}, true
}

var typeReflectValue = reflect.TypeOf(reflect.Value{})

func compareValue(v1, v2 reflect.Value) (Result, bool) {
	if v1.Type() != typeReflectValue {
		return nil, false
	}
	v1 = v1.Interface().(reflect.Value)
	v2 = v2.Interface().(reflect.Value)
	return compare(v1, v2), true
}
