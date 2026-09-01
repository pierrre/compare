package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Bool", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   true,
			V2:   true,
		},
		{
			Name: "Different",
			V1:   true,
			V2:   false,
		},
		{
			Name: "SupportDisabled",
			V1:   true,
			V2:   true,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Bool}
			},
			IgnoreBenchmark: true,
		},
	})
}
