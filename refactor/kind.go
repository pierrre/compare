package compare

import "reflect"

const kindsCount = reflect.UnsafePointer + 1

type KindValuesComparator struct {
	BaseInt   *IntValuesComparator
	BaseUint  *UintValuesComparator
	BaseFloat *FloatValuesComparator

	ValuesComparator [kindsCount]ValuesComparator
}

func New(vc ValuesComparator) *KindValuesComparator {
	kindVC := &KindValuesComparator{
		BaseInt:   NewIntValuesComparator(),
		BaseUint:  NewUintValuesComparator(),
		BaseFloat: NewFloatValuesComparator(),
	}
	kindVC.ValuesComparator = [kindsCount]ValuesComparator{
		reflect.Int:     ValuesComparatorFunc(kindVC.compareInt),
		reflect.Int8:    ValuesComparatorFunc(kindVC.compareInt),
		reflect.Int16:   ValuesComparatorFunc(kindVC.compareInt),
		reflect.Int32:   ValuesComparatorFunc(kindVC.compareInt),
		reflect.Int64:   ValuesComparatorFunc(kindVC.compareInt),
		reflect.Uint:    ValuesComparatorFunc(kindVC.compareUint),
		reflect.Uint8:   ValuesComparatorFunc(kindVC.compareUint),
		reflect.Uint16:  ValuesComparatorFunc(kindVC.compareUint),
		reflect.Uint32:  ValuesComparatorFunc(kindVC.compareUint),
		reflect.Uint64:  ValuesComparatorFunc(kindVC.compareUint),
		reflect.Float32: ValuesComparatorFunc(kindVC.compareFloat),
		reflect.Float64: ValuesComparatorFunc(kindVC.compareFloat),
	}
	_ = vc // TODO
	return kindVC
}

func (vc *KindValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) (handled bool) {
	return vc.ValuesComparator[v1.Kind()].CompareValues(st, v1, v2)
}

func (vc *KindValuesComparator) compareInt(st *State, v1, v2 reflect.Value) (handled bool) {
	return vc.BaseInt.CompareValues(st, v1, v2)
}

func (vc *KindValuesComparator) compareUint(st *State, v1, v2 reflect.Value) (handled bool) {
	return vc.BaseUint.CompareValues(st, v1, v2)
}

func (vc *KindValuesComparator) compareFloat(st *State, v1, v2 reflect.Value) (handled bool) {
	return vc.BaseFloat.CompareValues(st, v1, v2)
}
