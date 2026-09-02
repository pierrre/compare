package compare_test

import (
	"math"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Float", []*comparetest.Case{
		{
			Name: "Equal32",
			V1:   float32(1),
			V2:   float32(1),
		},
		{
			Name: "Equal64",
			V1:   float64(1),
			V2:   float64(1),
		},
		{
			Name: "Different",
			V1:   float64(1),
			V2:   float64(2),
		},
		{
			Name: "NaNEqual",
			V1:   math.NaN(),
			V2:   math.NaN(),
		},
		{
			Name: "NaNNotEqual",
			V1:   math.NaN(),
			V2:   math.NaN(),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Float.NaNEqual = false
			},
		},
		{
			Name: "SignedZeroEqual",
			V1:   0.0,
			V2:   math.Copysign(0, -1),
		},
		{
			Name: "SignedZeroNotEqual",
			V1:   0.0,
			V2:   math.Copysign(0, -1),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Float.SignedZeroEqual = false
			},
		},
		{
			Name: "SignedZeroNonZero",
			V1:   1.0,
			V2:   1.0,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind.Float.SignedZeroEqual = false
			},
		},
		{
			Name: "SupportDisabled",
			V1:   float64(1),
			V2:   float64(1),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Float}
			},
			IgnoreBenchmark: true,
		},
	})
}
