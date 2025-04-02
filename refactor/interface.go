package compare

import (
	"reflect"
)

// InterfaceComparator is a [ValuesComparator] for interface values.
type InterfaceComparator struct {
	ValuesComparator
}

// NewInterfaceComparator creates a new [InterfaceComparator].
func NewInterfaceComparator(vc ValuesComparator) *InterfaceComparator {
	return &InterfaceComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *InterfaceComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Interface {
		return nil, false
	}
	res, handled = compareNil(st, v1, v2)
	if handled {
		return res, true
	}
	v1 = v1.Elem()
	v2 = v2.Elem()
	res, handled = compareType(st, v1, v2)
	if handled {
		return res, true
	}
	return vc.ValuesComparator.CompareValues(st, v1, v2)
}
