package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Interface", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   [1]any{1},
			V2:   [1]any{1},
		},
		{
			Name:            "EqualNil",
			V1:              [1]any{nil},
			V2:              [1]any{nil},
			IgnoreBenchmark: true,
		},
		{
			Name: "NilMismatch",
			V1:   [1]any{1},
			V2:   [1]any{nil},
		},
		{
			Name: "DifferentType",
			V1:   [1]any{1},
			V2:   [1]any{"a"},
		},
		{
			Name: "Different",
			V1:   [1]any{1},
			V2:   [1]any{2},
		},
		{
			Name: "SupportDisabled",
			V1:   [1]any{1},
			V2:   [1]any{1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Interface}
			},
			IgnoreBenchmark: true,
		},
	})
}
