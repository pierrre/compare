package compare

import (
	"github.com/pierrre/go-libs/syncutil"
)

// State represents the state of a comparison.
//
// Functions must restore the original state when they return.
type State struct {
	Depth   int
	Visited []Visited
	Result  Result
}

func (st *State) addDifferences(ds ...Difference) {
	st.Result = append(st.Result, ds...)
}

var statePool = syncutil.Pool[*State]{
	New: func() *State {
		return &State{}
	},
}

func (st *State) reset() {
	st.Depth = 0
	st.Visited = st.Visited[:0]
}

// Visited represents a visited pair of values.
type Visited struct {
	V1, V2 uintptr
}
