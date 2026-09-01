package compare

import (
	"reflect"
)

// ArrayComparator is a [ValuesComparator] for array values.
//
// It should be created with [NewArrayComparator].
type ArrayComparator struct {
	ValuesComparator
	// MaxDifferences is the maximum number of different items to report.
	// If reached, the comparison is stopped for the current array.
	// Setting it to 0 disables it.
	// Default: 10.
	MaxDifferences int
}

// NewArrayComparator creates a new [ArrayComparator].
func NewArrayComparator(vc ValuesComparator) *ArrayComparator {
	return &ArrayComparator{
		ValuesComparator: vc,
		MaxDifferences:   10,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *ArrayComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if v1.Kind() != reflect.Array {
		return nil, false
	}
	return compareArray(vc.ValuesComparator, st, v1, v2, vc.MaxDifferences), true
}

// Supports implements [SupportChecker].
func (vc *ArrayComparator) Supports(typ reflect.Type) ValuesComparator {
	var res ValuesComparator
	if typ.Kind() == reflect.Array {
		res = vc
	}
	return res
}

func compareArray(vc ValuesComparator, st *State, v1, v2 reflect.Value, maxDifferences int) (res Result) {
	diffCount := 0
	for i := range v1.Len() {
		r := compareArrayIndex(vc, st, v1, v2, i)
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

func compareArrayIndex(vc ValuesComparator, st *State, v1, v2 reflect.Value, i int) (res Result) {
	res, _ = vc.CompareValues(st, v1.Index(i), v2.Index(i))
	if len(res) != 0 {
		res.AppendPathElem(PathElem{
			Index: new(i),
		})
	}
	return res
}
