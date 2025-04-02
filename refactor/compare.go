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
			yield:   yield,
			YieldOK: true,
		}
		c.ValuesComparator.CompareValues(st, v1, v2)
	}
}

type State struct {
	yield   func(*Difference) bool
	YieldOK bool
}

func (st *State) Yield(diff *Difference) {
	if st.YieldOK {
		st.YieldOK = st.yield(diff)
	}
}

type ValuesComparator interface {
	CompareValues(st *State, v1, v2 reflect.Value) (handled bool)
}

type ValuesComparatorFunc func(st *State, v1, v2 reflect.Value) (handled bool)

func (f ValuesComparatorFunc) CompareValues(st *State, v1, v2 reflect.Value) (handled bool) {
	return f(st, v1, v2)
}

type Difference struct {
	Message string
	V1      string
	V2      string
}
