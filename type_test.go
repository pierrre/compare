package compare_test

import (
	"reflect"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	type testByType struct {
		Value int
	}
	comparetest.AddCasesPrefix("ByType", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   testByType{Value: 1},
			V2:   testByType{Value: 1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ByType[reflect.TypeFor[testByType]()] = ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					return nil, true
				})
			},
		},
		{
			Name: "Different",
			V1:   testByType{Value: 1},
			V2:   testByType{Value: 2},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ByType[reflect.TypeFor[testByType]()] = ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					return Res("custom", v1.Interface(), v2.Interface()), true
				})
			},
		},
		{
			Name: "NotHandled",
			V1:   123,
			V2:   456,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ByType[reflect.TypeFor[testByType]()] = ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					return Res("custom", v1.Interface(), v2.Interface()), true
				})
			},
		},
		{
			Name: "NotHandledEmpty",
			V1:   123,
			V2:   456,
			ConfigureComparator: func(c *Comparator) {
				c.ValuesComparator = NewByTypeComparators()
			},
		},
	})
}
