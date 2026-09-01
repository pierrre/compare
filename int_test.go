package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Int", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   123,
			V2:   123,
		},
		{
			Name: "Equal8",
			V1:   int8(123),
			V2:   int8(123),
		},
		{
			Name: "Equal16",
			V1:   int16(123),
			V2:   int16(123),
		},
		{
			Name: "Equal32",
			V1:   int32(123),
			V2:   int32(123),
		},
		{
			Name: "Equal64",
			V1:   int64(123),
			V2:   int64(123),
		},
		{
			Name: "Different",
			V1:   123,
			V2:   456,
		},
		{
			Name: "SupportDisabled",
			V1:   123,
			V2:   123,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Int}
			},
			IgnoreBenchmark: true,
		},
	})
}
