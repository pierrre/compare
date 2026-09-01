package compare_test

import (
	"unsafe" //nolint:depguard // Required for test.

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	p := new(123)
	comparetest.AddCasesPrefix("UnsafePointer", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   unsafe.Pointer(p),
			V2:   unsafe.Pointer(p),
		},
		{
			Name:            "EqualNil",
			V1:              unsafe.Pointer(nil),
			V2:              unsafe.Pointer(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "Different",
			V1:              unsafe.Pointer(new(123)),
			V2:              unsafe.Pointer(new(456)),
			IgnoreResult:    true,
			IgnoreBenchmark: true,
		},
		{
			Name: "SupportDisabled",
			V1:   unsafe.Pointer(p),
			V2:   unsafe.Pointer(p),
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.UnsafePointer}
			},
			IgnoreBenchmark: true,
		},
	})
}
