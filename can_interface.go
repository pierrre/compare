package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// CanInterfaceComparator is a [ValuesComparator] that attempts to convert the [reflect.Value]s so they can be used with [reflect.Value.Interface].
//
// It should be created with [NewCanInterfaceComparator].
type CanInterfaceComparator struct {
	ValuesComparator
}

// NewCanInterfaceComparator creates a new [CanInterfaceComparator].
func NewCanInterfaceComparator(vc ValuesComparator) *CanInterfaceComparator {
	return &CanInterfaceComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *CanInterfaceComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	v1 = vc.convertValue(v1)
	v2 = vc.convertValue(v2)
	return vc.ValuesComparator.CompareValues(st, v1, v2)
}

func (vc *CanInterfaceComparator) convertValue(v reflect.Value) reflect.Value {
	if v.CanInterface() {
		return v
	}
	if v.Kind() == reflect.Pointer {
		return v
	}
	v, _ = reflectutil.ConvertValueCanInterface(v)
	return v
}
