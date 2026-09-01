package compare_test

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	. "github.com/pierrre/compare"
)

var testResult = Result{
	{
		Path: Path{
			{Index: new(123)},
		},
		Message: "not equal",
		V1:      1,
		V2:      2,
	},
	{
		Message: "not equal",
		V1:      true,
		V2:      false,
	},
}

var testDifference = testResult[0]

func TestResultString(t *testing.T) {
	s := testResult.String()
	assertauto.Equal(t, s)
}

func TestResultStringEmpty(t *testing.T) {
	s := Result(nil).String()
	assertauto.Equal(t, s)
}

func TestResultStringAllocs(t *testing.T) {
	assertauto.AllocsPerRun(t, 100, func() {
		_ = testResult.String()
	})
}

func BenchmarkResultString(b *testing.B) {
	for b.Loop() {
		_ = testResult.String()
	}
}

func TestResultAppendText(t *testing.T) {
	var b []byte
	b, err := testResult.AppendText(b)
	assert.NoError(t, err)
	assertauto.Equal(t, string(b))
}

func TestResultAppendTextAllocs(t *testing.T) {
	var b []byte
	assertauto.AllocsPerRun(t, 100, func() {
		b, _ = testResult.AppendText(b[:0])
	})
}

func BenchmarkResultAppendText(b *testing.B) {
	var buf []byte
	for b.Loop() {
		buf, _ = testResult.AppendText(buf[:0])
	}
}

func TestDifferenceString(t *testing.T) {
	s := testDifference.String()
	assertauto.Equal(t, s)
}

func TestDifferenceStringAllocs(t *testing.T) {
	assertauto.AllocsPerRun(t, 100, func() {
		_ = testDifference.String()
	})
}

func BenchmarkDifferenceString(b *testing.B) {
	for b.Loop() {
		_ = testDifference.String()
	}
}

func TestDifferenceAppendText(t *testing.T) {
	var b []byte
	b, err := testDifference.AppendText(b)
	assert.NoError(t, err)
	assertauto.Equal(t, string(b))
}

func TestDifferenceAppendTextAllocs(t *testing.T) {
	var b []byte
	assertauto.AllocsPerRun(t, 100, func() {
		b, _ = testDifference.AppendText(b[:0])
	})
}

func BenchmarkDifferenceAppendText(b *testing.B) {
	var buf []byte
	for b.Loop() {
		buf, _ = testDifference.AppendText(buf[:0])
	}
}
