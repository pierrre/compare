package compare_test

import (
	"testing"

	"github.com/pierrre/assert"
	"github.com/pierrre/assert/assertauto"
	. "github.com/pierrre/compare"
)

var testPath = Path{
	{
		Index: new(123),
	},
	{
		Key: new(any("test")),
	},
	{
		Name: new("test"),
	},
	{
		Pointer: true,
	},
}

func TestPathString(t *testing.T) {
	s := testPath.String()
	assertauto.Equal(t, s)
}

func TestPathStringEmpty(t *testing.T) {
	s := Path(nil).String()
	assertauto.Equal(t, s)
}

func TestPathStringAllocs(t *testing.T) {
	assertauto.AllocsPerRun(t, 100, func() {
		_ = testPath.String()
	})
}

func BenchmarkPathString(b *testing.B) {
	for b.Loop() {
		_ = testPath.String()
	}
}

func TestPathAppendText(t *testing.T) {
	var b []byte
	b, err := testPath.AppendText(b)
	assert.NoError(t, err)
	assertauto.Equal(t, string(b))
}

func TestPathAppendTextAllocs(t *testing.T) {
	var b []byte
	assertauto.AllocsPerRun(t, 100, func() {
		b, _ = testPath.AppendText(b[:0])
	})
}

func BenchmarkPathAppendText(b *testing.B) {
	var buf []byte
	for b.Loop() {
		buf, _ = testPath.AppendText(buf[:0])
	}
}

func TestPathElemString(t *testing.T) {
	for _, e := range testPath {
		s := e.String()
		assertauto.Equal(t, s)
	}
}

func TestPathElemStringAllocs(t *testing.T) {
	for _, e := range testPath {
		assertauto.AllocsPerRun(t, 100, func() {
			_ = e.String()
		})
	}
}

func BenchmarkPathElemString(b *testing.B) {
	for _, e := range testPath {
		b.Run("", func(b *testing.B) {
			for b.Loop() {
				_ = e.String()
			}
		})
	}
}

func TestPathElemAppendText(t *testing.T) {
	for _, e := range testPath {
		var b []byte
		b, err := e.AppendText(b)
		assert.NoError(t, err)
		assertauto.Equal(t, string(b))
	}
}

func TestPathElemAppendTextAllocs(t *testing.T) {
	for _, e := range testPath {
		var b []byte
		assertauto.AllocsPerRun(t, 100, func() {
			b, _ = e.AppendText(b[:0])
		})
	}
}

func BenchmarkPathElemAppendText(b *testing.B) {
	for _, e := range testPath {
		b.Run("", func(b *testing.B) {
			var buf []byte
			for b.Loop() {
				buf, _ = e.AppendText(buf[:0])
			}
		})
	}
}
