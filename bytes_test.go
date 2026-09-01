package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Bytes", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   []byte("test"),
			V2:   []byte("test"),
		},
		{
			Name: "Different",
			V1:   []byte("test"),
			V2:   []byte("test2"),
		},
		{
			Name:            "EqualNil",
			V1:              []byte(nil),
			V2:              []byte(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualEmpty",
			V1:              []byte{},
			V2:              []byte{},
			IgnoreBenchmark: true,
		},
		{
			Name: "NilMismatch",
			V1:   []byte(nil),
			V2:   []byte{1},
		},
		{
			Name: "SupportDisabled",
			V1:   []byte("test"),
			V2:   []byte("test2"),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   []int{1, 2},
			V2:   []int{1, 2},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Bytes}
			},
			IgnoreBenchmark: true,
		},
	})
}
