package compare_test

import (
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("MaxDepth", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   [][]int{{1, 2}, {3, 4}},
			V2:   [][]int{{1, 2}, {3, 4}},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.MaxDepth.Max = 10
			},
		},
		{
			Name: "Reached",
			V1:   [][]int{{1, 2}, {3, 4}},
			V2:   [][]int{{1, 2}, {3, 4}},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.MaxDepth.Max = 1
			},
		},
		{
			Name: "Standalone",
			V1:   1,
			V2:   1,
			ConfigureComparator: func(c *Comparator) {
				c.ValuesComparator = NewMaxDepthComparator(NewIntComparator())
			},
		},
		{
			Name: "StandaloneReentrant",
			V1:   [][]int{{1, 2}},
			V2:   [][]int{{1, 2}},
			ConfigureComparator: func(c *Comparator) {
				var m *MaxDepthComparator
				m = NewMaxDepthComparator(ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					return m.CompareValues(st, v1, v2)
				}))
				m.Max = 1
				c.ValuesComparator = m
			},
		},
		{
			Name: "StandaloneSiblings",
			V1:   [][][][]int{{{{1, 2}, {3, 4}}}, {{{5, 6}}}},
			V2:   [][][][]int{{{{1, 2}, {3, 4}}}, {{{5, 6}}}},
			ConfigureComparator: func(c *Comparator) {
				var m *MaxDepthComparator
				m = NewMaxDepthComparator(ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					if v1.Kind() != reflect.Slice {
						return nil, true
					}
					var res Result
					for i := range v1.Len() {
						r, _ := m.CompareValues(st, v1.Index(i), v2.Index(i))
						res = append(res, r...)
					}
					return res, true
				}))
				m.Max = 3
				c.ValuesComparator = m
			},
		},
	})
}

func TestMaxDepthRestoreState(t *testing.T) {
	var m *MaxDepthComparator
	m = NewMaxDepthComparator(ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
		return m.CompareValues(st, v1, v2)
	}))
	m.Max = 1
	st := &State{}
	res, handled := m.CompareValues(st, reflect.ValueOf(1), reflect.ValueOf(1))
	assert.True(t, handled)
	assert.SliceNotEmpty(t, res)
	assert.Zero(t, st.Depth)
}
