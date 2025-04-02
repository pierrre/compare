package compare_test

import (
	"testing"

	. "github.com/pierrre/compare/refactor"
)

func Test(t *testing.T) {
	vc := &IntComparator{}
	c := &Comparator{
		ValuesComparator: vc,
	}
	v1 := int64(1)
	v2 := int64(2)
	res := c.Compare(v1, v2)
	t.Log(res)
	// t.Fatal("aaa")
}

func Benchmark(b *testing.B) {
	vc := &IntComparator{}
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
