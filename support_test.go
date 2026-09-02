package compare_test

import (
	"reflect"

	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Support", []*comparetest.Case{
		{
			Name: "Empty",
			V1:   123,
			V2:   123,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support.Checkers = nil
			},
		},
		{
			Name: "NotSupported",
			V1:   123,
			V2:   123,
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support.Checkers = []SupportChecker{
					SupportCheckerFunc(func(typ reflect.Type) ValuesComparator {
						if typ.Kind() == reflect.String {
							return NewStringComparator()
						}
						return nil
					}),
				}
			},
		},
		{
			Name: "Custom",
			V1:   "test",
			V2:   "test2",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support.Checkers = []SupportChecker{
					SupportCheckerFunc(func(typ reflect.Type) ValuesComparator {
						if typ.Kind() == reflect.String {
							return NewStringComparator()
						}
						return nil
					}),
				}
			},
		},
	})
}
