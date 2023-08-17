package compare

import (
	"testing"

	"github.com/pierrre/assert"
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
			path: Path{
				StructPathElem{
					Field: "test",
				},
			},
			expected: ".test",
		},
		{
			name: "Map",
			path: Path{
				MapPathElem{
					Key: "test",
				},
			},
			expected: "[test]",
		},
		{
			name: "Indexed",
			path: Path{
				IndexedPathElem{
					Index: 1,
				},
			},
			expected: "[1]",
		},
		{
			name: "All",
			path: Path{
				IndexedPathElem{
					Index: 1,
				},
				MapPathElem{
					Key: "test",
				},
				StructPathElem{
					Field: "test",
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
