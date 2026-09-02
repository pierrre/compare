package compare

import (
	"fmt"
	"reflect"
)

// MaxDepthComparator is a [ValuesComparator] that limits the comparison depth.
//
// It should be created with [NewMaxDepthComparator].
type MaxDepthComparator struct {
	ValuesComparator
	// Max is the maximum depth.
	// Default: 0 (no limit).
	Max int
	// Report controls whether a difference is reported if the max depth is reached.
	// Default: true.
	Report bool
}

// NewMaxDepthComparator creates a new [MaxDepthComparator].
func NewMaxDepthComparator(vc ValuesComparator) *MaxDepthComparator {
	return &MaxDepthComparator{
		ValuesComparator: vc,
		Report:           true,
	}
}

// CompareValues implements [ValuesComparator].
func (vc *MaxDepthComparator) CompareValues(st *State, v1, v2 reflect.Value) (res Result, handled bool) {
	res, maxReached := vc.checkMaxDepth(st, v1, v2)
	defer vc.postMaxDepth(st)
	if maxReached {
		return res, true
	}
	return vc.ValuesComparator.CompareValues(st, v1, v2)
}

func (vc *MaxDepthComparator) checkMaxDepth(st *State, v1, v2 reflect.Value) (res Result, maxReached bool) {
	if vc.Max > 0 && st.Depth >= vc.Max {
		maxReached = true
		if vc.Report {
			res = Res(fmt.Sprintf("max depth reached: %d", vc.Max), v1.Type(), v2.Type())
		}
	}
	st.Depth++
	return res, maxReached
}

func (vc *MaxDepthComparator) postMaxDepth(st *State) {
	st.Depth--
}
