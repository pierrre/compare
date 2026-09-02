package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("String", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   "test",
			V2:   "test",
		},
		{
			Name: "Different",
			V1:   "test",
			V2:   "test2",
		},
		{
			Name: "SupportDisabled",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   123,
			V2:   123,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.String}
			},
			IgnoreBenchmark: true,
		},
	})
}
