package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// StructComparator is a [ValuesComparator] for struct values.
//
// It should be created with [NewStructComparator].
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
			res = append(res, r...)
		}
	}
	return res, true
}

// Supports implements [SupportChecker].
func (vc *StructComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.Struct {
		res = vc
	}
	return res
}

func (vc *StructComparator) compareField(st *State, v1, v2 reflect.Value, i int, sf reflect.StructField) Result { //nolint:gocritic // StructField is large.
	res, _ := vc.ValuesComparator.CompareValues(st, v1.Field(i), v2.Field(i))
	if len(res) != 0 {
		res.AppendPathElem(PathElem{
			Name: new(sf.Name),
		})
	}
	return res
}
