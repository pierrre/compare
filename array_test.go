package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Array", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   [3]int{1, 2, 3},
			V2:   [3]int{1, 2, 3},
		},
		{
			Name: "Different",
			V1:   [3]int{1, 2, 3},
			V2:   [3]int{1, 0, 3},
		},
		{
			Name: "MaxDifferences",
			V1:   [20]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			V2:   [20]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name: "MaxDifferencesDisabled",
			V1:   [20]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			V2:   [20]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Array.MaxDifferences = 0
			},
		},
		{
			Name: "MaxDifferencesCustom",
			V1:   [20]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			V2:   [20]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Array.MaxDifferences = 3
			},
		},
		{
			Name: "SupportDisabled",
			V1:   [3]int{1, 2, 3},
			V2:   [3]int{1, 2, 3},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Array}
			},
			IgnoreBenchmark: true,
		},
	})
}
