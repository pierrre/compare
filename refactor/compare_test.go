package compare_test

import (
	"fmt"
	"testing"

	. "github.com/pierrre/compare/refactor"
)

func Test(t *testing.T) {
	vc := &IntValuesComparator{}
	c := &Comparator{
		ValuesComparator: vc,
	}
	v1 := int64(1)
	v2 := int64(2)
	result := c.Compare(v1, v2)
	for d := range result {
		fmt.Println(d.Message)
	}
	// t.Fatal("aaa")
}

func Benchmark(b *testing.B) {
	vc := &IntValuesComparator{}
	c := &Comparator{
		ValuesComparator: vc,
	}
	v1 := int64(1)
	v2 := int64(2)
	for b.Loop() {
		result := c.Compare(v1, v2)
		for d := range result {
			_ = d
		}
	}
}
