package compare

import (
	"bytes"
	"fmt"
	"image"
	"math/big"
	"net"
	"reflect"
	"testing"
	"time"
	"unsafe" //nolint:depguard // Used for unsafe.Pointer comparison.

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/ext/pierrreerrors"
	"github.com/pierrre/assert/ext/pierrrepretty"
)

func init() {
	// Prevent import cycle.
	assert.DeepEqualer = func(v1, v2 any) (diff string, equal bool) {
		res := Compare(v1, v2)
		if len(res) == 0 {
			return "", true
		}
		diff = fmt.Sprintf("%+v", res)
		return diff, false
	}
	pierrrepretty.ConfigureDefault()
	pierrreerrors.Configure()
}

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

var compareTestCases = []struct {
	name     string
	v1       any
	v2       any
	expected Result
}{
	{
		name: "EqualNotValid",
		v1:   nil,
		v2:   nil,
	},
	{
		name: "NotEqualOnlyOneIsValid",
		v1:   nil,
		v2:   true,
		expected: Result{
			Difference{
				Message: msgOnlyOneIsValid,
				V1:      false,
				V2:      true,
			},
		},
	},
	{
		name: "NotEqualDifferentType",
		v1:   int32(1),
		v2:   int64(1),
		expected: Result{
			Difference{
				Message: msgTypeNotEqual,
				V1:      reflect.TypeOf(int32(0)),
				V2:      reflect.TypeOf(int64(0)),
			},
		},
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
		expected: Result{
			Difference{
				Message: msgBoolNotEqual,
				V1:      true,
				V2:      false,
			},
		},
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
		expected: Result{
			Difference{
				Message: msgIntNotEqual,
				V1:      int64(1),
				V2:      int64(2),
			},
		},
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
		expected: Result{
			Difference{
				Message: msgUintNotEqual,
				V1:      uint64(1),
				V2:      uint64(2),
			},
		},
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
		expected: Result{
			Difference{
				Message: msgFloatNotEqual,
				V1:      float64(1),
				V2:      float64(2),
			},
		},
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
		expected: Result{
			Difference{
				Message: msgComplexNotEqual,
				V1:      complex(1, 1),
				V2:      complex(2, 2),
			},
		},
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
		expected: Result{
			Difference{
				Message: msgStringNotEqual,
				V1:      "a",
				V2:      "b",
			},
		},
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
		expected: Result{
			Difference{
				Path: Path{
					{
						Index: toPtr(1),
					},
				},
				Message: msgIntNotEqual,
				V1:      int64(2),
				V2:      int64(0),
			},
		},
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
		name: "SliceNotEqual",
		v1:   []int{1, 2, 3},
		v2:   []int{1, 0, 3},
		expected: Result{
			Difference{
				Path: Path{
					{
						Index: toPtr(1),
					},
				},
				Message: msgIntNotEqual,
				V1:      int64(2),
				V2:      int64(0),
			},
		},
	},
	{
		name: "SliceNotEqualOnlyOneIsNil",
		v1:   []int{1, 2, 3},
		v2:   []int(nil),
		expected: Result{
			Difference{
				Message: msgOnlyOneIsNil,
				V1:      false,
				V2:      true,
			},
		},
	},
	{
		name: "SliceNotEqualLength",
		v1:   []int{1, 2},
		v2:   []int{1, 2, 3},
		expected: Result{
			Difference{
				Message: msgLengthNotEqual,
				V1:      2,
				V2:      3,
			},
		},
	},
	{
		name: "SliceByteNotEqual",
		v1:   make([]byte, 1<<20),
		v2: func() []byte {
			s := make([]byte, 1<<20)
			s[0] = 1
			return s
		}(),
		expected: Result{
			Difference{
				Path: Path{
					{
						Index: toPtr(0),
					},
				},
				Message: msgUintNotEqual,
				V1:      uint64(0),
				V2:      uint64(1),
			},
		},
	},
	{
		name: "SliceNotEqualMaxDifferences",
		v1: func() []int {
			s := make([]int, MaxSliceDifferences*2)
			for i := range s {
				s[i] = i
			}
			return s
		}(),
		v2: func() []int {
			s := make([]int, MaxSliceDifferences*2)
			for i := range s {
				s[i] = i + 1
			}
			return s
		}(),
		expected: func() Result {
			r := make(Result, MaxSliceDifferences)
			for i := range r {
				r[i] = Difference{
					Path: Path{
						{
							Index: toPtr(i),
						},
					},
					Message: msgIntNotEqual,
					V1:      int64(i),
					V2:      int64(i + 1),
				}
			}
			return r
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
		expected: Result{
			Difference{
				Path: Path{
					{
						Index: toPtr(0),
					},
				},
				Message: msgOnlyOneIsNil,
				V1:      false,
				V2:      true,
			},
		},
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
		name: "PointerNotEqual",
		v1: func() *int {
			i := 1
			return &i
		}(),
		v2: func() *int {
			i := 2
			return &i
		}(),
		expected: Result{
			Difference{
				Message: msgIntNotEqual,
				V1:      int64(1),
				V2:      int64(2),
			},
		},
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
		expected: Result{
			Difference{
				Path: Path{
					{
						Struct: toPtr("Exported"),
					},
				},
				Message: msgIntNotEqual,
				V1:      int64(1),
				V2:      int64(2),
			},
		},
	},
	{
		name: "StructEqualNotEqualUnexported",
		v1: &testStruct{
			Exported:   1,
			unexported: 1,
		},
		v2: &testStruct{
			Exported:   1,
			unexported: 2,
		},
		expected: Result{
			Difference{
				Path: Path{
					{
						Struct: toPtr("unexported"),
					},
				},
				Message: msgIntNotEqual,
				V1:      int64(1),
				V2:      int64(2),
			},
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
		name: "MapNotEqualValue",
		v1: map[string]int{
			"i": 1,
		},
		v2: map[string]int{
			"i": 2,
		},
		expected: Result{
			Difference{
				Path: Path{
					{
						Map: toPtr("i"),
					},
				},
				Message: msgIntNotEqual,
				V1:      int64(1),
				V2:      int64(2),
			},
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
		expected: Result{
			Difference{
				Path: Path{
					{
						Map: toPtr("a"),
					},
				},
				Message: msgMapKeyNotDefined,
				V1:      true,
				V2:      false,
			},
			Difference{
				Path: Path{
					{
						Map: toPtr("b"),
					},
				},
				Message: msgMapKeyNotDefined,
				V1:      false,
				V2:      true,
			},
		},
	},
	{
		name: "MapNotEqualOnlyOneIsNil",
		v1: map[string]int{
			"i": 1,
		},
		v2: map[string]int(nil),
		expected: Result{
			Difference{
				Message: msgOnlyOneIsNil,
				V1:      false,
				V2:      true,
			},
		},
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
		expected: Result{
			Difference{
				Message: msgLengthNotEqual,
				V1:      1,
				V2:      2,
			},
		},
	},
	{
		name: "UnsafePointerEqual",
		v1:   unsafe.Pointer(&testInt), //nolint:gosec // Ignore for testing.
		v2:   unsafe.Pointer(&testInt), //nolint:gosec // Ignore for testing.
	},
	{
		name: "UnsafePointerNotEqual",
		v1:   unsafe.Pointer(&testInt),   //nolint:gosec // Ignore for testing.
		v2:   unsafe.Pointer(&testSlice), //nolint:gosec // Ignore for testing.
		expected: Result{
			Difference{
				Message: msgUnsafePointerNotEqual,
				V1:      uintptr(unsafe.Pointer(&testInt)),   //nolint:gosec // Ignore for testing.
				V2:      uintptr(unsafe.Pointer(&testSlice)), //nolint:gosec // Ignore for testing.
			},
		},
	},
	{
		name: "ChanEqual",
		v1:   make(chan int),
		v2:   make(chan int),
	},
	{
		name: "ChanEqualNil",
		v1:   [1]chan int{},
		v2:   [1]chan int{},
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
		expected: Result{
			Difference{
				Message: msgOnlyOneIsNil,
				V1:      false,
				V2:      true,
			},
		},
	},
	{
		name: "ChanNotEqualCapacity",
		v1:   make(chan int, 1),
		v2:   make(chan int, 2),
		expected: Result{
			Difference{
				Message: msgCapacityNotEqual,
				V1:      1,
				V2:      2,
			},
		},
	},
	{
		name: "ChanNotEqualLength",
		v1:   make(chan int, 1),
		v2: func() chan int {
			chn := make(chan int, 1)
			chn <- 1
			return chn
		}(),
		expected: Result{
			Difference{
				Message: msgLengthNotEqual,
				V1:      0,
				V2:      1,
			},
		},
	},
	{
		name: "FuncEqual",
		v1:   testFunc,
		v2:   testFunc,
	},
	{
		name: "FuncEqualNil",
		v1:   [1]func(){},
		v2:   [1]func(){},
	},
	{
		name: "FuncNotEqual",
		v1:   [1]func(){},
		v2:   [1]func(){testFunc},
		expected: Result{
			Difference{
				Path: Path{
					{
						Index: toPtr(0),
					},
				},
				Message: msgFuncPointerNotEqual,
				V1:      uintptr(0),
				V2:      reflect.ValueOf(testFunc).Pointer(),
			},
		},
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
		expected: Result{
			Difference{
				Message: fmt.Sprintf(msgMethodNotEqual, "Equal"),
				V1:      time.Unix(1136239445, 0),
				V2:      time.Unix(1136239446, 0),
			},
		},
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
		expected: Result{
			Difference{
				Message: fmt.Sprintf(msgMethodNotEqual, "Equal"),
				V1:      net.ParseIP("111.111.111.111"),
				V2:      net.ParseIP("222.222.222.222"),
			},
		},
	},
	{
		name: "ImageRectangeEqual",
		v1:   image.Rect(0, 0, 1, 1),
		v2:   image.Rect(0, 0, 1, 1),
	},
	{
		name: "ImageRectangeEqualEmpty",
		v1:   image.Rect(1, 1, 1, 1),
		v2:   image.Rect(2, 2, 2, 2),
	},
	{
		name: "ImageRectangeNotEqual",
		v1:   image.Rect(0, 0, 1, 1),
		v2:   image.Rect(0, 0, 2, 2),
		expected: Result{
			Difference{
				Message: fmt.Sprintf(msgMethodNotEqual, "Eq"),
				V1:      image.Rect(0, 0, 1, 1),
				V2:      image.Rect(0, 0, 2, 2),
			},
		},
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
		expected: Result{
			Difference{
				Message: fmt.Sprintf(msgMethodCmpNotEqual, -1),
				V1:      big.NewInt(1),
				V2:      big.NewInt(2),
			},
		},
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
		expected: Result{
			Difference{
				Message: fmt.Sprintf(msgMethodCmpNotEqual, -1),
				V1:      big.NewRat(1, 1),
				V2:      big.NewRat(2, 1),
			},
		},
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
		expected: Result{
			Difference{
				Message: fmt.Sprintf(msgMethodCmpNotEqual, -1),
				V1:      big.NewFloat(12.34),
				V2:      big.NewFloat(56.78),
			},
		},
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
		expected: Result{
			Difference{
				Message: msgIntNotEqual,
				V1:      int64(1),
				V2:      int64(2),
			},
		},
	},
}

func TestCompare(t *testing.T) {
	for _, tc := range compareTestCases {
		t.Run(tc.name, func(t *testing.T) {
			r := Compare(tc.v1, tc.v2)
			assert.DeepEqual(t, r, tc.expected)
		})
	}
}

func BenchmarkCompare(b *testing.B) {
	for _, tc := range compareTestCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Compare(tc.v1, tc.v2)
			}
		})
	}
}

var (
	testSlice = []int{1, 2, 3}
	testInt   = 1
	testMap   = map[string]int{"i": 1}
	testChan  = make(chan int)
	testFunc  = func() {}
)

type testStruct struct {
	Exported   int
	unexported int
}

var testResult = Result{
	Difference{
		Message: "test1",
		V1:      1,
		V2:      2,
	},
	Difference{
		Message: "test2",
		V1:      "a",
		V2:      "b",
	},
}

func TestResultFormat(t *testing.T) {
	s := fmt.Sprintf("%+v", testResult)
	expected := ".: test1\n\tv1=1\n\tv2=2\n.: test2\n\tv1=\"a\"\n\tv2=\"b\""
	assert.Equal(t, s, expected)
}

func BenchmarkResultFormat(b *testing.B) {
	var it any = testResult
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = fmt.Fprintf(buf, "%+v", it)
	}
}

func TestResultFormatEmpty(t *testing.T) {
	var r Result
	s := fmt.Sprintf("%+v", r)
	expected := "<none>"
	assert.Equal(t, s, expected)
}

func TestResultFormatUnsupportedVerb(t *testing.T) {
	var r Result
	s := fmt.Sprintf("%s", r)
	expected := "%!s(compare.Result)"
	assert.Equal(t, s, expected)
}

func TestDifferenceFormatUnsupportedVerb(t *testing.T) {
	var d Difference
	s := fmt.Sprintf("%s", d)
	expected := "%!s(compare.Difference)"
	assert.Equal(t, s, expected)
}

func TestPathString(t *testing.T) {
	for _, tc := range []struct {
		name     string
		path     Path
		expected string
	}{
		{
			name:     "Empty",
			expected: ".",
		},
		{
			name: "Struct",
			path: Path{
				{
					Struct: toPtr("test"),
				},
			},
			expected: ".test",
		},
		{
			name: "Map",
			path: Path{
				{
					Map: toPtr("test"),
				},
			},
			expected: "[test]",
		},
		{
			name: "Index",
			path: Path{
				{
					Index: toPtr(1),
				},
			},
			expected: "[1]",
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
			expected: ".test[test][1]",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.path.String()
			assert.Equal(t, s, tc.expected)
		})
	}
}

var sortValuesTestCases = []struct {
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

func TestSortValues(t *testing.T) {
	for _, tc := range sortValuesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			sortValues(tc.values, tc.typ)
			assert.DeepEqual(t, tc.values, tc.expected)
		})
	}
}
