package compare

import (
	"reflect"
	"sync/atomic"
)

// Compare compares 2 values with [DefaultComparator].
func Compare(iv1, iv2 any) Result {
	return DefaultComparator.Load().Compare(iv1, iv2)
}

// DefaultComparator is the default [Comparator].
//
// It uses [DefaultValuesComparator].
var DefaultComparator atomic.Pointer[Comparator]

func init() {
	DefaultComparator.Store(NewComparator(DefaultValuesComparator.Load()))
}

// Comparator is a value comparator.
//
// It should be created with [NewComparator].
type Comparator struct {
	ValuesComparator ValuesComparator
}

// NewComparator creates a new [Comparator].
func NewComparator(vc ValuesComparator) *Comparator {
	return &Comparator{
		ValuesComparator: vc,
	}
}

// Compare compares 2 values.
func (c *Comparator) Compare(iv1, iv2 any) Result {
	v1 := reflect.ValueOf(iv1)
	v2 := reflect.ValueOf(iv2)
	vl1 := v1.IsValid()
	vl2 := v2.IsValid()
	if !vl1 && !vl2 {
		return nil
	}
	if vl1 != vl2 {
		return Res("nil mismatch", !vl1, !vl2)
	}
	st := statePool.Get()
	defer statePool.Put(st)
	st.Reset()
	res, handled := compareType(v1, v2)
	if handled {
		return res
	}
	res, handled = c.ValuesComparator.CompareValues(st, v1, v2)
	if !handled {
		return Res("not handled", v1.Type(), v2.Type())
	}
	return res
}
