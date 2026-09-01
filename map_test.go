package compare_test

import (
	. "github.com/pierrre/compare"
	"github.com/pierrre/compare/internal/comparetest"
)

func init() {
	comparetest.AddCasesPrefix("Map", []*comparetest.Case{
		{
			Name: "Equal",
			V1:   map[string]int{"a": 1, "b": 2},
			V2:   map[string]int{"b": 2, "a": 1},
		},
		{
			Name:            "EqualNil",
			V1:              map[string]int(nil),
			V2:              map[string]int(nil),
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualEmpty",
			V1:              map[string]int{},
			V2:              map[string]int{},
			IgnoreBenchmark: true,
		},
		{
			Name:            "EqualSame",
			V1:              testMap,
			V2:              testMap,
			IgnoreBenchmark: true,
		},
		{
			Name: "NilMismatch",
			V1:   map[string]int{"a": 1},
			V2:   map[string]int(nil),
		},
		{
			Name: "Different",
			V1:   map[string]int{"a": 1},
			V2:   map[string]int{"a": 2},
		},
		{
			Name: "KeyMissing",
			V1:   map[string]int{"a": 1, "b": 2},
			V2:   map[string]int{"a": 1, "c": 2},
		},
		{
			Name: "KeyMissingReversed",
			V1:   map[string]int{"b": 1},
			V2:   map[string]int{"a": 1},
		},
		{
			Name: "LengthDifferent",
			V1:   map[string]int{"a": 1},
			V2:   map[string]int{"a": 1, "b": 2},
		},
		{
			Name: "MaxDifferences",
			V1:   map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16, "q": 17, "r": 18, "s": 19, "t": 20},
			V2:   map[string]int{"a": 0, "b": 0, "c": 0, "d": 0, "e": 0, "f": 0, "g": 0, "h": 0, "i": 0, "j": 0, "k": 0, "l": 0, "m": 0, "n": 0, "o": 0, "p": 0, "q": 0, "r": 0, "s": 0, "t": 0},
		},
		{
			Name: "MaxDifferencesKeyMissing",
			V1:   map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16, "q": 17, "r": 18, "s": 19, "t": 20},
			V2:   map[string]int{"A": 1, "B": 2, "C": 3, "D": 4, "E": 5, "F": 6, "G": 7, "H": 8, "I": 9, "J": 10, "K": 11, "L": 12, "M": 13, "N": 14, "O": 15, "P": 16, "Q": 17, "R": 18, "S": 19, "T": 20},
		},
		{
			Name: "SupportDisabled",
			V1:   map[string]int{"a": 1},
			V2:   map[string]int{"a": 1},
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.Support = nil
			},
		},
		{
			Name: "Not",
			V1:   "test",
			V2:   "test",
			ConfigureValuesComparator: func(vc *CommonComparator) {
				vc.ValuesComparators = ValuesComparators{vc.Kind.Map}
			},
			IgnoreBenchmark: true,
		},
	})
}

var testMap = map[string]int{"a": 1}
