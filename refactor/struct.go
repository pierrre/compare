package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// StructComparator is a [ValuesComparator] for struct values.
type StructComparator struct {
	ValuesComparator
}

// NewStructComparator creates a new [NewStructComparator].
func NewStructComparator(vc ValuesComparator) *StructComparator {
	return &StructComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *StructComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Struct {
		return nil, false
	}
	sfs := reflectutil.GetStructFields(v1.Type())
	for i, sf := range sfs.Range {
		r := vc.compareField(st, v1, v2, i, sf)
		if len(r) > 0 {
			// TODO: max diff
			res = append(res, r...)
		}
	}
	return res, true
}

func (vc *StructComparator) compareField(st *State, v1, v2 reflect.Value, i int, sf reflect.StructField) Result { //nolint:gocritic // StructField is large.
	st.Path.Push(PathElem{
		Name: sf.Name,
	})
	defer st.Path.Pop()
	res, _ := vc.ValuesComparator.CompareValues(st, v1.Field(i), v2.Field(i))
	return res
}
