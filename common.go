package compare

import (
	"reflect"
	"sync/atomic"
)

// DefaultValuesComparator is the default [CommonComparator].
var DefaultValuesComparator atomic.Pointer[CommonComparator]

func init() {
	DefaultValuesComparator.Store(NewCommonComparator())
}

// CommonComparator is a [ValuesComparator] with common [ValuesComparator]s.
//
// Any [ValuesComparator] can be configured, and can be set to nil in order to disable it.
//
// It should be created with [NewCommonComparator].
type CommonComparator struct {
	Recursion         *RecursionComparator
	MaxDepth          *MaxDepthComparator
	CanInterface      *CanInterfaceComparator
	ByType            ByTypeComparators
	ValuesComparators ValuesComparators
	Support           *SupportComparator
	ReflectValue      *ReflectValueComparator
	Bytes             *BytesComparator
	Equal             *FilterComparator[*EqualComparator]
	Compare           *FilterComparator[*CompareComparator]
	Kind              *KindComparator
}

// NewCommonComparator creates a new [CommonComparator] initialized with default values.
func NewCommonComparator() *CommonComparator {
	vc := &CommonComparator{}
	vc.Recursion = NewRecursionComparator(nil)
	vc.MaxDepth = NewMaxDepthComparator(nil)
	vc.CanInterface = NewCanInterfaceComparator(nil)
	vc.ByType = NewByTypeComparators()
	vc.Support = NewSupportComparator()
	vc.Support.Checkers = []SupportChecker{vc}
	vc.ReflectValue = NewReflectValueComparator(vc)
	vc.Bytes = NewBytesComparator()
	vc.Equal = NewFilterComparator(NewEqualComparator(), nil)
	vc.Compare = NewFilterComparator(NewCompareComparator(), nil)
	vc.Kind = NewKindComparator(vc)
	return vc
}

// CompareValues implements [ValuesComparator].
func (vc *CommonComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if vc.Recursion != nil {
		recursionEntry, visitedAdded, recursionDetected := vc.Recursion.checkRecursion(st, v1, v2)
		if recursionDetected {
			return nil, true
		}
		if visitedAdded {
			defer vc.Recursion.postRecursion(st, recursionEntry)
		}
	}
	if vc.MaxDepth != nil {
		res, maxReached := vc.MaxDepth.checkMaxDepth(st, v1, v2)
		defer vc.MaxDepth.postMaxDepth(st)
		if maxReached {
			return res, true
		}
	}
	if vc.CanInterface != nil {
		v1 = vc.CanInterface.convertValue(v1)
		v2 = vc.CanInterface.convertValue(v2)
	}
	return vc.compareValues(st, v1, v2)
}

//nolint:gocyclo // Sequential optional comparator pipeline.
func (vc *CommonComparator) compareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if len(vc.ByType) != 0 {
		res, handled = vc.ByType.CompareValues(st, v1, v2)
		if handled {
			return res, true
		}
	}
	if len(vc.ValuesComparators) != 0 {
		res, handled = vc.ValuesComparators.CompareValues(st, v1, v2)
		if handled {
			return res, true
		}
	}
	if res, handled = callValuesComparatorPointer(vc.Support, st, v1, v2); handled {
		return res, true
	}
	if res, handled = callValuesComparatorPointer(vc.ReflectValue, st, v1, v2); handled {
		return res, true
	}
	if res, handled = callValuesComparatorPointer(vc.Bytes, st, v1, v2); handled {
		return res, true
	}
	if res, handled = callValuesComparatorPointer(vc.Equal, st, v1, v2); handled {
		return res, true
	}
	if res, handled = callValuesComparatorPointer(vc.Compare, st, v1, v2); handled {
		return res, true
	}
	if res, handled = callValuesComparatorPointer(vc.Kind, st, v1, v2); handled {
		return res, true
	}
	return nil, false
}

func callValuesComparatorPointer[VC interface {
	*T
	ValuesComparator
}, T any](p VC, st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if p != nil {
		return p.CompareValues(st, v1, v2)
	}
	return nil, false
}

// Supports implements [SupportChecker].
func (vc *CommonComparator) Supports(typ reflect.Type) ValuesComparator {
	if w := callSupportsCheckerPointer(vc.ReflectValue, typ); w != nil {
		return w
	}
	if w := callSupportsCheckerPointer(vc.Equal, typ); w != nil {
		return w
	}
	if w := callSupportsCheckerPointer(vc.Compare, typ); w != nil {
		return w
	}
	if w := callSupportsCheckerPointer(vc.Bytes, typ); w != nil {
		return w
	}
	if w := callSupportsCheckerPointer(vc.Kind, typ); w != nil {
		return w
	}
	return nil
}
