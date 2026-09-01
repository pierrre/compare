package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/reflectutil"
)

// MapComparator is a [ValuesComparator] for map values.
//
// It should be created with [NewMapComparator].
type MapComparator struct {
	ValuesComparator
	// MaxDifferences is the maximum number of different items to report.
	// If reached, the comparison is stopped for the current map.
	// Setting it to 0 disables it.
	// Default: 10.
	MaxDifferences int
}

// NewMapComparator creates a new [MapComparator].
func NewMapComparator(vc ValuesComparator) *MapComparator {
	return &MapComparator{
		ValuesComparator: vc,
		MaxDifferences:   10,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *MapComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Map {
		return nil, false
	}
	res, handled = compareNilLenPointer(v1, v2)
	if handled {
		return res, true
	}
	return compareMap(vc.ValuesComparator, st, v1, v2, vc.MaxDifferences), true
}

// Supports implements [SupportChecker].
func (vc *MapComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.Map {
		res = vc
	}
	return res
}

func compareMap(vc ValuesComparator, st *State, v1, v2 reflect.Value, maxDifferences int) (res Result) {
	es1 := reflectutil.GetSortedMap(v1)
	es2 := reflectutil.GetSortedMap(v2)
	defer es1.Release()
	defer es2.Release()
	cmpFunc := reflectutil.GetCompareFunc(v1.Type().Key())
	diffCount := 0
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
			res = append(res, r...)
			diffCount++
			if diffCount >= maxDifferences && maxDifferences > 0 {
				break
			}
		}
	}
	return res
}

func compareMapKey(vc ValuesComparator, st *State, v1, v2, key reflect.Value) Result {
	res, _ := vc.CompareValues(st, v1, v2)
	if len(res) != 0 {
		res.AppendPathElem(PathElem{
			Key: new(getMapKeyInterface(key)),
		})
	}
	return res
}

func compareMapKeyMissing(key reflect.Value, inV1 bool) Result {
	r := Res("map key not defined", inV1, !inV1)
	r.AppendPathElem(PathElem{
		Key: new(getMapKeyInterface(key)),
	})
	return r
}

func getMapKeyInterface(v reflect.Value) any {
	i, _ := reflectutil.TryValueInterface(v)
	return i
}
