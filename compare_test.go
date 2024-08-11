package compare_test

import (
	"fmt"
	"io"
	"math/big"
	"net"
	"reflect"
	"testing"
	"time"
	"unsafe" //nolint:depguard // Used for unsafe.Pointer comparison.

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	. "github.com/pierrre/compare"
)

func ExampleCompare() {
	type T struct {
		String string
		Int    int
		Map    map[string]any
		Slice  []int
	}
	v1 := T{
		String: "aaa",
		Int:    1,
		Map: map[string]any{
			"a": "a",
			"b": "b",
			"c": "c",
		},
		Slice: []int{1, 2, 3},
	}
	v2 := T{
		String: "bbb",
		Int:    2,
		Map: map[string]any{
			"a": "z",
			"b": 5,
			"d": "c",
		},
		Slice: []int{1, 2, 4},
	}
	diff := Compare(v1, v2)
	if len(diff) != 0 {
		fmt.Printf("%+v", diff)
	}
	// Output:
	// .String: string not equal
	// 	v1="aaa"
	// 	v2="bbb"
	// .Int: int not equal
	// 	v1=1
	// 	v2=2
	// .Map[a]: string not equal
	// 	v1="a"
	// 	v2="z"
	// .Map[b]: type not equal
	// 	v1=string
	// 	v2=int
	// .Map[c]: map key not defined
	// 	v1=true
	// 	v2=false
	// .Map[d]: map key not defined
	// 	v1=false
	// 	v2=true
	// .Slice[2]: int not equal
	// 	v1=3
	// 	v2=4
}

type testCase struct {
	name      string
	v1        any
	v2        any
	configure func(c *Comparator)
}

func (tc testCase) newComparator() *Comparator {
	c := NewComparator()
	if tc.configure != nil {
		tc.configure(c)
	}
	return c
}

var compareTestCases = []testCase{
	{
		name: "EqualNotValid",
		v1:   nil,
		v2:   nil,
	},
	{
		name: "NotEqualOnlyOneIsValid",
		v1:   nil,
		v2:   true,
	},
	{
		name: "NotEqualDifferentType",
		v1:   int32(1),
		v2:   int64(1),
	},
	{
		name: "BoolEqual",
		v1:   true,
		v2:   true,
	},
	{
		name: "BoolNotEqual",
		v1:   true,
		v2:   false,
	},
	{
		name: "IntEqual",
		v1:   int(1),
		v2:   int(1),
	},
	{
		name: "IntNotEqual",
		v1:   int(1),
		v2:   int(2),
	},
	{
		name: "UintEqual",
		v1:   uint(1),
		v2:   uint(1),
	},
	{
		name: "UintNotEqual",
		v1:   uint(1),
		v2:   uint(2),
	},
	{
		name: "FloatEqual",
		v1:   float64(1),
		v2:   float64(1),
	},
	{
		name: "FloatNotEqual",
		v1:   float64(1),
		v2:   float64(2),
	},
	{
		name: "ComplexEqual",
		v1:   complex(1, 1),
		v2:   complex(1, 1),
	},
	{
		name: "ComplexNotEqual",
		v1:   complex(1, 1),
		v2:   complex(2, 2),
	},
	{
		name: "StringEqual",
		v1:   "a",
		v2:   "a",
	},
	{
		name: "StringNotEqual",
		v1:   "a",
		v2:   "b",
	},
	{
		name: "ArrayEqual",
		v1:   [3]int{1, 2, 3},
		v2:   [3]int{1, 2, 3},
	},
	{
		name: "ArrayNotEqual",
		v1:   [3]int{1, 2, 3},
		v2:   [3]int{1, 0, 3},
	},
	{
		name: "SliceEqual",
		v1:   []int{1, 2, 3},
		v2:   []int{1, 2, 3},
	},
	{
		name: "SliceEqualNil",
		v1:   []int(nil),
		v2:   []int(nil),
	},
	{
		name: "SliceEqualLengthZero",
		v1:   []int{},
		v2:   []int{},
	},
	{
		name: "SliceEqualPointer",
		v1:   testSlice,
		v2:   testSlice,
	},
	{
		name: "SliceByteEqual",
		v1:   make([]byte, 1<<20),
		v2:   make([]byte, 1<<20),
	},
	{
		name: "SliceEqualRecursion",
		v1: func() []any {
			v := make([]any, 1)
			v[0] = v
			return v
		}(),
		v2: func() any {
			v := make([]any, 1)
			v[0] = v
			return v
		}(),
	},
	{
		name: "SliceNotEqual",
		v1:   []int{1, 2, 3},
		v2:   []int{1, 0, 3},
	},
	{
		name: "SliceNotEqualOnlyOneIsNil",
		v1:   []int{1, 2, 3},
		v2:   []int(nil),
	},
	{
		name: "SliceNotEqualLength",
		v1:   []int{1, 2},
		v2:   []int{1, 2, 3},
	},
	{
		name: "SliceByteNotEqual",
		v1:   make([]byte, 10),
		v2: func() []byte {
			s := make([]byte, 10)
			s[0] = 1
			return s
		}(),
	},
	{
		name: "SliceNotEqualMaxDifferences",
		v1: func() []int {
			s := make([]int, DefaultComparator.SliceMaxDifferences*2)
			for i := range s {
				s[i] = i
			}
			return s
		}(),
		v2: func() []int {
			s := make([]int, DefaultComparator.SliceMaxDifferences*2)
			for i := range s {
				s[i] = i + 1
			}
			return s
		}(),
	},
	{
		name: "InterfaceEqual",
		v1:   [1]any{1},
		v2:   [1]any{1},
	},
	{
		name: "InterfaceEqualNil",
		v1:   [1]any{nil},
		v2:   [1]any{nil},
	},
	{
		name: "InterfaceNotEqualOnlyOneIsNil",
		v1:   [1]any{1},
		v2:   [1]any{nil},
	},
	{
		name: "PointerEqual",
		v1: func() *int {
			i := 1
			return &i
		}(),
		v2: func() *int {
			i := 1
			return &i
		}(),
	},
	{
		name: "PointerEqualPointer",
		v1:   &testInt,
		v2:   &testInt,
	},
	{
		name: "PointerEqualRecursion",
		v1: func() any {
			var v any
			v = &v
			return v
		}(),
		v2: func() any {
			var v any
			v = &v
			return v
		}(),
	},
	{
		name: "PointerNotEqual",
		v1: func() *int {
			i := 1
			return &i
		}(),
		v2: func() *int {
			i := 2
			return &i
		}(),
	},
	{
		name: "StructEqual",
		v1: &testStruct{
			Exported:   1,
			unexported: 2,
		},
		v2: &testStruct{
			Exported:   1,
			unexported: 2,
		},
	},
	{
		name: "StructNotEqualExported",
		v1: &testStruct{
			Exported:   1,
			unexported: 2,
		},
		v2: &testStruct{
			Exported:   2,
			unexported: 2,
		},
	},
	{
		name: "StructNotEqualUnexported",
		v1: &testStruct{
			Exported:   1,
			unexported: 1,
		},
		v2: &testStruct{
			Exported:   1,
			unexported: 2,
		},
	},
	{
		name: "MapEqual",
		v1: map[string]int{
			"i": 1,
			"j": 2,
		},
		v2: map[string]int{
			"j": 2,
			"i": 1,
		},
	},
	{
		name: "MapEqualNil",
		v1:   map[string]int(nil),
		v2:   map[string]int(nil),
	},
	{
		name: "MapEqualLengthZero",
		v1:   map[string]int{},
		v2:   map[string]int{},
	},
	{
		name: "MapEqualPointer",
		v1:   testMap,
		v2:   testMap,
	},
	{
		name: "MapEqualRecursion",
		v1: func() map[int]any {
			v := make(map[int]any)
			v[0] = v
			return v
		}(),
		v2: func() any {
			v := make(map[int]any)
			v[0] = v
			return v
		}(),
	},
	{
		name: "MapNotEqualValue",
		v1: map[string]int{
			"i": 1,
		},
		v2: map[string]int{
			"i": 2,
		},
	},
	{
		name: "MapNotEqualKey",
		v1: map[string]int{
			"a": 1,
		},
		v2: map[string]int{
			"b": 1,
		},
	},
	{
		name: "MapNotEqualOnlyOneIsNil",
		v1: map[string]int{
			"i": 1,
		},
		v2: map[string]int(nil),
	},
	{
		name: "MapNotEqualLength",
		v1: map[string]int{
			"i": 1,
		},
		v2: map[string]int{
			"i": 1,
			"j": 1,
		},
	},
	{
		name: "MapNotEqualMaxDifferences",
		v1: func() map[int]int {
			m := make(map[int]int)
			for i := range DefaultComparator.SliceMaxDifferences * 2 {
				m[i] = i
			}
			return m
		}(),
		v2: func() map[int]int {
			m := make(map[int]int)
			for i := range DefaultComparator.SliceMaxDifferences * 2 {
				m[i] = i + 1
			}
			return m
		}(),
	},
	{
		name: "UnsafePointerEqual",
		v1:   unsafe.Pointer(&testInt),
		v2:   unsafe.Pointer(&testInt),
	},
	{
		name: "ChanEqual",
		v1:   make(chan int),
		v2:   make(chan int),
	},
	{
		name: "ChanEqualNil",
		v1:   chan int(nil),
		v2:   chan int(nil),
	},
	{
		name: "ChanEqualPointer",
		v1:   testChan,
		v2:   testChan,
	},
	{
		name: "ChanNotEqualOnlyOneIsNil",
		v1:   make(chan int),
		v2:   chan int(nil),
	},
	{
		name: "ChanNotEqualCapacity",
		v1:   make(chan int, 1),
		v2:   make(chan int, 2),
	},
	{
		name: "ChanNotEqualLength",
		v1:   make(chan int, 1),
		v2: func() chan int {
			chn := make(chan int, 1)
			chn <- 1
			return chn
		}(),
	},
	{
		name: "FuncEqual",
		v1:   testFunc1,
		v2:   testFunc1,
	},
	{
		name: "FuncEqualNil",
		v1:   (func())(nil),
		v2:   (func())(nil),
	},
	{
		name: "FuncNotEqual",
		v1:   testFunc1,
		v2:   testFunc2,
	},
	{
		name: "TimeEqual",
		v1:   time.Unix(1136239445, 0),
		v2:   time.Unix(1136239445, 0),
	},
	{
		name: "TimeEqualDifferentLocation",
		v1:   time.Unix(1136239445, 0).UTC(),
		v2: func() time.Time {
			loc, err := time.LoadLocation("Europe/Paris")
			if err != nil {
				panic(err)
			}
			return time.Unix(1136239445, 0).In(loc)
		}(),
	},
	{
		name: "TimeNotEqual",
		v1:   time.Unix(1136239445, 0),
		v2:   time.Unix(1136239446, 0),
	},
	{
		name: "NetIPEqual",
		v1:   net.ParseIP("111.111.111.111"),
		v2:   net.ParseIP("111.111.111.111"),
	},
	{
		name: "NetIPNotEqualEqual",
		v1:   net.ParseIP("111.111.111.111"),
		v2:   net.ParseIP("222.222.222.222"),
	},
	{
		name: "MatBigIntEqual",
		v1:   big.NewInt(1),
		v2:   big.NewInt(1),
	},
	{
		name: "MatBigIntNotEqual",
		v1:   big.NewInt(1),
		v2:   big.NewInt(2),
	},
	{
		name: "MathBigRatEqual",
		v1:   big.NewRat(1, 1),
		v2:   big.NewRat(2, 2),
	},
	{
		name: "MathBigIntNotEqual",
		v1:   big.NewRat(1, 1),
		v2:   big.NewRat(2, 1),
	},
	{
		name: "MathBigFloatEqual",
		v1:   big.NewFloat(12.34),
		v2:   big.NewFloat(12.34),
	},
	{
		name: "MathBigFloatNotEqual",
		v1:   big.NewFloat(12.34),
		v2:   big.NewFloat(56.78),
	},
	{
		name: "ReflectValueEqual",
		v1:   reflect.ValueOf(1),
		v2:   reflect.ValueOf(1),
	},
	{
		name: "ReflectValueNotEqual",
		v1:   reflect.ValueOf(1),
		v2:   reflect.ValueOf(2),
	},
	{
		name: "MaxDepth",
		v1: map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": "test1",
				},
			},
		},
		v2: map[string]any{
			"a": map[string]any{
				"b": map[string]any{
					"c": "test2",
				},
			},
		},
		configure: func(c *Comparator) {
			c.MaxDepth = 6
		},
	},
}

func TestCompare(t *testing.T) {
	for _, tc := range compareTestCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.newComparator()
			r := c.Compare(tc.v1, tc.v2)
			assertauto.DeepEqual(t, r)
		})
	}
}

func TestCompareAllocs(t *testing.T) {
	for _, tc := range compareTestCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.newComparator()
			assertauto.AllocsPerRun(t, 100, func() {
				c.Compare(tc.v1, tc.v2)
			})
		})
	}
}

func BenchmarkCompare(b *testing.B) {
	for _, tc := range compareTestCases {
		b.Run(tc.name, func(b *testing.B) {
			c := tc.newComparator()
			for range b.N {
				c.Compare(tc.v1, tc.v2)
			}
		})
	}
}

func TestCompareUnsafePointerNotEqual(t *testing.T) {
	v1 := unsafe.Pointer(&testInt)
	v2 := unsafe.Pointer(&testSlice)
	r := Compare(v1, v2)
	assert.SliceLen(t, r, 1)
	d := r[0]
	assertauto.Equal(t, d.Message)
	assert.NotEqual(t, d.V1, d.V2)
}

var (
	testSlice = []int{1, 2, 3}
	testInt   = 1
	testMap   = map[string]int{"i": 1}
	testChan  = make(chan int)
)

func testFunc1() {}

func testFunc2() {}

type testStruct struct {
	Exported   int
	unexported int
}

var testResult = Result{
	Difference{
		Message: "test1",
		V1:      "1",
		V2:      "2",
	},
	Difference{
		Message: "test2",
		V1:      "a",
		V2:      "b",
	},
}

func TestResultFormat(t *testing.T) {
	s := fmt.Sprintf("%+v", testResult)
	assertauto.Equal(t, s)
}

func TestResultFormatAllocs(t *testing.T) {
	var it any = testResult
	assert.AllocsPerRun(t, 100, func() {
		fmt.Fprintf(io.Discard, "%+v", it)
	}, 0)
}

func BenchmarkResultFormat(b *testing.B) {
	var it any = testResult
	for range b.N {
		_, _ = fmt.Fprintf(io.Discard, "%+v", it)
	}
}

func TestResultFormatEmpty(t *testing.T) {
	var r Result
	s := fmt.Sprintf("%+v", r)
	assertauto.Equal(t, s)
}

func TestResultFormatUnsupportedVerb(t *testing.T) {
	var r Result
	s := fmt.Sprintf("%s", r)
	assertauto.Equal(t, s)
}

func TestDifferenceFormatUnsupportedVerb(t *testing.T) {
	var d Difference
	s := fmt.Sprintf("%s", d)
	assertauto.Equal(t, s)
}

var pathTestCases = []struct {
	name string
	path Path
}{
	{
		name: "Empty",
	},
	{
		name: "Struct",
		path: Path{
			{
				Struct: toPtr("test"),
			},
		},
	},
	{
		name: "Map",
		path: Path{
			{
				Map: toPtr("test"),
			},
		},
	},
	{
		name: "Index",
		path: Path{
			{
				Index: toPtr(1),
			},
		},
	},
	{
		name: "All",
		path: Path{
			{
				Index: toPtr(1),
			},
			{
				Map: toPtr("test"),
			},
			{
				Struct: toPtr("test"),
			},
		},
	},
}

func TestPathFormat(t *testing.T) {
	for _, tc := range pathTestCases {
		t.Run(tc.name, func(t *testing.T) {
			s := fmt.Sprintf("%v", tc.path)
			assertauto.DeepEqual(t, s)
		})
	}
}

func TestPathFormatAllocs(t *testing.T) {
	for _, tc := range pathTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var it any = tc.path
			assert.AllocsPerRun(t, 100, func() {
				fmt.Fprintf(io.Discard, "%v", it)
			}, 0)
		})
	}
}

func BenchmarkPathFormat(b *testing.B) {
	for _, tc := range pathTestCases {
		b.Run(tc.name, func(b *testing.B) {
			var it any = tc.path
			for range b.N {
				_, _ = fmt.Fprintf(io.Discard, "%v", it)
			}
		})
	}
}

func toPtr[V any](v V) *V {
	return &v
}
