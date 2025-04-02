package compare

import (
	"iter"
	"reflect"
)

type Comparator struct {
	ValuesComparator ValuesComparator
}

func (c *Comparator) Compare(iv1, iv2 any) iter.Seq[*Difference] {
	return func(yield func(*Difference) bool) {
		v1 := reflect.ValueOf(iv1)
		v2 := reflect.ValueOf(iv2)
		st := &State{
			Yield: yield,
		}
		c.ValuesComparator.CompareValues(st, v1, v2)
	}
}

type State struct {
	Yield func(*Difference) bool
}

type ValuesComparator interface {
	CompareValues(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool)
}

type ValuesComparatorFunc func(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool)

func (f ValuesComparatorFunc) CompareValues(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool) {
	return f(st, v1, v2)
}

type Difference struct {
	Message string
	V1      string
	V2      string
}
