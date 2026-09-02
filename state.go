package compare

import (
	"reflect"

	"github.com/pierrre/go-libs/syncutil"
)

// State represents the state of a running comparison.
//
// Functions must restore the original state when they return.
type State struct {
	// Depth is the current comparison depth.
	Depth int
	// Visited contains the visited value pairs, to prevent recursion.
	Visited map[VisitedEntry]struct{}
}

// Reset resets the state.
func (st *State) Reset() {
	st.Depth = 0
	clear(st.Visited)
}

var statePool = syncutil.Pool[*State]{
	New: func() *State {
		return &State{}
	},
}

// VisitedEntry represents a visited value pair.
type VisitedEntry struct {
	Type  reflect.Type
	Addr1 uintptr
	Addr2 uintptr
}
