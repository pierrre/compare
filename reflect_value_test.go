package compare_test

import (
	"reflect"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

type testReflectValueUnexported struct {
	rv reflect.Value
}

type testStruct2 struct {
	Foo int
	Bar float64
}

func init() {
	comparetest.AddCasesPrefix("ReflectValue", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   reflect.ValueOf(1),
			V2:   reflect.ValueOf(1),
		},
		{
			Name: "Different",
			V1:   reflect.ValueOf(1),
			V2:   reflect.ValueOf(2),
		},
		{
			Name: "DifferentWrappedType",
			V1:   reflect.ValueOf(1),
			V2:   reflect.ValueOf("x"),
		},
		{
			Name: "DifferentWrappedStructType",
			V1: reflect.ValueOf(testStruct{
				Foo: 1,
			}),
			V2: reflect.ValueOf(testStruct2{
				Foo: 1,
			}),
		},
		{
			Name: "EqualStruct",
			V1: reflect.ValueOf(testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			}),
			V2: reflect.ValueOf(testStruct{
				Foo:        1,
				Bar:        2,
				unexported: 3,
			}),
		},
		{
			Name:            "EqualInvalid",
			V1:              reflect.ValueOf(reflect.Value{}),
			V2:              reflect.ValueOf(reflect.Value{}),
			IgnoreBenchmark: true,
		},
		{
			Name:            "InvalidMismatch",
			V1:              reflect.ValueOf(reflect.ValueOf(1)),
			V2:              reflect.ValueOf(reflect.Value{}),
			IgnoreBenchmark: true,
		},
		{
			Name: "SupportDisabled",
			V1:   reflect.ValueOf(1),
			V2:   reflect.ValueOf(1),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "UnexportedField",
			V1:   testReflectValueUnexported{rv: reflect.ValueOf(1)},
			V2:   testReflectValueUnexported{rv: reflect.ValueOf(1)},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.CanInterface = nil
			},
		},
		{
			Name: "UnexportedFieldDifferent",
			V1:   testReflectValueUnexported{rv: reflect.ValueOf(1)},
			V2:   testReflectValueUnexported{rv: reflect.ValueOf(2)},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.CanInterface = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.ReflectValue}
			},
			IgnoreBenchmark: true,
		},
	})
}
