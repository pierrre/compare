package comparetest_test

import (
	"testing"

	compare "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("CompareTest", []*comparetest.Case{
		{
			Name:                      "Equal",
			V1:                        123,
			V2:                        123,
			ConfigureComparator:       func(c *compare.Comparator) {},
			ConfigureValuesComparator: func(vc *compare.CommonComparator) {},
		},
		{
			Name: "Different",
			V1:   123,
			V2:   456,
		},
	})
}

func Test(t *testing.T) {
	comparetest.Test(t)
	testing.Benchmark(comparetest.Benchmark)
}
