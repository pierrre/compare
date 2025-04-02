package compare

import (
	"reflect"
)

const kindsCount = reflect.UnsafePointer + 1

// KindComparator is a [ValuesComparator] for different kinds of values.
type KindComparator struct {
	Bool          *BoolComparator
	Int           *IntComparator
	Uint          *UintComparator
	Float         *FloatComparator
	Complex       *ComplexComparator
	Array         *ArrayComparator
	Chan          *ChanComparator
	Func          *FuncComparator
	Interface     *InterfaceComparator
	Pointer       *PointerComparator
	Slice         *SliceComparator
	Map           *MapComparator
	String        *StringComparator
	Struct        *StructComparator
	UnsafePointer *UnsafePointerComparator

	ValuesComparators [kindsCount]ValuesComparator
}

// NewKindComparator creates a new KindValuesComparator.
func NewKindComparator(vc ValuesComparator) *KindComparator {
	kindVC := &KindComparator{
		Bool:          NewBoolComparator(),
		Int:           NewIntComparator(),
		Uint:          NewUintComparator(),
		Float:         NewFloatComparator(),
		Complex:       NewComplexComparator(),
		Array:         NewArrayComparator(vc),
		Chan:          NewChanComparator(),
		Func:          NewFuncComparator(),
		Interface:     NewInterfaceComparator(vc),
		Pointer:       NewPointerComparator(vc),
		Slice:         NewSliceComparator(vc),
		Map:           NewMapComparator(vc),
		String:        NewStringComparator(),
		Struct:        NewStructComparator(vc),
		UnsafePointer: NewUnsafePointerComparator(),
	}
	kindVC.ValuesComparators = [kindsCount]ValuesComparator{
		reflect.Bool:          kindVC.Bool,
		reflect.Int:           kindVC.Int,
		reflect.Int8:          kindVC.Int,
		reflect.Int16:         kindVC.Int,
		reflect.Int32:         kindVC.Int,
		reflect.Int64:         kindVC.Int,
		reflect.Uint:          kindVC.Uint,
		reflect.Uint8:         kindVC.Uint,
		reflect.Uint16:        kindVC.Uint,
		reflect.Uint32:        kindVC.Uint,
		reflect.Uint64:        kindVC.Uint,
		reflect.Float32:       kindVC.Float,
		reflect.Float64:       kindVC.Float,
		reflect.Complex64:     kindVC.Complex,
		reflect.Complex128:    kindVC.Complex,
		reflect.Array:         kindVC.Array,
		reflect.Chan:          kindVC.Chan,
		reflect.Func:          kindVC.Func,
		reflect.Interface:     kindVC.Interface,
		reflect.Pointer:       kindVC.Pointer,
		reflect.Slice:         kindVC.Slice,
		reflect.Map:           kindVC.Map,
		reflect.String:        kindVC.String,
		reflect.Struct:        kindVC.Struct,
		reflect.UnsafePointer: kindVC.UnsafePointer,
	}
	return kindVC
}

// CompareValues implements [ValuesComparator].
func (kvc *KindComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	return kvc.ValuesComparators[v1.Kind()].CompareValues(st, v1, v2)
}
