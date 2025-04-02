package compare

import (
	"reflect"
)

const kindsCount = reflect.UnsafePointer + 1

type KindValuesComparator struct {
	Bool      *BoolValuesComparator
	BaseInt   *IntValuesComparator
	BaseUint  *UintValuesComparator
	BaseFloat *FloatValuesComparator
	Complex   *ComplexValuesComparator
	Pointer   *PointerValuesComparator
	String    *StringValuesComparator

	ValuesComparator [kindsCount]ValuesComparator
}

func New(vc ValuesComparator) *KindValuesComparator {
	kindVC := &KindValuesComparator{
		Bool:      NewBoolValuesComparator(),
		BaseInt:   NewIntValuesComparator(),
		BaseUint:  NewUintValuesComparator(),
		BaseFloat: NewFloatValuesComparator(),
		Complex:   NewComplexValuesComparator(),
		Pointer:   NewPointerValuesComparator(vc),
		String:    NewStringValuesComparator(),
	}
	kindVC.ValuesComparator = [kindsCount]ValuesComparator{
		reflect.Bool:       kindVC.Bool,
		reflect.Int:        kindVC.BaseInt,
		reflect.Int8:       kindVC.BaseInt,
		reflect.Int16:      kindVC.BaseInt,
		reflect.Int32:      kindVC.BaseInt,
		reflect.Int64:      kindVC.BaseInt,
		reflect.Uint:       kindVC.BaseUint,
		reflect.Uint8:      kindVC.BaseUint,
		reflect.Uint16:     kindVC.BaseUint,
		reflect.Uint32:     kindVC.BaseUint,
		reflect.Uint64:     kindVC.BaseUint,
		reflect.Float32:    kindVC.BaseFloat,
		reflect.Float64:    kindVC.BaseFloat,
		reflect.Complex64:  kindVC.Complex,
		reflect.Complex128: kindVC.Complex,
		reflect.Pointer:    kindVC.Pointer,
		reflect.String:     kindVC.String,
	}
	return kindVC
}
