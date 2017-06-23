package compare

import (
	"fmt"
	"testing"
)

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
			path: StructPath{
				Field: "test",
			},
			expected: ".test",
		},
		{
			name: "Map",
			path: MapPath{
				Key: "test",
			},
			expected: "[test]",
		},
		{
			name: "Indexed",
			path: IndexedPath{
				Index: 1,
			},
			expected: "[1]",
		},
		{
			name: "All",
			path: StructPath{
				Field: "test",
				Next: MapPath{
					Key: "test",
					Next: IndexedPath{
						Index: 1,
					},
				},
			},
			expected: ".test[test][1]",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var s string
			if tc.path != nil {
				s = tc.path.(fmt.Stringer).String()
			} else {
				s = PathString(tc.path)
			}
			if s != tc.expected {
				t.Fatalf("unexpected result: got %q, want %q", s, tc.expected)
			}
		})
	}
}
