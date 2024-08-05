package internal_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	. "github.com/pierrre/compare/internal"
)

var sortMapsKeysTestCases = []struct {
	name     string
	values   []reflect.Value
	typ      reflect.Type
	expected []reflect.Value
}{
	{
		name: "Bool",
		values: []reflect.Value{
			reflect.ValueOf(true),
			reflect.ValueOf(false),
		},
		typ: reflect.TypeFor[bool](),
		expected: []reflect.Value{
			reflect.ValueOf(false),
			reflect.ValueOf(true),
		},
	},
	{
		name: "Int",
		values: []reflect.Value{
			reflect.ValueOf(int(2)),
			reflect.ValueOf(int(1)),
		},
		typ: reflect.TypeFor[int](),
		expected: []reflect.Value{
			reflect.ValueOf(int(1)),
			reflect.ValueOf(int(2)),
		},
	},
	{
		name: "Uint",
		values: []reflect.Value{
			reflect.ValueOf(uint(2)),
			reflect.ValueOf(uint(1)),
		},
		typ: reflect.TypeFor[uint](),
		expected: []reflect.Value{
			reflect.ValueOf(uint(1)),
			reflect.ValueOf(uint(2)),
		},
	},
	{
		name: "Float",
		values: []reflect.Value{
			reflect.ValueOf(float64(2)),
			reflect.ValueOf(float64(1)),
		},
		typ: reflect.TypeFor[float64](),
		expected: []reflect.Value{
			reflect.ValueOf(float64(1)),
			reflect.ValueOf(float64(2)),
		},
	},
	{
		name: "Complex",
		values: []reflect.Value{
			reflect.ValueOf(complex(1, 1)),
			reflect.ValueOf(complex(2, 2)),
			reflect.ValueOf(complex(2, 1)),
			reflect.ValueOf(complex(1, 2)),
		},
		typ: reflect.TypeFor[complex128](),
		expected: []reflect.Value{
			reflect.ValueOf(complex(1, 1)),
			reflect.ValueOf(complex(1, 2)),
			reflect.ValueOf(complex(2, 1)),
			reflect.ValueOf(complex(2, 2)),
		},
	},
	{
		name: "String",
		values: []reflect.Value{
			reflect.ValueOf("b"),
			reflect.ValueOf("a"),
		},
		typ: reflect.TypeFor[string](),
		expected: []reflect.Value{
			reflect.ValueOf("a"),
			reflect.ValueOf("b"),
		},
	},
	{
		name: "NetIP",
		values: []reflect.Value{
			reflect.ValueOf(net.ParseIP("2.2.2.2")),
			reflect.ValueOf(net.ParseIP("1.1.1.1")),
		},
		typ: reflect.TypeFor[net.IP](),
		expected: []reflect.Value{
			reflect.ValueOf(net.ParseIP("1.1.1.1")),
			reflect.ValueOf(net.ParseIP("2.2.2.2")),
		},
	},
}

func TestSortMapsKeys(t *testing.T) {
	for _, tc := range sortMapsKeysTestCases {
		t.Run(tc.name, func(t *testing.T) {
			SortMapsKeys(tc.typ, tc.values)
			assert.DeepEqual(t, tc.values, tc.expected)
		})
	}
}
