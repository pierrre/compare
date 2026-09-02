package compare_test

import (
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/compare"
)

func TestStateReset(t *testing.T) {
	st := &State{
		Depth: 42,
		Visited: map[VisitedEntry]struct{}{
			{}: {},
		},
	}
	st.Reset()
	assert.Zero(t, st.Depth)
	assert.MapLen(t, st.Visited, 0)
	st.Reset()
	assert.Zero(t, st.Depth)
	assert.MapLen(t, st.Visited, 0)
}
