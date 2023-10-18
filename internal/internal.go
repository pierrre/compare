package internal

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
)

// SortMapsKeys sorts maps keys.
func SortMapsKeys(typ reflect.Type, s []reflect.Value) {
	sortCmp := getMapsKeysSortCmp(typ)
	slices.SortFunc(s, sortCmp)
}

func getMapsKeysSortCmp(typ reflect.Type) func(a, b reflect.Value) int {
	switch typ.Kind() { //nolint:exhaustive // Optimized for common kinds, the default case is less optimized.
	case reflect.Bool:
		return func(a, b reflect.Value) int {
			ab := a.Bool()
			bb := b.Bool()
			if ab == bb {
				return 0
			}
			if !ab {
				return -1
			}
			return 1
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.Int(), b.Int())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.Uint(), b.Uint())
		}
	case reflect.Float32, reflect.Float64:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.Float(), b.Float())
		}
	case reflect.String:
		return func(a, b reflect.Value) int {
			return cmp.Compare(a.String(), b.String())
		}
	default:
		return func(a, b reflect.Value) int {
			return cmp.Compare(fmt.Sprint(a), fmt.Sprint(b))
		}
	}
}
