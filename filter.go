package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// FilterComparator is a [ValuesComparator] that calls the [ValuesComparator] if the filter returns true.
//
// It should be created with [NewFilterComparator].
type FilterComparator[VC ValuesComparator] struct {
	ValuesComparator VC
	// Filter filters types.
	// The value is handled if it returns true or if it is nil.
	Filter func(typ reflect.Type) bool
}

// NewFilterComparator creates a new [FilterComparator].
func NewFilterComparator[VC ValuesComparator](vc VC, f func(typ reflect.Type) bool) *FilterComparator[VC] {
	return &FilterComparator[VC]{
		ValuesComparator: vc,
		Filter:           f,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *FilterComparator[VC]) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if vc.Filter == nil || vc.Filter(v1.Type()) {
		return vc.ValuesComparator.CompareValues(st, v1, v2)
	}
	return nil, false
}

// Supports implements [SupportChecker].
func (vc *FilterComparator[VC]) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if vc.Filter == nil || vc.Filter(typ) {
		res = supportsValuesComparator(typ, vc.ValuesComparator)
	}
	return res
}

// FilterTypes returns a new filter function that returns true if the type is in the given list or if it implements any of the given interface types.
func FilterTypes(typs ...reflect.Type) func(typ reflect.Type) bool {
	set := make(map[reflect.Type]struct{}, len(typs))
	var ics []*reflectutil.ImplementsCache
	for _, typ := range typs {
		if _, ok := set[typ]; !ok {
			set[typ] = struct{}{}
			if typ.Kind() == reflect.Interface {
				ics = append(ics, reflectutil.NewImplementsCache(typ))
			}
		}
	}
	return func(typ reflect.Type) bool {
		_, ok := set[typ]
		if ok {
			return true
		}
		for _, ic := range ics {
			if ic.ImplementedBy(typ) {
				return true
			}
		}
		return false
	}
}
