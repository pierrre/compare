package compare_test

import (
	"reflect"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

type testCanInterface struct {
	private int
}

func init() {
	comparetest.AddCasesPrefix("CanInterface", []*comparetest.Case{
		{
			Name: "UnexportedField",
			V1:   &testCanInterface{private: 1},
			V2:   &testCanInterface{private: 2},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ByType[reflect.TypeFor[int]()] = ValuesComparatorFunc(func(st *State, v1, v2 reflect.Value) (Result, bool) {
					return Res("custom", v1.Interface(), v2.Interface()), true
				})
			},
		},
		{
			Name: "Standalone",
			V1:   1,
			V2:   1,
			ConfigureComparator: func(c *Comparator) {
				c.ValuesComparator = NewCanInterfaceComparator(NewIntComparator())
			},
		},
		{
			Name: "UnexportedPointer",
			V1:   &testCanInterfacePtr{},
			V2:   &testCanInterfacePtr{},
		},
	})
}

type testCanInterfacePtr struct {
	privatePtr *int //nolint:unused // Read only through reflection during the comparison.
}
