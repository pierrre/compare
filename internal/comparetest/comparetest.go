// Package comparetest provides utilities for testing the compare package.
package comparetest

import (
	"testing"

	"github.com/pierrre/assert/assertauto"
	compare "github.com/pierrre/compare"
)

// Case represents a test case.
type Case struct {
	Name                      string
	V1, V2                    any
	ConfigureComparator       func(c *compare.Comparator)
	ConfigureValuesComparator func(vc *compare.CommonComparator)
	IgnoreResult              bool
	IgnoreAllocs              bool
	IgnoreBenchmark           bool
}

func (tc *Case) newComparator() *compare.Comparator {
	vc := compare.NewCommonComparator()
	if tc.ConfigureValuesComparator != nil {
		tc.ConfigureValuesComparator(vc)
	}
	c := compare.NewComparator(vc)
	if tc.ConfigureComparator != nil {
		tc.ConfigureComparator(c)
	}
	return c
}

var testCases []*Case

// AddCases adds test cases to the list of test cases.
func AddCases(tcs []*Case) {
	testCases = append(testCases, tcs...)
}

// AddCasesPrefix adds test cases with a prefix to the list of test cases.
func AddCasesPrefix(prefix string, tcs []*Case) {
	for _, tc := range tcs {
		tc.Name = prefix + "/" + tc.Name
	}
	AddCases(tcs)
}

// Test runs the tests.
func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			c := tc.newComparator()
			res := c.Compare(tc.V1, tc.V2)
			t.Log(res)
			if !tc.IgnoreResult {
				assertauto.Equal(t, res, assertauto.ValueStringer(func(a any) string {
					r, _ := a.(compare.Result)
					if len(r) == 0 {
						return "equal"
					}
					return r.String()
				}))
			}
			if !tc.IgnoreAllocs {
				allocs, _ := assertauto.AllocsPerRun(t, 100, func() {
					t.Helper()
					_ = c.Compare(tc.V1, tc.V2)
				})
				t.Logf("allocs: %g", allocs)
			}
		})
	}
}

// Benchmark runs the benchmarks.
func Benchmark(b *testing.B) {
	for _, tc := range testCases {
		if !tc.IgnoreBenchmark {
			b.Run(tc.Name, func(b *testing.B) {
				c := tc.newComparator()
				for b.Loop() {
					c.Compare(tc.V1, tc.V2)
				}
			})
		}
	}
}
