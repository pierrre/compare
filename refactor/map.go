package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// MapComparator is a [ValuesComparator] for map values.
type MapComparator struct {
	ValuesComparator
}

// NewMapComparator creates a new [MapComparator].
func NewMapComparator(vc ValuesComparator) *MapComparator {
	return &MapComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *MapComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Map {
		return nil, false
	}
	res, handled = compareNilPointerLen(st, v1, v2)
	if handled {
		return res, true
	}
	return compareMap(vc.ValuesComparator, st, v1, v2), true
}

func compareMap(vc ValuesComparator, st *State, v1, v2 reflect.Value) (res Result) {
	es1 := reflectutil.GetSortedMap(v1)
	es2 := reflectutil.GetSortedMap(v2)
	defer es1.Release()
	defer es2.Release()
	cmpFunc := reflectutil.GetCompareFunc(v1.Type().Key())
	i1 := 0
	i2 := 0
	for i1 < len(es1) || i2 < len(es2) {
		var cm int
		switch {
		case i1 >= len(es1):
			cm = 1
		case i2 >= len(es2):
			cm = -1
		default:
			cm = cmpFunc(es1[i1].Key, es2[i2].Key)
		}
		var r Result
		switch {
		case cm < 0:
			r = compareMapKeyMissing(es1[i1].Key, true)
			i1++
		case cm > 0:
			r = compareMapKeyMissing(es2[i2].Key, false)
			i2++
		default:
			r = compareMapKey(vc, st, es1[i1].Value, es2[i2].Value, es1[i1].Key)
			i1++
			i2++
		}
		if len(r) > 0 {
			// TODO: max diff
			res = append(res, r...)
		}
	}
	return res
}

func compareMapKey(vc ValuesComparator, st *State, v1, v2, key reflect.Value) Result {
	res, _ := vc.CompareValues(st, v1, v2)
	if len(res) != 0 {
		res.AppendPathElem(PathElem{
			Key: new(any(copyReflectValue(key))),
		})
	}
	return res
}

func compareMapKeyMissing(key reflect.Value, inV1 bool) Result {
	r := Res("map key not defined", inV1, !inV1)
	r.AppendPathElem(PathElem{
		Key: new(any(copyReflectValue(key))),
	})
	return r
}

func copyReflectValue(v reflect.Value) reflect.Value {
	c := reflect.New(v.Type()).Elem()
	c.Set(v)
	return c
}
