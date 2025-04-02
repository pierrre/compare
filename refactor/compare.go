package compare

import (
	"iter"
	"reflect"
	"strconv"
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

type ValuesSupportComparator interface {
	ValuesComparator
	Supports(typ reflect.Type) bool
}

type Difference struct {
	Message string
	V1      string
	V2      string
}

type IntValuesComparator struct{}

func (vc *IntValuesComparator) CompareValues(st *State, v1, v2 reflect.Value) (yieldOK bool, handled bool) {
	switch v1.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return false, false
	}
	i1 := v1.Int()
	i2 := v2.Int()
	if i1 == i2 {
		return true, true
	}
	return st.Yield(&Difference{
		Message: "different int",
		V1:      strconv.FormatInt(i1, 10),
		V2:      strconv.FormatInt(i2, 10),
	}), true
}
