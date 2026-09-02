package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/syncutil"
)

// SupportChecker checks if a [reflect.Type] is supported.
// If the [reflect.Type] is supported, it returns a non-nil [ValuesComparator].
type SupportChecker interface {
	Supports(typ reflect.Type) ValuesComparator
}

// SupportCheckerFunc is a [SupportChecker] function.
type SupportCheckerFunc func(typ reflect.Type) ValuesComparator

// Supports implements [SupportChecker].
func (f SupportCheckerFunc) Supports(typ reflect.Type) ValuesComparator {
	return f(typ)
}

// SupportCheckerValuesComparator implements [ValuesComparator] and [SupportChecker].
type SupportCheckerValuesComparator struct {
	ValuesComparator
	SupportChecker
}

// SupportComparator is a [ValuesComparator] that selects a [ValuesComparator] based on the [reflect.Type] of the values.
// It selects the first [SupportChecker] that supports the [reflect.Type].
//
// It should be created with [NewSupportComparator].
type SupportComparator struct {
	cache    syncutil.Map[reflect.Type, valuesComparatorContainer]
	Checkers []SupportChecker
}

// NewSupportComparator creates a new [SupportComparator].
func NewSupportComparator() *SupportComparator {
	return &SupportComparator{}
}

// CompareValues implements [ValuesComparator].
func (vc *SupportComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	if len(vc.Checkers) == 0 {
		return nil, false
	}
	typ := v1.Type()
	vcc, ok := vc.cache.Load(typ)
	if !ok {
		vcc = valuesComparatorContainer{}
		for _, c := range vc.Checkers {
			vc2 := c.Supports(typ)
			if vc2 != nil {
				vcc.vc = vc2
				break
			}
		}
		vc.cache.Store(typ, vcc)
	}
	vc2 := vcc.vc
	if vc2 == nil {
		return nil, false
	}
	return vc2.CompareValues(st, v1, v2)
}

type valuesComparatorContainer struct {
	vc ValuesComparator
}

func supportsValuesComparator(typ reflect.Type, vc ValuesComparator) ValuesComparator {
	var res ValuesComparator
	c, ok := vc.(SupportChecker)
	if ok {
		res = c.Supports(typ)
	}
	return res
}

func callSupportsCheckerPointer[P interface {
	*T
	SupportChecker
}, T any](p P, typ reflect.Type) ValuesComparator {
	if p != nil {
		return p.Supports(typ)
	}
	return nil
}
