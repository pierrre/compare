package compare_test

import (
	"reflect"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

type testRecursionNode struct {
	Next *testRecursionNode
}

func init() {
	comparetest.AddCasesPrefix("Recursion", []*comparetest.Case{
		{
			Name: "Equal",
			V1: &testRecursionNode{
				Next: &testRecursionNode{},
			},
			V2: &testRecursionNode{
				Next: &testRecursionNode{},
			},
		},
		{
			Name: "Different",
			V1: &testRecursionNode{
				Next: &testRecursionNode{Next: &testRecursionNode{}},
			},
			V2: &testRecursionNode{
				Next: &testRecursionNode{},
			},
		},
		{
			Name: "Recursion",
			V1: func() *testRecursionNode {
				n := &testRecursionNode{}
				n.Next = n
				return n
			}(),
			V2: func() *testRecursionNode {
				n := &testRecursionNode{}
				n.Next = n
				return n
			}(),
		},
		{
			Name: "NilSlices",
			V1:   [][]int{nil, nil},
			V2:   [][]int{nil, nil},
		},
		{
			Name: "Standalone",
			V1:   1,
			V2:   1,
			ConfigureComparator: func(c *Comparator) {
				c.ValuesComparator = NewRecursionComparator(NewIntComparator())
			},
		},
		{
			Name: "StandaloneRecursion",
			V1: func() *testRecursionNode {
				n := &testRecursionNode{}
				n.Next = n
				return n
			}(),
			V2: func() *testRecursionNode {
				n := &testRecursionNode{}
				n.Next = n
				return n
			}(),
			ConfigureComparator: func(c *Comparator) {
				var r *RecursionComparator
				r = NewRecursionComparator(ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					return r.CompareValues(st, v1, v2)
				}))
				c.ValuesComparator = r
			},
		},
	})
}
