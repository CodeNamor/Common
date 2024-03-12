package errors

import (
	"encoding/json"
	"fmt"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_New(t *testing.T) {
	testcases := []struct {
		name     string
		errstr   string
		expected ErrorLog
	}{
		{
			name:   "with string",
			errstr: "mystring",
			expected: ErrorLog{
				Err: pkgerrors.New("mystring"),
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected.Error(), New(tc.errstr).Error())
		})
	}
}

func Test_NewRootMsgStatusCode(t *testing.T) {
	testcases := []struct {
		name       string
		rootCause  string
		msg        string
		statusCode string
		expected   string
	}{
		{
			name:       "empty root, msg, statuscode",
			rootCause:  "",
			msg:        "",
			statusCode: "",
			expected:   "",
		},
		{
			name:       "root",
			rootCause:  "myroot",
			msg:        "",
			statusCode: "",
			expected:   "myroot",
		},
		{
			name:       "root, msg",
			rootCause:  "myroot",
			msg:        "mymsg",
			statusCode: "",
			expected:   "myroot mymsg",
		},
		{
			name:       "root, msg, status",
			rootCause:  "myroot",
			msg:        "mymsg",
			statusCode: "206",
			expected:   "myroot mymsg StatusCode:206",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, NewRootMsgStatusCode(tc.rootCause, tc.msg, tc.statusCode).Error())
		})
	}
}

func Test_Errorf(t *testing.T) {
	testcases := []struct {
		name     string
		format   string
		args     []interface{}
		expected ErrorLog
	}{
		{
			name:   "only format",
			format: "myformat",
			expected: ErrorLog{
				Err: pkgerrors.New("myformat"),
			},
		},
		{
			name:   "format, 1 arg",
			format: "myformat %d",
			args:   []interface{}{10},
			expected: ErrorLog{
				Err: pkgerrors.New("myformat 10"),
			},
		},
		{
			name:   "format, 2 args",
			format: "myformat %d %s",
			args:   []interface{}{10, "anotherarg"},
			expected: ErrorLog{
				Err: pkgerrors.New("myformat 10 anotherarg"),
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected.Error(), Errorf(tc.format, tc.args...).Error())
		})
	}
}
func Test_FromError(t *testing.T) {
	err1 := pkgerrors.New("myerror1")
	testcases := []struct {
		name     string
		err      error
		expected *ErrorLog
	}{
		{
			name:     "nil err",
			err:      nil,
			expected: nil,
		},
		{
			name: "with err",
			err:  err1,
			expected: &ErrorLog{
				Err: err1,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, FromError(tc.err))
		})
	}
}

func Test_WithErrorAndCause(t *testing.T) {
	err1 := pkgerrors.New("myerror1")
	testcases := []struct {
		name      string
		err       error
		rootCause string
		expected  *ErrorLog
	}{
		{
			name:      "nil err, empty cause",
			err:       nil,
			rootCause: "",
			expected:  nil,
		},
		{
			name:      "nil err, non-empty cause",
			err:       nil,
			rootCause: "just because",
			expected:  nil,
		},
		{
			name:      "with err, empty cause",
			err:       err1,
			rootCause: "",
			expected: &ErrorLog{
				RootCause: "",
				Err:       err1,
			},
		},
		{
			name:      "with err, non-empty cause",
			err:       err1,
			rootCause: "just because",
			expected: &ErrorLog{
				RootCause: "just because",
				Err:       err1,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, WithErrorAndCause(tc.err, tc.rootCause))
		})
	}
}

func Test_ErrorAndString(t *testing.T) {
	testcases := []struct {
		name     string
		errorLog ErrorLog
		expected string
	}{
		{
			name: "only err",
			errorLog: ErrorLog{
				Err: fmt.Errorf("myerror"),
			},
			expected: "myerror",
		},
		{
			name: "err, root",
			errorLog: ErrorLog{
				RootCause: "myroot",
				Err:       fmt.Errorf("myerror"),
			},
			expected: "myroot myerror",
		},
		{
			name: "trace, err",
			errorLog: ErrorLog{
				Trace: "mytrace",
				Err:   fmt.Errorf("myerror"),
			},
			expected: "mytrace myerror",
		},
		{
			name: "trace",
			errorLog: ErrorLog{
				Trace: "mytrace",
			},
			expected: "mytrace",
		},
		{
			name: "all fields",
			errorLog: ErrorLog{
				RootCause:             "myroot",
				Trace:                 "mytrace",
				StatusCode:            "206",
				Source:                "mysource",
				Scope:                 "myscope",
				Query:                 "myquery",
				AdditionalInformation: "myaddinfo",
				ExceptionType:         "myexc",
				Err:                   fmt.Errorf("myerror"),
			},
			expected: "myroot mytrace myerror StatusCode:206 Source:mysource Scope:myscope Query:myquery AdditionalInformation:myaddinfo ExceptionType:myexc",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.errorLog.Error())
			require.Equal(t, tc.expected, tc.errorLog.String())
		})
	}
}

func Test_errorLogToErrorString(t *testing.T) {
	testcases := []struct {
		name               string
		errorLog           ErrorLog
		expectedIncludeErr string
		expectedExcludeErr string
	}{
		{
			name: "only err",
			errorLog: ErrorLog{
				Err: fmt.Errorf("myerror"),
			},
			expectedIncludeErr: "myerror",
			expectedExcludeErr: "",
		},
		{
			name: "err, root",
			errorLog: ErrorLog{
				RootCause: "myroot",
				Err:       fmt.Errorf("myerror"),
			},
			expectedIncludeErr: "myroot myerror",
			expectedExcludeErr: "myroot",
		},
		{
			name: "trace, err",
			errorLog: ErrorLog{
				Trace: "mytrace",
				Err:   fmt.Errorf("myerror"),
			},
			expectedIncludeErr: "mytrace myerror",
			expectedExcludeErr: "mytrace",
		},
		{
			name: "trace",
			errorLog: ErrorLog{
				Trace: "mytrace",
			},
			expectedIncludeErr: "mytrace",
			expectedExcludeErr: "mytrace",
		},
		{
			name: "all fields",
			errorLog: ErrorLog{
				RootCause:             "myroot",
				Trace:                 "mytrace",
				StatusCode:            "206",
				Source:                "mysource",
				Scope:                 "myscope",
				Query:                 "myquery",
				AdditionalInformation: "myaddinfo",
				ExceptionType:         "myexc",
				Err:                   fmt.Errorf("myerror"),
			},
			expectedIncludeErr: "myroot mytrace myerror StatusCode:206 Source:mysource Scope:myscope Query:myquery AdditionalInformation:myaddinfo ExceptionType:myexc",
			expectedExcludeErr: "myroot mytrace StatusCode:206 Source:mysource Scope:myscope Query:myquery AdditionalInformation:myaddinfo ExceptionType:myexc",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedIncludeErr, errorLogToErrorString(true, &tc.errorLog))
			require.Equal(t, tc.expectedExcludeErr, errorLogToErrorString(false, &tc.errorLog))
		})
	}
}

func Test_Format(t *testing.T) {
	testcases := []struct {
		name          string
		errorLog      ErrorLog
		expected      string
		expectedPlusV string
	}{
		{
			name: "only err",
			errorLog: ErrorLog{
				Err: fmt.Errorf("myerror"),
			},
			expected:      "myerror",
			expectedPlusV: "myerror",
		},
		{
			name: "err, root",
			errorLog: ErrorLog{
				RootCause: "myroot",
				Err:       fmt.Errorf("myerror"),
			},
			expected:      "myroot myerror",
			expectedPlusV: "myroot myerror",
		},
		{
			name: "trace",
			errorLog: ErrorLog{
				Trace: "mytrace",
			},
			expected:      "mytrace",
			expectedPlusV: "mytrace",
		},
		{
			name: "trace, err",
			errorLog: ErrorLog{
				Trace: "mytrace",
				Err:   fmt.Errorf("myerror"),
			},
			expected:      "mytrace myerror",
			expectedPlusV: "mytrace myerror",
		},
		{
			name: "all fields",
			errorLog: ErrorLog{
				RootCause:             "myroot",
				Trace:                 "mytrace",
				StatusCode:            "206",
				Source:                "mysource",
				Scope:                 "myscope",
				Query:                 "myquery",
				AdditionalInformation: "myaddinfo",
				ExceptionType:         "myexc",
				Err:                   fmt.Errorf("myerror"),
			},
			expected:      "myroot mytrace myerror StatusCode:206 Source:mysource Scope:myscope Query:myquery AdditionalInformation:myaddinfo ExceptionType:myexc",
			expectedPlusV: "myroot mytrace StatusCode:206 Source:mysource Scope:myscope Query:myquery AdditionalInformation:myaddinfo ExceptionType:myexc myerror",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// normal v's
			v := fmt.Sprintf("%v", tc.errorLog)
			require.Equal(t, tc.expected, v)
			v = fmt.Sprintf("%v", &tc.errorLog)
			require.Equal(t, tc.expected, v, "using pointer")

			// plus v's
			pv := fmt.Sprintf("%+v", tc.errorLog)
			require.Equal(t, tc.expectedPlusV, pv)
			pv = fmt.Sprintf("%+v", &tc.errorLog)
			require.Equal(t, tc.expectedPlusV, pv, "using pointer")
		})
	}
}

func Test_JSONMarshal(t *testing.T) {
	testcases := []struct {
		name     string
		errorLog ErrorLog
		expected string
	}{
		{
			name: "only err",
			errorLog: ErrorLog{
				Err: fmt.Errorf("myerror"),
			},
			expected: "{\"Trace\":\"myerror\"}",
		},
		{
			name: "err, root",
			errorLog: ErrorLog{
				RootCause: "myroot",
				Err:       fmt.Errorf("myerror"),
			},
			expected: "{\"Trace\":\"myerror\",\"RootCause\":\"myroot\"}",
		},
		{
			name: "trace",
			errorLog: ErrorLog{
				Trace: "mytrace",
			},
			expected: "{\"Trace\":\"mytrace\"}",
		},
		{
			name: "trace, err",
			errorLog: ErrorLog{
				Trace: "mytrace",
				Err:   fmt.Errorf("myerror"),
			},
			expected: "{\"Trace\":\"mytrace myerror\"}",
		},
		{
			name: "all fields",
			errorLog: ErrorLog{
				RootCause:             "myroot",
				Trace:                 "mytrace",
				StatusCode:            "206",
				Source:                "mysource",
				Scope:                 "myscope",
				Query:                 "myquery",
				AdditionalInformation: "myaddinfo",
				ExceptionType:         "myexc",
				Err:                   fmt.Errorf("myerror"),
			},
			expected: "{\"Trace\":\"mytrace myerror\",\"RootCause\":\"myroot\",\"StatusCode\":\"206\",\"Source\":\"mysource\",\"Scope\":\"myscope\",\"Query\":\"myquery\",\"AdditionalInformation\":\"myaddinfo\",\"ExceptionType\":\"myexc\"}",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.errorLog)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(bytes))

			// repeat using pointer to ErrorLog
			bytes, err = json.Marshal(&tc.errorLog)
			require.NoError(t, err, "using pointer")
			require.Equal(t, tc.expected, string(bytes), "using pointer")
		})
	}
}

func Test_combineStringAndError(t *testing.T) {
	testcases := []struct {
		name     string
		str      string
		err      error
		expected string
	}{
		{
			name:     "empty string, nil err",
			str:      "",
			err:      nil,
			expected: "",
		},
		{
			name:     "string, nil err",
			str:      "mystr",
			err:      nil,
			expected: "mystr",
		},
		{
			name:     "string, empty err",
			str:      "mystr",
			err:      pkgerrors.New(""),
			expected: "mystr",
		},
		{
			name:     "empty string, err",
			str:      "",
			err:      pkgerrors.New("myerror"),
			expected: "myerror",
		},
		{
			name:     "string, err",
			str:      "mystr",
			err:      pkgerrors.New("myerror"),
			expected: "mystr myerror",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, combineStringAndError(tc.str, tc.err))
		})
	}
}
