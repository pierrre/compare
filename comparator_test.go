package compare_test

import (
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Comparator", []*comparetest.Case{
		{
			Name:            "EqualNotValid",
			V1:              nil,
			V2:              nil,
			IgnoreBenchmark: true,
		},
		{
			Name:            "NotEqualOnlyOneIsValid",
			V1:              nil,
			V2:              true,
			IgnoreBenchmark: true,
		},
		{
			Name: "NotEqualDifferentType",
			V1:   int32(1),
			V2:   int64(1),
		},
		{
			Name: "NotHandled",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Kind = nil
			},
		},
		{
			Name: "NotHandledSecond",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Int, vc.Kind.String}
			},
		},
	})
}

func Test(t *testing.T) {
	comparetest.Test(t)
}

func Benchmark(b *testing.B) {
	comparetest.Benchmark(b)
}

func TestCompare(t *testing.T) {
	res := Compare(123, 123)
	assert.SliceEmpty(t, res)
}
