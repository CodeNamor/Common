package transform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MyStruct struct {
	Foo string
	Bar int
	Cat string
}

func newDefaultMyStruct() *MyStruct {
	return &MyStruct{
		Foo: "abc",
		Bar: 123,
	}
}

func Test_DecodeJSONOntoDefaultStruct(t *testing.T) {
	testcases := []struct {
		name            string
		bytes           []byte
		targetStructPtr *MyStruct
		expectedResult  MyStruct
		expectedError   string
	}{
		{
			name:          "1 - valid bytes, nil targetStructPtr",
			bytes:         []byte(`{"Foo":"abc"}`),
			expectedError: "result must be addressable (a pointer)",
		},
		{
			name:            "2 - nil bytes",
			targetStructPtr: newDefaultMyStruct(),
			expectedError:   "unexpected end of JSON input",
		},
		{
			name:            "3 - empty bytes",
			targetStructPtr: newDefaultMyStruct(),
			expectedError:   "unexpected end of JSON input",
		},
		{
			name:            "4 - valid bytes",
			bytes:           []byte(`{"Foo":"abc"}`),
			targetStructPtr: newDefaultMyStruct(),
			expectedResult: MyStruct{
				Foo: "abc",
				Bar: 123,
				Cat: "",
			},
		},
		{
			name:            "5 - non-matching fields in struct",
			bytes:           []byte(`{"Foo": "abc", "Oops": "bad", "Another": "alsobad"}`),
			targetStructPtr: newDefaultMyStruct(),
			expectedError:   "errored trying to match JSON data to struct: 1 error(s) decoding:\n\n* '' has invalid keys: Another, Oops",
		},
		{
			name:            "6 - invalid JSON",
			bytes:           []byte(`{Foo: "abc"}`),
			targetStructPtr: newDefaultMyStruct(),
			expectedError:   "invalid character 'F' looking for beginning of object key string",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := DecodeJSONToStruct(tc.bytes, tc.targetStructPtr)
			if tc.expectedError != "" { // expected to fail
				require.EqualError(t, err, tc.expectedError)
			} else { // expected to succeed
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, *tc.targetStructPtr)
			}
		})
	}
}
