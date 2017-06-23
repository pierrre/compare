package compare

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"sync"
)

func getMapsKeys(v1, v2 reflect.Value) []reflect.Value {
	ks := v1.MapKeys()
	for _, k2 := range v2.MapKeys() {
		found := false
		for _, k := range ks {
			if len(compare(k2, k)) == 0 {
				found = true
				break
			}
		}
		if !found {
			ks = append(ks, k2)
		}
	}
	sortValues(ks, v1.Type().Key())
	return ks
}

func sortValues(s []reflect.Value, t reflect.Type) {
	sort.Slice(s, newSortLess(s, t))
}

func newSortLess(s []reflect.Value, t reflect.Type) func(i, j int) bool {
	switch t.Kind() {
	case reflect.Bool:
		return newSortLessBool(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return newSortLessInt(s)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return newSortLessUint(s)
	case reflect.Float32, reflect.Float64:
		return newSortLessFloat(s)
	case reflect.Complex64, reflect.Complex128:
		return newSortLessComplex(s)
	case reflect.String:
		return newSortLessString(s)
	default:
		return newSortLessGeneric(s)
	}
}

func newSortLessBool(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return !s[i].Bool() && s[j].Bool()
	}
}

func newSortLessInt(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].Int() < s[j].Int()
	}
}

func newSortLessUint(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].Uint() < s[j].Uint()
	}
}

func newSortLessFloat(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].Float() < s[j].Float()
	}
}

func newSortLessComplex(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		ci := s[i].Complex()
		cj := s[j].Complex()
		if real(ci) < real(cj) {
			return true
		}
		if real(ci) > real(cj) {
			return false
		}
		return imag(ci) < imag(cj)
	}
}

func newSortLessString(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return s[i].String() < s[j].String()
	}
}

func newSortLessGeneric(s []reflect.Value) func(i, j int) bool {
	return func(i, j int) bool {
		return fmt.Sprint(s[i]) < fmt.Sprint(s[j])
	}
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
