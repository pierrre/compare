package compare

import (
	"reflect"
)

// RecursionComparator is a [ValuesComparator] that prevents infinite recursion.
//
// It should be created with [NewRecursionComparator].
type RecursionComparator struct {
	ValuesComparator
}

// NewRecursionComparator creates a new [RecursionComparator].
func NewRecursionComparator(vc ValuesComparator) *RecursionComparator {
	return &RecursionComparator{
		ValuesComparator: vc,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *RecursionComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	e, visitedAdded, recursionDetected := vc.checkRecursion(st, v1, v2)
	if recursionDetected {
		return nil, true
	}
	if visitedAdded {
		defer vc.postRecursion(st, e)
	}
	return vc.ValuesComparator.CompareValues(st, v1, v2)
}

func (vc *RecursionComparator) checkRecursion(st *State, v1, v2 reflect.Value) (e VisitedEntry, visitedAdded bool, recursionDetected bool) {
	switch v1.Kind() { //nolint:exhaustive // Only pointer-like kinds can be recursive.
	case reflect.Pointer, reflect.Map, reflect.Slice:
	default:
		return VisitedEntry{}, false, false
	}
	if v1.IsNil() || v2.IsNil() {
		return VisitedEntry{}, false, false
	}
	e = VisitedEntry{
		Type:  v1.Type(),
		Addr1: uintptr(v1.UnsafePointer()),
		Addr2: uintptr(v2.UnsafePointer()),
	}
	_, ok := st.Visited[e]
	if ok {
		return VisitedEntry{}, false, true
	}
	if st.Visited == nil {
		st.Visited = make(map[VisitedEntry]struct{})
	}
	st.Visited[e] = struct{}{}
	return e, true, false
}

func (vc *RecursionComparator) postRecursion(st *State, e VisitedEntry) {
	delete(st.Visited, e)
}
