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
