package internal_test

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/ext/pierrrepretty"
	"github.com/pierrre/compare"
	. "github.com/pierrre/compare/internal"
)

func init() {
	// Prevent import cycle.
	assert.DeepEqualer = func(v1, v2 any) (diff string, equal bool) {
		res := compare.Compare(v1, v2)
		if len(res) == 0 {
			return "", true
		}
		diff = fmt.Sprintf("%+v", res)
		return diff, false
	}
	pierrrepretty.ConfigureDefault()
}

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
		typ: reflect.TypeOf(false),
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
		typ: reflect.TypeOf(int(0)),
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
		typ: reflect.TypeOf(uint(0)),
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
		typ: reflect.TypeOf(float64(0)),
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
		typ: reflect.TypeOf(complex(0, 0)),
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
		typ: reflect.TypeOf(""),
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
		typ: reflect.TypeOf(net.IP{}),
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
