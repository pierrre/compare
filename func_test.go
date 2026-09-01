package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Func", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   testFunc1,
			V2:   testFunc1,
		},
		{
			Name:            "EqualNil",
			V1:              (func())(nil),
			V2:              (func())(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "Different",
			V1:              testFunc1,
			V2:              testFunc2,
			IgnoreBenchmark: true,
		},
		{
			Name: "NilMismatch",
			V1:   testFunc1,
			V2:   (func())(nil),
		},
		{
			Name: "SupportDisabled",
			V1:   testFunc1,
			V2:   testFunc1,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Func}
			},
			IgnoreBenchmark: true,
		},
	})
}

func testFunc1() {}

func testFunc2() {}
