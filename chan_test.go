package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	testChan := make(chan int)
	comparetest.AddCasesPrefix("Chan", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   make(chan int),
			V2:   make(chan int),
		},
		{
			Name:            "EqualSame",
			V1:              testChan,
			V2:              testChan,
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualNil",
			V1:              chan int(nil),
			V2:              chan int(nil),
			IgnoreBenchmark: true,
		},
		{
			Name: "NilMismatch",
			V1:   make(chan int),
			V2:   chan int(nil),
		},
		{
			Name: "CapacityDifferent",
			V1:   make(chan int, 1),
			V2:   make(chan int, 2),
		},
		{
			Name: "LengthDifferent",
			V1:   make(chan int, 1),
			V2: func() chan int {
				chn := make(chan int, 1)
				chn <- 1
				return chn
			}(),
		},
		{
			Name: "SupportDisabled",
			V1:   make(chan int),
			V2:   make(chan int),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Chan}
			},
			IgnoreBenchmark: true,
		},
	})
}
