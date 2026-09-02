package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Pointer", []*comparetest.Case{
		{
			Name: "Equal",
			V1: func() *int {
				i := 1
				return &i
			}(),
			V2: func() *int {
				i := 1
				return &i
			}(),
		},
		{
			Name:            "EqualSame",
			V1:              &testIntPtr,
			V2:              &testIntPtr,
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualNil",
			V1:              (*int)(nil),
			V2:              (*int)(nil),
			IgnoreBenchmark: true,
		},
		{
			Name: "NilMismatch",
			V1:   &testIntPtr,
			V2:   (*int)(nil),
		},
		{
			Name: "Different",
			V1: func() *int {
				i := 1
				return &i
			}(),
			V2: func() *int {
				i := 2
				return &i
			}(),
		},
		{
			Name: "SupportDisabled",
			V1:   &testIntPtr,
			V2:   &testIntPtr,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Pointer}
			},
			IgnoreBenchmark: true,
		},
	})
}

var testIntPtr = 1
