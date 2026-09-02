package compare_test

import (
	"math"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Complex", []*comparetest.Case{
		{
			Name: "Equal64",
			V1:   complex64(1 + 2i),
			V2:   complex64(1 + 2i),
		},
		{
			Name: "Equal128",
			V1:   complex(1, 2),
			V2:   complex(1, 2),
		},
		{
			Name: "Different",
			V1:   complex(1, 2),
			V2:   complex(3, 4),
		},
		{
			Name: "NaNEqual",
			V1:   complex(math.NaN(), 1),
			V2:   complex(math.NaN(), 1),
		},
		{
			Name: "NaNNotEqual",
			V1:   complex(math.NaN(), 1),
			V2:   complex(math.NaN(), 1),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Complex.NaNEqual = false
			},
		},
		{
			Name: "SignedZeroEqual",
			V1:   complex(0, 0),
			V2:   complex(math.Copysign(0, -1), 0),
		},
		{
			Name: "SignedZeroNotEqual",
			V1:   complex(0, 0),
			V2:   complex(math.Copysign(0, -1), 0),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Complex.SignedZeroEqual = false
			},
		},
		{
			Name: "SignedZeroNonZero",
			V1:   complex(1, 1),
			V2:   complex(1, 1),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Complex.SignedZeroEqual = false
			},
		},
		{
			Name: "SupportDisabled",
			V1:   complex(1, 2),
			V2:   complex(1, 2),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Complex}
			},
			IgnoreBenchmark: true,
		},
	})
}
