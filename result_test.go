package compare

import "testing"

func TestResult(t *testing.T) {
	for _, tc := range []struct {
		name     string
		result   Result
		expected string
	}{
		{
			name:     "None",
			expected: "<none>",
		},
		{
			name: "1",
			result: Result{
				Difference{
					Message: "test",
					V1:      1,
					V2:      2,
				},
			},
			expected: ".: test: v1=1 v2=2",
		},
		{
			name: "2",
			result: Result{
				Difference{
					Message: "test1",
					V1:      1,
					V2:      2,
				},
				Difference{
					Message: "test2",
					V1:      3,
					V2:      4,
				},
			},
			expected: ".: test1: v1=1 v2=2\n.: test2: v1=3 v2=4",
		},
		{
			name: "String",
			result: Result{
				Difference{
					Message: "test",
					V1:      "a",
					V2:      "b",
				},
			},
			expected: ".: test: v1=\"a\" v2=\"b\"",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.result.String()
			if s != tc.expected {
				t.Fatalf("unexpected result: got %q, want %q", s, tc.expected)
			}
		})
	}
}
