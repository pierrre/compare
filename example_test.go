package compare

import "fmt"

func ExampleCompare() {
	type T struct {
		String string
		Int    int
		Map    map[string]interface{}
		Slice  []int
	}
	v1 := T{
		String: "aaa",
		Int:    1,
		Map: map[string]interface{}{
			"a": "a",
			"b": "b",
			"c": "c",
		},
		Slice: []int{1, 2, 3},
	}
	v2 := T{
		String: "bbb",
		Int:    2,
		Map: map[string]interface{}{
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
