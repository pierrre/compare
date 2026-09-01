package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Uint", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   uint(123),
			V2:   uint(123),
		},
		{
			Name: "Equal8",
			V1:   uint8(123),
			V2:   uint8(123),
		},
		{
			Name: "Equal16",
			V1:   uint16(123),
			V2:   uint16(123),
		},
		{
			Name: "Equal32",
			V1:   uint32(123),
			V2:   uint32(123),
		},
		{
			Name: "Equal64",
			V1:   uint64(123),
			V2:   uint64(123),
		},
		{
			Name: "EqualPtr",
			V1:   uintptr(123),
			V2:   uintptr(123),
		},
		{
			Name: "Different",
			V1:   uint(123),
			V2:   uint(456),
		},
		{
			Name: "SupportDisabled",
			V1:   uint(123),
			V2:   uint(123),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Uint}
			},
			IgnoreBenchmark: true,
		},
	})
}
