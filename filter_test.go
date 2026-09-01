package compare_test

import (
	"reflect"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Filter", []*comparetest.Case{
		{
			Name: "Match",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = []ValuesComparator{
					NewFilterComparator(
						ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
							return Res("handled", v1.Interface(), v2.Interface()), true
						}),
						FilterTypes(reflect.TypeFor[string]()),
					),
				}
			},
		},
		{
			Name: "MatchInterface",
			V1:   filterTestError("test"),
			V2:   filterTestError("test"),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = []ValuesComparator{
					NewFilterComparator(
						ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
							return Res("handled", v1.Interface(), v2.Interface()), true
						}),
						FilterTypes(reflect.TypeFor[error]()),
					),
				}
			},
		},
		{
			Name: "NoMatch",
			V1:   123,
			V2:   123,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = []ValuesComparator{
					NewFilterComparator(
						ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
							panic("should not be called")
						}),
						FilterTypes(reflect.TypeFor[string]()),
					),
				}
			},
		},
		{
			Name: "Nil",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = []ValuesComparator{
					NewFilterComparator(
						ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
							return Res("handled", v1.Interface(), v2.Interface()), true
						}),
						nil,
					),
				}
			},
		},
		{
			Name: "Standalone",
			V1:   123,
			V2:   123,
			ConfigureComparator: func(c *Comparator) {
				c.ValuesComparator = NewFilterComparator(NewIntComparator(), FilterTypes(reflect.TypeFor[int]()))
			},
		},
		{
			Name: "StandaloneNoMatch",
			V1:   123,
			V2:   456,
			ConfigureComparator: func(c *Comparator) {
				c.ValuesComparator = NewFilterComparator(NewIntComparator(), FilterTypes(reflect.TypeFor[string]()))
			},
		},
	})
}

type filterTestError string

func (e filterTestError) Error() string {
	return string(e)
}
