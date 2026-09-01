package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Compare", []*comparetest.Case{
		{
			Name: "Equal",
			V1: compareValueStruct{
				Foo: 1,
				bar: 1,
			},
			V2: compareValueStruct{
				Foo: 1,
				bar: 2,
			},
		},
		{
			Name: "Different",
			V1: compareValueStruct{
				Foo: 1,
			},
			V2: compareValueStruct{
				Foo: 2,
			},
		},
		{
			Name: "PointerEqual",
			V1: &comparePointerStruct{
				Foo: 1,
			},
			V2: &comparePointerStruct{
				Foo: 1,
			},
		},
		{
			Name: "PointerDifferent",
			V1: &comparePointerStruct{
				Foo: 1,
			},
			V2: &comparePointerStruct{
				Foo: 2,
			},
		},
		{
			Name: "PointerEqualNil",
			V1:   (*comparePointerStruct)(nil),
			V2:   (*comparePointerStruct)(nil),
		},
		{
			Name: "PointerDifferentNil",
			V1:   (*comparePointerStruct)(nil),
			V2:   &comparePointerStruct{Foo: 1},
		},
		{
			Name: "Cmp",
			V1: cmpValueStruct{
				Foo: 1,
			},
			V2: cmpValueStruct{
				Foo: 1,
			},
		},
		{
			Name: "CmpDifferent",
			V1: cmpValueStruct{
				Foo: 1,
			},
			V2: cmpValueStruct{
				Foo: 2,
			},
		},
		{
			Name: "Custom",
			V1: customCompareStruct{
				Foo: 1,
			},
			V2: customCompareStruct{
				Foo: 1,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Compare.ValuesComparator.CompareNames = []string{"CustomCompare"}
			},
		},
		{
			Name: "SupportDisabled",
			V1: compareValueStruct{
				Foo: 1,
			},
			V2: compareValueStruct{
				Foo: 1,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "WrongSignature",
			V1:   compareWrongSignatureStruct{Foo: 1},
			V2:   compareWrongSignatureStruct{Foo: 1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   testStruct{Foo: 1},
			V2:   testStruct{Foo: 1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Compare}
			},
			IgnoreBenchmark: true,
		},
		{
			Name: "Disabled",
			V1: compareValueStruct{
				Foo: 1,
				bar: 1,
			},
			V2: compareValueStruct{
				Foo: 1,
				bar: 2,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Compare = nil
			},
		},
	})
}

type compareValueStruct struct {
	Foo int
	bar int
}

func (c compareValueStruct) Compare(o compareValueStruct) int {
	return c.Foo - o.Foo
}

type comparePointerStruct struct {
	Foo int
}

func (c *comparePointerStruct) Compare(o *comparePointerStruct) int {
	return c.Foo - o.Foo
}

type cmpValueStruct struct {
	Foo int
}

func (c cmpValueStruct) Cmp(o cmpValueStruct) int {
	return c.Foo - o.Foo
}

type customCompareStruct struct {
	Foo int
}

func (c customCompareStruct) CustomCompare(o customCompareStruct) int {
	return c.Foo - o.Foo
}

type compareWrongSignatureStruct struct {
	Foo int
}

func (c compareWrongSignatureStruct) Compare(o compareWrongSignatureStruct, _ int) int {
	return c.Foo - o.Foo
}
