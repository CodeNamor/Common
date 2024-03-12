package url

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FormatQuery(t *testing.T) {
	testcases := []struct {
		name        string
		prefix      string
		queryParams []QueryParam
		expected    string
	}{
		{
			name:        "0 - all empty",
			prefix:      "",
			queryParams: []QueryParam{},
			expected:    "",
		},
		{
			name:        "1 - only prefix",
			prefix:      "Foo",
			queryParams: []QueryParam{},
			expected:    "Foo",
		},
		{
			name:   "2 - all populated",
			prefix: "Params: ",
			queryParams: []QueryParam{
				{Key: "foo", Value: "123"},
				{Key: "bar", Value: "456"},
				{Key: "cat", Value: "789"},
			},
			expected: "Params: foo: 123, bar: 456, cat: 789",
		},
		{
			name:   "3 - some populated",
			prefix: "Params: ",
			queryParams: []QueryParam{
				{Key: "empty1", Value: ""},
				{Key: "foo", Value: "123"},
				{Key: "bar", Value: "456"},
				{Key: "empty2", Value: ""},
				{Key: "cat", Value: "789"},
				{Key: "empty3", Value: ""},
			},
			expected: "Params: foo: 123, bar: 456, cat: 789",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatQuery(tc.prefix, tc.queryParams)
			require.Equal(t, tc.expected, result)
		})
	}
}
