package compare_test

import (
	"time"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Equal", []*comparetest.Case{
		{
			Name: "Equal",
			V1: equalValueStruct{
				Foo: 1,
				bar: 1,
			},
			V2: equalValueStruct{
				Foo: 1,
				bar: 2,
			},
		},
		{
			Name: "Different",
			V1: equalValueStruct{
				Foo: 1,
			},
			V2: equalValueStruct{
				Foo: 2,
			},
		},
		{
			Name: "PointerEqual",
			V1: &equalPointerStruct{
				Foo: 1,
			},
			V2: &equalPointerStruct{
				Foo: 1,
			},
		},
		{
			Name: "PointerDifferent",
			V1: &equalPointerStruct{
				Foo: 1,
			},
			V2: &equalPointerStruct{
				Foo: 2,
			},
		},
		{
			Name: "PointerEqualNil",
			V1:   (*equalPointerStruct)(nil),
			V2:   (*equalPointerStruct)(nil),
		},
		{
			Name: "PointerDifferentNil",
			V1:   (*equalPointerStruct)(nil),
			V2:   &equalPointerStruct{Foo: 1},
		},
		{
			Name: "Eq",
			V1: eqValueStruct{
				Foo: 1,
			},
			V2: eqValueStruct{
				Foo: 1,
			},
		},
		{
			Name: "EqDifferent",
			V1: eqValueStruct{
				Foo: 1,
			},
			V2: eqValueStruct{
				Foo: 2,
			},
		},
		{
			Name: "Custom",
			V1: customValueStruct{
				Foo: 1,
			},
			V2: customValueStruct{
				Foo: 1,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Equal.ValuesComparator.Names = []string{"Custom"}
			},
		},
		{
			Name: "UnexportedField",
			V1: equalNestedUnexportedStruct{
				w: equalValueStruct{
					Foo: 1,
					bar: 1,
				},
			},
			V2: equalNestedUnexportedStruct{
				w: equalValueStruct{
					Foo: 1,
					bar: 2,
				},
			},
		},
		{
			Name: "UnexportedFieldDifferent",
			V1: equalNestedUnexportedStruct{
				w: equalValueStruct{
					Foo: 1,
				},
			},
			V2: equalNestedUnexportedStruct{
				w: equalValueStruct{
					Foo: 2,
				},
			},
		},
		{
			Name: "UnexportedTime",
			V1:   unexportedTimeStruct{t: time.Unix(0, 0)},
			V2:   unexportedTimeStruct{t: time.Unix(0, 0)},
		},
		{
			Name: "UnexportedTimeDifferent",
			V1:   unexportedTimeStruct{t: time.Unix(0, 0)},
			V2:   unexportedTimeStruct{t: time.Unix(1, 0)},
		},
		{
			Name: "SupportDisabled",
			V1: equalValueStruct{
				Foo: 1,
			},
			V2: equalValueStruct{
				Foo: 1,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "WrongSignature",
			V1:   equalWrongSignatureStruct{Foo: 1},
			V2:   equalWrongSignatureStruct{Foo: 1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   testStruct{Foo: 1},
			V2:   testStruct{Foo: 1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Equal}
			},
			IgnoreBenchmark: true,
		},
		{
			Name: "Disabled",
			V1: equalValueStruct{
				Foo: 1,
				bar: 1,
			},
			V2: equalValueStruct{
				Foo: 1,
				bar: 2,
			},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Equal = nil
			},
		},
	})
}

type equalValueStruct struct {
	Foo int
	bar int
}

func (e equalValueStruct) Equal(o equalValueStruct) bool {
	return e.Foo == o.Foo
}

type equalPointerStruct struct {
	Foo int
}

func (e *equalPointerStruct) Equal(o *equalPointerStruct) bool {
	return e.Foo == o.Foo
}

type eqValueStruct struct {
	Foo int
}

func (e eqValueStruct) Eq(o eqValueStruct) bool {
	return e.Foo == o.Foo
}

type customValueStruct struct {
	Foo int
}

func (e customValueStruct) Custom(o customValueStruct) bool {
	return e.Foo == o.Foo
}

type equalWrongSignatureStruct struct {
	Foo int
}

func (e equalWrongSignatureStruct) Equal(o equalWrongSignatureStruct, _ int) bool {
	return e.Foo == o.Foo
}

type equalNestedUnexportedStruct struct {
	w equalValueStruct
}

type unexportedTimeStruct struct {
	t time.Time
}
