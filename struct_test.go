package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Struct", []*comparetest.Case{
		{
			Name: "Equal",
			V1: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			},
			V2: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			},
		},
		{
			Name: "DifferentExported",
			V1: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			},
			V2: testStruct{
				Foo:        2,
				Bar:        2,
				unexported: 3,
			},
		},
		{
			Name: "DifferentUnexported",
			V1: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			},
			V2: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 4,
			},
		},
		{
			Name: "SupportDisabled",
			V1: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			},
			V2: testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Struct}
			},
			IgnoreBenchmark: true,
		},
	})
}

type testStruct struct {
	Foo        int
	Bar        float64
	unexported int
}
