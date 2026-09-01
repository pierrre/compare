package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Slice", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   []int{1, 2, 3},
			V2:   []int{1, 2, 3},
		},
		{
			Name:            "EqualNil",
			V1:              []int(nil),
			V2:              []int(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualEmpty",
			V1:              []int{},
			V2:              []int{},
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualSame",
			V1:              testSlice,
			V2:              testSlice,
			IgnoreBenchmark: true,
		},
		{
			Name: "Overlapping",
			V1:   testSlice[:2],
			V2:   testSlice[1:3],
		},
		{
			Name: "SameBackingDifferentLength",
			V1:   testSlice,
			V2:   testSlice[:2],
		},
		{
			Name: "NilMismatch",
			V1:   []int{1},
			V2:   []int(nil),
		},
		{
			Name: "LengthDifferent",
			V1:   []int{1, 2},
			V2:   []int{1, 2, 3},
		},
		{
			Name: "Different",
			V1:   []int{1, 2, 3},
			V2:   []int{1, 0, 3},
		},
		{
			Name: "MaxDifferences",
			V1:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			V2:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name: "SupportDisabled",
			V1:   []int{1, 2, 3},
			V2:   []int{1, 2, 3},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Slice}
			},
			IgnoreBenchmark: true,
		},
	})
}

var testSlice = []int{1, 2, 3}
